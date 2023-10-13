package main


const (
	phoneSize = 10 // Change this to the desired length of the phone number
)

// var (
// 	randomSource = rand.NewSource(time.Now().UnixNano())
// 	random       = rand.New(randomSource)
// )

func randomPhoneNumber() string {
	phone := "+336"

	// Generate area code (XXX)
	for i := 0; i < 8; i++ {
		phone += string(random.Intn(10) + '0')
	}
	return phone
}
