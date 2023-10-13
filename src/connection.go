package main

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"strings"

	// "log"
	"os"
	"time"

	// "github.com/fatih/color"
)

// Connect to the database and make some checkings
func connect(login string, password string, ip string, db_name string) (*sql.DB, error) {
	log.Printf("Connecting to database...\n")

	var database *sql.DB
	var err error

	database, err = sql.Open("mysql", login+":"+password+"@tcp("+ip+")/"+db_name)
	if err != nil {
		log.Printf("Error while connecting to database\n")
		return nil, err
	}

	// Ping the database to check if the connection is working
	err = database.Ping()
	if err != nil {
		log.Printf("Error while Ping the database\n")
		return nil, err
	}
	log.Printf("Ping is ok")

	// Check if the database exists
	log.Printf("Check if database exist")
	query := fmt.Sprintf("SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = '%s'", db_name)
	var result string
	err = database.QueryRow(query).Scan(&result)
	if err != nil {
		if err == sql.ErrNoRows {
			// Database does not exist
			log.Printf("Database does not exist\n")
			return nil, nil
		}
		return nil, err
	}
	log.Printf("Databse exist")
	log.Printf("Start to Tranfert Data")

	database, err = transferInfo(login, password, ip, db_name)
	if err != nil {
		log.Printf("Error in transferInfo function\n")
		return nil, err
	}

	// defer database.Close()
	return database, nil
}

// Create a new database with feeding the tables
func transferInfo(login string, password string, ip string, db_name string) (*sql.DB, error) {

	log.Printf("Creating database...\n")
	database, err := initDatabase(login, password, ip, db_name)
	if err != nil {
		log.Printf("Error while creating database\n")
		return nil, err
	}

	err = createTables(database)
	if err != nil {
		log.Printf("Error while creating tables\n")
		return nil, err
	}

	log.Printf("Tables created\n")
	log.Printf("Transferring data...\n")
	database, err = transferData(database)
	if err != nil {
		log.Printf("Error while transferring data\n")
		return nil, err
	}
	log.Printf("Data transferred\n")
	return database, nil
}

func initDatabase(login string, password string, ip string, db_name string) (*sql.DB, error) {

	//New connection to mysql
	db, err := sql.Open("mysql", login+":"+password+"@tcp("+ip+")/")
	if err != nil {
		return nil, err
	}

	// make sure connection is available and a 5 seconds timeout
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	_, err = db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+db_name)
	if err != nil {
		db.Close()
		return nil, err
	}

	// Connect to the newly created database
	database, err := sql.Open("mysql", login+":"+password+"@tcp("+ip+")/"+db_name)
	if err != nil {
		database.Close()
		return nil, err
	}
	return database, nil
}

