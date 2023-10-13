package main

import (
	// "database/sql"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

// load .env file
func loadEnv(key string) string {

	// check if .env file exists
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error reading .env file : %v\n", err)
	}
	//return the value of the key
	return os.Getenv(key)
}

func displayProgressBar(current, total int) {
	const progressBarWidth = 40
	progress := float64(current) / float64(total)
	barLength := int(progress * progressBarWidth)

	fmt.Print("\r[")
	for i := 0; i < barLength; i++ {
		fmt.Print("=")
	}
	for i := barLength; i < progressBarWidth; i++ {
		fmt.Print(" ")
	}
	fmt.Printf("] %3d%%", int(progress*100))
}


func printQuantileInfo(Quantiles []float64, InfoQuantiles map[float64]top_customer){
	i := 1
	color.HiGreen("\nQUANTILES")
	for _, value := range Quantiles {
		color.Blue("Quantile: %d\n", i)
		color.Green("%v\n", value)
		i++
	}
	i = 1
	fmt.Print("\n==================================================\n\n")
	color.HiGreen("CUSTOMER INFO\n")

	// for _, value := range InfoQuantiles {
	// 	color.Blue("Quantile: %d\n", i)
	// 	color.HiYellow("%v\n", value)
	// 	i++
	// }
	// Create a slice to store the keys (quantile values)
	keys := make([]float64, 0, len(InfoQuantiles))
	for key := range InfoQuantiles {
		keys = append(keys, key)
	}

	// Sort the keys
	sort.Float64s(keys)

	// Iterate over the sorted keys and access values from the map in order
	for _, key := range keys {
		value := InfoQuantiles[key]
		color.Blue("Quantile: %d\n", i)
		color.Green("%+v\n", value)
		i++
	}
}