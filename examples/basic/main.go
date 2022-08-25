package main

import (
	"fmt"
	"log"
	"os"

	"github.com/treethought/cam-sweeney-sdk/sdk"
)

var API_KEY = os.Getenv("ONE_API_KEY")

func main() {

	// The following functionality does not require authenticated
	// so we cna use NewUnAuthenticated

	client := sdk.NewUnAuthenticated()

	// get all books
	books, err := client.Books().List()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("LOTR Books:")
	for _, b := range books {
		fmt.Printf("Title: %s, ID: %s\n", b.Name, b.ID)
	}

	// lets get the chapters of the two towers, limited to 3 items
	twoTowersID := "5cf58077b53e011a64671583"
	chapters, err := client.Books().GetChapters(twoTowersID, sdk.WithLimit(3))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nTwo Towers Chapters: (limited to 3")
	for _, c := range chapters {
		fmt.Printf("Title: %s, ID: %s\n", c.Name, c.ID)
	}

	// we can also sort them by a field
	sortOpt := sdk.WithSort("name", "asc")
	books, err = client.Books().List(sortOpt)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nLOTR Books, sorted:")
	for _, b := range books {
		fmt.Printf("Title: %s, ID: %s\n", b.Name, b.ID)
	}

}
