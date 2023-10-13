package main

import (
	"database/sql"
	// "fmt"
	"log"

	// "sort"

	// "time"
	// "fmt"
	"github.com/fatih/color"

	_ "github.com/go-sql-driver/mysql"
)

func main(){
	var login string = loadEnv("LOGIN")
	var password string = loadEnv("PASSWORD")
	var ip string = loadEnv("IP")
	var db_name string = loadEnv("DB_NAME")

	var database *sql.DB

	// ? =========================================CONNECT==========================================

	database, err := connect(login, password, ip, db_name)
	if err != nil {
		log.Fatal(err)
	}
	
	// ? =========================================TREAT===========================================
	
	// Create a slice of the total price of each purchase per customer
	TotalPurchase, err := sumAllPurchase(database)
	if err != nil {
		log.Fatal(err)
	}

	//Sort the slice by total price
	sorted := sortTotalPurchase(TotalPurchase)
	color.HiGreen("SORTED LIST OF TOTAL PURCHASE")
	for _, value := range sorted {
		color.Blue("%f\n", value.TotalPrice)
	}
	// map for the best customers
	mapBestCustomer(sorted)

	// Calculate the quantiles + info
	Quantiles, InfoQuantiles, err := calculateAllQuantiles(sorted, 40)
	if err != nil {
		log.Fatal(err)
	}
	printQuantileInfo(Quantiles, InfoQuantiles)
	
	// ? =========================================EXPORT==========================================

	exportToNewTable(database, TotalPurchase)
	
	//!ATTTENTION its just for test purpose
	// dropAllTables(database)
	
	
	defer database.Close()
}