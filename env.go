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

func ReadEnvOrFile(variable string) []byte {
	val, _ := readEnvVariableOrFile(variable, false)
	return val
}

func RequireEnvOrFile(variable string) []byte {
	val, err := readEnvVariableOrFile(variable, true)
	if err != nil {
		log.Fatal(err)
	}

	return val
}

