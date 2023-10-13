package main


const (
	mobileIDSize = 10 // Change this to the desired length of the MobileID
)

// var (
// 	randomSource = rand.NewSource(time.Now().UnixNano())
// 	random       = rand.New(randomSource)
// )

func randomMobileID() string {
	mobileID := ""

	// Generate the MobileID digits
	for i := 0; i < mobileIDSize; i++ {
		mobileID += string(random.Intn(10) + '0')
	}

	return mobileID
}
