package main

import (
	"database/sql"
	"fmt"
	// "image/color"
	"log"
	"time"
	"github.com/fatih/color"
	"strings"
)

// Get Email by CustomerID
func getEmailByCustomerID(db *sql.DB, customerID int) (string, error) {
	query := `
		SELECT ct.Name
		FROM CustomerData cd
		JOIN ChannelType ct ON cd.ChannelTypeID = ct.ChannelTypeID
		WHERE cd.CustomerID = ? AND cd.ChannelTypeID = 1
	`

	var name string

	err := db.QueryRow(query, customerID).Scan(&name)
	if err != nil {
		return "", err
	}

	return name, nil
}

// Export the result of the query to a new table
func exportToNewTable(db *sql.DB, TotalPurchase map[int]TotalPrice) error {
	
	tableName := "test_export_" + strings.Replace(time.Now().Format("2006-01-02"), "-", "_", -1)
	createTableQuery := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		CustomerID INT,
		Email VARCHAR(255),
		CA DECIMAL(8,2)
	);`, tableName)

	// Create the new table
	_, err := db.Exec(createTableQuery)
	if err != nil {
		color.Red("Error while creating table %s\n", tableName)
		log.Fatal(err)
	}

	// Insert the result of the query into the new table
	for _, value := range TotalPurchase {
		Email, err := getEmailByCustomerID(db, value.CustomerID)
		if err != nil {
			Email = "test@test.com"
		}
		query := fmt.Sprintf("INSERT INTO %s (CustomerID, Email, CA) VALUES (%d, '%s', %f) ON DUPLICATE KEY UPDATE CA = %f", tableName, value.CustomerID, Email, value.TotalPrice, value.TotalPrice)
		_, err = db.Exec(query)
		if err != nil {
			color.Red("Error while inserting into table %s\n", tableName)
			log.Fatal(err)
		}
	}

	return nil
}