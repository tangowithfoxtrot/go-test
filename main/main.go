package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	bw "github.com/tangowithfoxtrot/go-module-test"
)

/*
#cgo linux LDFLAGS: -L/usr/local/lib -L/usr/lib -L ../lib/linux-amd64
#cgo darwin LDFLAGS: -L/usr/local/lib -L/usr/lib -L ../lib/macos-arm64
*/
import "C"

func main() {
	apiURL := "http://localhost:4000"
	identityURL := "http://localhost:33656"

	organizationIDStr := os.Getenv("ORGANIZATION_ID")

	bitwardenClient, _ := bw.NewBitwardenClient(&apiURL, &identityURL)

	accessToken := os.Getenv("ACCESS_TOKEN")
	if accessToken == "" {
		accessToken = os.Getenv("BWS_ACCESS_TOKEN")
	}
	projectName := os.Getenv("PROJECT_NAME")
	if projectName == "" {
		projectName = "NewTestProject" // default value
	}

	err := bitwardenClient.AccessTokenLogin(accessToken, nil)
	if err != nil {
		panic(err)
	}

	secretIdentifiers, err := bitwardenClient.Secrets.List(organizationIDStr)
	if err != nil {
		panic(err)
	}

	secretIDs := make([]string, len(secretIdentifiers.Data))
	for i, identifier := range secretIdentifiers.Data {
		secretIDs[i] = identifier.ID
	}

	secrets, err := bitwardenClient.Secrets.GetByIDS(secretIDs)
	if err != nil {
		log.Fatalf("Error getting secrets: %v", err)
	}

	jsonSecrets, err := json.MarshalIndent(secrets, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling secrets to JSON: %v", err)
	}

	fmt.Println(string(jsonSecrets))
}
