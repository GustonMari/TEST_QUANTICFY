package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"
)

// Just utils function to drop all tables
func dropAllTables(db *sql.DB) error {
	color.HiGreen("Dropping all tables...")
	// Disable foreign key checks
	_, err := db.Exec("SET FOREIGN_KEY_CHECKS = 0;")
	if err != nil {
		return err
	}

	tables := []string{
		"CustomerEventData",
		"ContentPrice",
		"CustomerData",
		"CustomerEvent",
		"Customer",
		"ChannelType",
		"EventType",
		"Content",
		"test_export" + strings.Replace(time.Now().Format("2006-01-02"), "-", "_", -1),
	}

	// Drop the tables in reverse order to avoid foreign key constraints
	for i := len(tables) - 1; i >= 0; i-- {
		query := fmt.Sprintf("DROP TABLE IF EXISTS %s;", tables[i])
		_, err := db.Exec(query)
		if err != nil {
			return err
		}
	}

	// Re-enable foreign key checks
	_, err = db.Exec("SET FOREIGN_KEY_CHECKS = 1;")
	if err != nil {
		return err
	}
	color.HiGreen("All tables dropped")
	return nil
}


func loadEnv(key string) string {

	// check if .env file exists
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error reading .env file : %v\n", err)
	}
	//return the value of the key
	return os.Getenv(key)
}

func main() {
	var login string = loadEnv("LOGIN")
	var password string = loadEnv("PASSWORD")
	var ip string = loadEnv("IP")
	var db_name string = loadEnv("DB_NAME")

	var database *sql.DB

	database, err := sql.Open("mysql", login+":"+password+"@tcp("+ip+")/"+db_name)
	if err != nil {
		log.Fatal("Error while connecting to database\n")
	}
	dropAllTables(database)
	defer database.Close()
}
