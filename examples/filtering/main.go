package main

import (
	"fmt"
	"log"
	"os"

	"github.com/treethought/cam-sweeney-sdk/sdk"
)

// NOTE: be sure to set your ONE_API_KEY env var
var API_KEY = os.Getenv("ONE_API_KEY")

func main() {

	// The following functionality does not require authenticated
	// so we cna use NewUnAuthenticated

	// let's set pagination limit to be 2000 across all API calls
	// so we get everything we want

	opts := []sdk.RequestOption{sdk.WithLimit(1000)}
	config := sdk.ClientConfig{PersistentOptions: opts, ApiKey: API_KEY}
	client := sdk.NewWithConfig(config)

	// lets get all human characters
	rOpts := []sdk.RequestOption{sdk.WithFilterInclude("race", "Human")}

	// lets also sort by realm
	rOpts = append(opts, sdk.WithSort("realm", "asc"))

	humans, err := client.Characters().List(rOpts...)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Human Characters:")
	for _, h := range humans {
		fmt.Printf("Name: %s, Realm: %s, Race: %s\n", h.Name, h.Realm, h.Race)
	}

}
