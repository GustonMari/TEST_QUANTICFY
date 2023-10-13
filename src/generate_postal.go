package main

const (
	postalSize = 5 // Change this to the desired length of the postal code
)

// var (
// 	randomSource = rand.NewSource(time.Now().UnixNano())
// 	random       = rand.New(randomSource)
// )

func randomPostalCode() string {
	postal := ""

	// Generate the postal code digits
	for i := 0; i < postalSize; i++ {
		postal += string(random.Intn(10) + '0')
	}

	return postal
}
