package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gomods/athens/cmd/gitauth/credentials"
)

func main() {
	secretsFile := flag.String("file", "", "the path of the json file holding the credentials")
	flag.Parse()

	if *secretsFile == "" {
		log.Fatal("-file parameter not found")
	}

	params, err := credentials.ParseInput(os.Stdin)
	if err != nil {
		log.Fatalf("Failed to read input %v", err)
	}

	f, err := os.Open(*secretsFile)
	if err != nil {
		log.Fatalf("Failed to open file %v", err)
	}

	credentials, err := credentials.FromJSON(f, params)
	if err != nil {
		log.Fatalf("Failed to read credentials %v", err)
	}

	fmt.Printf("username=%s\n", credentials.Username)
	fmt.Printf("password=%s\n", credentials.Password)
	os.Exit(0)
}
