package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func createVaultSnapshot() {
	caCert := ReadEnvOrFile("VAULT_CACERT")
	vaultToken := RequireEnvOrFile("VAULT_TOKEN")
	vaultAddr := RequireEnvOrFile("VAULT_ADDR")

	client := CreateVaultClient(string(vaultAddr), string(vaultToken), caCert)

	snapshotName := "vault.snapshot"
	err := client.CreateSnapshot(snapshotName)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Wrote vault snapshot to %s\n", snapshotName)
}

func RunRestic(arg ...string) (string, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("restic", arg...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", errors.New(fmt.Sprintf("%s: %s", err, stderr.String()))
	}

	return stdout.String(), nil
}

func main() {
	// createVaultSnapshot()

	// Initialize restic repo if required
	_, err := RunRestic("cat", "config")
	if err != nil {
		if strings.Index(err.Error(), "unable to open config file") != -1 {
			_, err := RunRestic("init")
			if err != nil {
				log.Fatalf("Could not initialize restic repository: %s", err)
			}

		}
		log.Println("Initialized restic repository")
	} else {
		log.Println("Found existing restic repository")
	}
}