// Create the tables in the database
func createTables(database *sql.DB) error {

	queries := []string {
	`CREATE TABLE IF NOT EXISTS ChannelType (
		ChannelTypeID smallint UNSIGNED AUTO_INCREMENT NOT NULL,
		Name varchar(30) NOT NULL,
		PRIMARY KEY (ChannelTypeID)
	);`,

	`CREATE TABLE IF NOT EXISTS EventType (
		EventTypeID smallint UNSIGNED AUTO_INCREMENT NOT NULL,
		Name varchar(30) NOT NULL,
		PRIMARY KEY (EventTypeID)
	);`,
	
	`CREATE TABLE IF NOT EXISTS Content (
		ContentID int UNSIGNED AUTO_INCREMENT NOT NULL,
		ClientContentID bigint UNSIGNED NOT NULL,
		InsertDate timestamp NOT NULL,
		PRIMARY KEY (ContentID)
	);`,
	
	`CREATE TABLE IF NOT EXISTS Customer (
		CustomerID bigint UNSIGNED AUTO_INCREMENT NOT NULL,
		ClientCustomerID bigint UNSIGNED NOT NULL,
		InsertDate timestamp NOT NULL,
		PRIMARY KEY (CustomerID)
	);`,
	
	`CREATE TABLE IF NOT EXISTS CustomerData (
		CustomerChannelID bigint UNSIGNED AUTO_INCREMENT NOT NULL,
		CustomerID bigint UNSIGNED NOT NULL,
		ChannelTypeID smallint UNSIGNED NOT NULL,
		ChannelValue varchar(600) NOT NULL,
		InsertDate timestamp NOT NULL,
		PRIMARY KEY (CustomerChannelID),
		FOREIGN KEY (CustomerID) REFERENCES Customer (CustomerID),
		FOREIGN KEY (ChannelTypeID) REFERENCES ChannelType (ChannelTypeID)
	);`,
	
	`CREATE TABLE IF NOT EXISTS CustomerEvent (
		EventID bigint UNSIGNED AUTO_INCREMENT NOT NULL,
		ClientEventID bigint NOT NULL,
		InsertDate timestamp NOT NULL,
		PRIMARY KEY (EventID)
	);`,
	
	`CREATE TABLE IF NOT EXISTS CustomerEventData (
		EventDataId bigint UNSIGNED AUTO_INCREMENT NOT NULL,
		EventID bigint UNSIGNED NOT NULL,
		ContentID int UNSIGNED NOT NULL,
		CustomerID bigint UNSIGNED NOT NULL,
		EventTypeID smallint UNSIGNED NOT NULL,
		EventDate timestamp NOT NULL,
		Quantity smallint UNSIGNED NOT NULL,
		InsertDate timestamp NOT NULL,
		PRIMARY KEY (EventDataId),
		FOREIGN KEY (EventID) REFERENCES CustomerEvent (EventID),
		FOREIGN KEY (ContentID) REFERENCES Content (ContentID),
		FOREIGN KEY (CustomerID) REFERENCES Customer (CustomerID),
		FOREIGN KEY (EventTypeID) REFERENCES EventType (EventTypeID)
	);`,
	
	`CREATE TABLE IF NOT EXISTS ContentPrice (
		ContentPriceID mediumint UNSIGNED AUTO_INCREMENT NOT NULL,
		ContentID int UNSIGNED NOT NULL,
		Price decimal(8,2) UNSIGNED NOT NULL,
		Currency char(3) NOT NULL,
		InsertDate timestamp NOT NULL,
		PRIMARY KEY (ContentPriceID),
		FOREIGN KEY (ContentID) REFERENCES Content (ContentID)
	);`,
	}

	for _, query := range queries {
		_, err := database.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}

// Fill the tables with the data from the csv files + generative data
func transferData(database *sql.DB) (*sql.DB, error){
	log.Printf("Start Generative Data\n")
	err := generativeFill(database)
	if err != nil {
		log.Printf("Error %s when inserting generative data\n", err)
		return nil, err
	}
	log.Printf("Start Fill Content\n")
	err = fillContent(database)
	if err != nil {
		log.Printf("Error %s when Fill Content\n", err)
		return nil, err
	}
	return database, nil
}

func fillContent(database *sql.DB) error {
	tableName := []string{"Customer", "CustomerData", "CustomerEvent", "Content", "ContentPrice", "CustomerEventData"}
	csvPath := []string{"./csv/Customer.csv", "./csv/CustomerData.csv", "./csv/CustomerEvent.csv", "./csv/Content.csv", "./csv/ContentPrice.csv", "./csv/CustomerEventData.csv"}
	// csvPath := []string{"../csv/Customer.csv", "../csv/CustomerData.csv", "../csv/CustomerEvent.csv", "../csv/Content.csv", "../csv/ContentPrice.csv", "../csv/CustomerEventData.csv"}

	for i, table := range tableName {
		err := importDataFromCSV(database, csvPath[i], table)
		if err != nil {
			log.Printf("Error %s when importing %s\n", err, table)
			return err
		}
		displayProgressBar(i+1, len(tableName))
	}
	fmt.Printf("\n")
	return nil
}

// import data from csv file and insert into database
func importDataFromCSV(db *sql.DB, csvPath, tableName string) error {
	// log.Printf("Importing data from %s\n", csvPath)
	file, err := os.Open(csvPath)
	if err != nil {
		log.Printf("Error %s when opening file %s\n", err, csvPath)
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("Error %s when reading file %s\n", err, csvPath)
		return err
	}

	// Read the first row (header) from the CSV file
	header := records[0]

	// Store the column headers in the 'columnNames' slice
	columnNames := header

	// Create the SQL statement for the INSERT operation
	numColumns := len(columnNames)
	placeholders := make([]string, numColumns)
	for i := 0; i < numColumns; i++ {
		// log.Printf("columnNames[i]: %s\n", columnNames[i])
		placeholders[i] = "?"
	}

	// Insert the records in the database with querry
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, strings.Join(columnNames, ", "), strings.Join(strings.Split(strings.Repeat("?", len(columnNames)), ""), ", "))
	// log.Printf("query: %s\n", query)
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement\n", err)
		return err
	}
	defer stmt.Close()

	for i, record := range records {
		//skip header
		if i == 0 {
			continue
		}

		// Check if the record has the correct number of fields
		if len(record) != len(columnNames) {
			return fmt.Errorf("record at row %d has an incorrect number of columns", i+1)
		}

		// Convert the record values to interface values,
		// * make([]interface{}, len(record)) creates a slice of interfaces with the same length as the record slice
		values := make([]interface{}, len(record))
		for j, colValue := range record {
			values[j] = colValue
		}

		// Insert the record into the database,
		// * values... unpacks the values slice into separate arguments for the function call
		_, err := stmt.Exec(values...)
		if err != nil {
			log.Printf("Error %s when inserting %s\n", err, record)
			return err
		}
	}
	// log.Printf("Data imported from %s\n", csvPath)
	return nil
}


// Create generative random numbers, email, etc...
func generativeFill(database *sql.DB) error {

	channelTypes := []string{randomEmail(), randomPhoneNumber(), randomPostalCode(), randomMobileID(), randomCookie()}
	if err := InsertRowsInDb(database, channelTypes, "ChannelType"); err != nil {
		log.Printf("Error %s when inserting ChannelTypes\n", err)
		return err
	}
	eventTypes := []string{"sent", "view", "click", "visit", "cart", "purchase"}
	if err := InsertRowsInDb(database, eventTypes, "EventType"); err != nil {
		log.Printf("Error %s when inserting EventTypes\n", err)
		return err
	}
	return nil
}

// Insert string content inside the database
func InsertRowsInDb(database *sql.DB, rows []string, tableName string) error {
	for _, row := range rows {
		query := fmt.Sprintf("INSERT INTO %s (Name) VALUES (?)", tableName)
		_, err := database.Exec(query, row)
		if err != nil {
			log.Printf("Error %s when inserting %s in tableName:%s\n", err, row, tableName)
			return err
		}
	}
	return nil
}