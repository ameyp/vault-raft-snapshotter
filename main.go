package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func readEnvVariableOrFile(variable string, required bool) ([]byte, error) {
	var valueContent []byte
	value := os.Getenv(variable)
	if value != "" {
		if string(value[0]) == "@" {
			filename := value[1:]
			log.Printf("Reading %s from file: %s", variable, filename)
			fileContent, err := ioutil.ReadFile(filename)
			if err != nil {
				return nil, err
			}
			valueContent = fileContent
		} else {
			valueContent = []byte(value)
		}
	} else if required {
		return nil, errors.New(fmt.Sprintf("%s environment variable must be set.", variable))
	}

	return valueContent, nil
}

func main() {
	caCert, err := readEnvVariableOrFile("VAULT_CACERT", false)
	if err != nil {
		log.Fatal(err)
	}

	vaultToken, err := readEnvVariableOrFile("VAULT_TOKEN", true)
	if err != nil {
		log.Fatal(err)
	}

	client := CreateVaultClient("https://vault.wirywolf.com", string(vaultToken), caCert)

	// now := time.Now()
	// datetime := now.Format(time.RFC3339)
	// snapshotName := fmt.Sprintf("%s.snap", datetime)

	snapshotName := "vault.snapshot"
	err = client.CreateSnapshot(snapshotName)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Wrote vault snapshot to %s\n", snapshotName)
}

