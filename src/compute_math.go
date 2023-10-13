package main

import (
	"log"
	"sort"
	"github.com/fatih/color"
	// "github.com/tdewolff/quantile"
)

type top_customer struct {
	NumberOfCustomer int
	MaxPurchase float64
	MinPurchase float64
}

// Take my struct and return a sorted slice, ascending order
func sortTotalPurchase(TotalPurchase map[int]TotalPrice) []TotalPrice {
	// Create a slice from the map
	totalPurchaseSlice := make([]TotalPrice, 0, len(TotalPurchase))
	for _, totalPrice := range TotalPurchase {
		totalPurchaseSlice = append(totalPurchaseSlice, totalPrice)
	}

	// Define a sorting function based on TotalPrice
	sort.Slice(totalPurchaseSlice, func(i, j int) bool {
		return totalPurchaseSlice[i].TotalPrice < totalPurchaseSlice[j].TotalPrice
	})

	return totalPurchaseSlice
}

// Take a sorted slice and return a map of the 2.5% best customers
func mapBestCustomer(TotalPurchase []TotalPrice) map[int]float64 {
	// Create a map from the slice
	bestCustomer := make(map[int]float64)
	// calculate the index of the 2.5% best customers (first quantile)
	index := int(float64(len(TotalPurchase)) * (1 - 0.025))

	for i := index; i < len(TotalPurchase); i++ {
		bestCustomer[TotalPurchase[i].CustomerID] = TotalPurchase[i].TotalPrice
	}

	if len(bestCustomer) == 0 {
		log.Fatal("No best customer found")
	} else {
		color.HiGreen("\nBEST CUSTOMER MAP")
		log.Printf("%v\n\n", bestCustomer)
	}
	
	return bestCustomer
}

// convert slice of TotalPrice to slice of float64
func TotalPriceTofloat64(TotalPurchase []TotalPrice) ([]float64) {
	var data []float64
	for _, value := range TotalPurchase {
		data = append(data, value.TotalPrice)
	}
	return data
}

// calculate all the quantiles and return a slice of them + a map of the quantiles info
func calculateAllQuantiles(TotalPurchase []TotalPrice, quantileCount int) ([]float64, map[float64]top_customer, error) {
	
	data := TotalPriceTofloat64(TotalPurchase)

	// Sort the data
	sort.Float64s(data)

	// Create slice to store quantiles + additional information
	quantiles := make([]float64, quantileCount)
	// customersInfo := make(map[float64]top_customer, quantileCount)

	// Iterate through all the quantiles and compute them
	for i := 0; i < quantileCount; i++ {
		quantileValue := calculateOneQuantile(data, i, quantileCount)
		quantiles[i] = quantileValue

		// Calculate additional information for this quantile
		// customersInfo[float64(i)] = calculateCustomerInfo(data, quantileValue)
	}
    customersInfo := calculateCustomerInfo(data, 40)

	return quantiles, customersInfo, nil
}

func calculateOneQuantile(data []float64, k, n int) float64 {
	// Calculate the position of the desired quantile
	position := float64(k) / float64(n-1)
	var upperValue float64
	// Calculate the index and fractional part
	index := position * float64(len(data)-1)
	lowerIndex := int(index)
	fraction := index - float64(lowerIndex)

	// Calculate the quantile value by linear interpolation
	lowerValue := data[lowerIndex]
	if (lowerIndex + 1) >= len(data) {
		upperValue = data[lowerIndex]
	}else {
		upperValue = data[lowerIndex+1]
	}
	quantileValue := lowerValue + fraction*(upperValue-lowerValue)

	return quantileValue
}

// Print the quantiles and their information
func calculateCustomerInfo(data []float64, quantileCount int) map[float64]top_customer {
	// Create a map to store quantile information
	customerInfo := make(map[float64]top_customer)
	// Calculate the number of clients per quantile
	clientsPerQuantile := len(data) / quantileCount
	// Create a slice to store quantiles for sorting
	quantiles := make([]float64, quantileCount)

	for i := 0; i < quantileCount; i++ {
		// Calculate the index range for the current quantile
		startIndex := i * clientsPerQuantile
		endIndex := (i + 1) * clientsPerQuantile
		if i == quantileCount-1 {
			endIndex = len(data)
		}
		// Initialize minPurchase and maxRevenue with the first value in the current quantile
		minPurchase := data[startIndex]
		maxPurchase := data[startIndex]
		customersInQuantile := endIndex - startIndex
		// Calculate max and min revenue in the current quantile
		for j := startIndex; j < endIndex; j++ {
			if data[j] > maxPurchase {
				maxPurchase = data[j]
			}
			if data[j] < minPurchase {
				minPurchase = data[j]
			}
		}
		// Calculate the quantile value
		quantile := float64(i+1) / float64(quantileCount)
		quantiles[i] = quantile
		// Store quantile information in the map
		customerInfo[quantile] = top_customer{
			NumberOfCustomer: customersInQuantile,
			MinPurchase:      minPurchase,
			MaxPurchase:      maxPurchase,
		}
	}
	return customerInfo
}