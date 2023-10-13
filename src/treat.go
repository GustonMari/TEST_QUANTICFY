package main

import (
	"log"
	"fmt"
	"database/sql"
	// "time"
	"github.com/fatih/color"

)

type PriceInfo struct {
    Price    float64
    Currency string
}

type TotalPrice struct {
	CustomerID int
	TotalPrice float64
}


// Return ContentID of all content purchased after 2020-04-01 00:00:00 +0000 UTC
func getContentIDByPurchase(db *sql.DB) ([]int, error) {
	var contentID []int
	var err error

	query := fmt.Sprintf("SELECT ContentID FROM CustomerEventData WHERE EventTypeID = 6 AND InsertDate >= '2020-04-01 00:00:00 +0000 UTC'")
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error while querying database: %s\n", query)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			log.Printf("Error while scanning rows\n")
			return nil, err
		}
		contentID = append(contentID, id)
	}
	return contentID, nil
}

func convertToEuro(price float64, currency string) float64 {
	if currency == "EUR" {
		return price
	}

	ExchangeRate := map[string]float64{
		"EUR": 1.0,
		"USD": 1.06,
		"GBP": 0.86,
		"BRL": 5.33,
		"XAF": 655.95,
		"XOF": 655.95,
		"IDR": 17200.0,
		"UZS": 10200.0,
		"UAH": 30.0,
		"PHP": 55.0,
		"CNY": 7.5,
		"AFN": 80.0,
		"PAB": 1.0,
		"ALL": 120.0,
		"MKD": 61.0,
		"RUB": 80.0,
		"ISK": 150.0,
		"GEL": 3.5,
		"PLN": 4.5,
		"VND": 25000.0,
		"SEK": 10.0,
		"KMF": 500.0,
		"HTG": 100.0,
		"XPF": 120.0,
		"KRW": 1200.0,
		"THB": 35.0,
		"COP": 4000.0,
		"PEN": 3.5,
	}
	if currency == "" {
		color.Red("Errror currency is empty\n")
	}
	return price / ExchangeRate[currency]
}

func getPriceInEuro(db *sql.DB, contentID int) (map[int]PriceInfo, error) {
    result := make(map[int]PriceInfo)
    query := fmt.Sprintf("SELECT Price, Currency FROM ContentPrice WHERE ContentID = %d", contentID)
    rows, err := db.Query(query)
    if err != nil {
        log.Printf("Error while querying database: %s\n", query)
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var price float64
        var currency string
        err := rows.Scan(&price, &currency)
        if err != nil {
            log.Printf("Error while scanning rows\n")
            return nil, err
        }
        result[contentID] = PriceInfo{
            Price:    convertToEuro(price, currency),
            Currency: "EUR",
        }
    }

    return result, nil
}

func sumAllPurchase(db *sql.DB) (map[int]TotalPrice, error) {
	TotalPurchase := make(map[int]TotalPrice)
	var customerID int
	var quantity int

	contentIDs, err := getContentIDByPurchase(db)
	if err != nil {
		log.Printf("Error while getting contentID by purchase\n")
		return nil, err
	}
	for _, contentID := range contentIDs {
		// fetch wich customer bought this content and the quantity
		query := fmt.Sprintf("SELECT CustomerID, Quantity FROM CustomerEventData WHERE ContentID = %d AND EventTypeID = 6 AND InsertDate >= '2020-04-01 00:00:00 +0000 UTC'", contentID)
		rows, err := db.Query(query)
		if err != nil {
			log.Printf("Error while querying database: %s\n", query)
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			if err := rows.Scan(&customerID, &quantity); err != nil {
				log.Printf("Error while scanning rows\n")
				return nil, err
			}
			// fmt.Printf("CustomerID: %d, Quantity: %d\n", customerID, quantity)
			mapPrice, err := getPriceInEuro(db, contentID)
			if err != nil {
				log.Printf("Error while getting price in euro\n")
				return nil, err
			}

			for _, priceInfo := range mapPrice {
				// fmt.Printf("Price: %f, Currency: %s\n", priceInfo.Price, priceInfo.Currency)
				purchaseAmount := priceInfo.Price * float64(quantity)
				// Update the TotalPurchase map for the customer
				total, exists := TotalPurchase[customerID]
				if !exists {
					total = TotalPrice{CustomerID: customerID, TotalPrice: 0}
				}
				total.TotalPrice += purchaseAmount
				TotalPurchase[customerID] = total
			}
		}
	}
	return TotalPurchase, nil
}