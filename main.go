package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os/exec"
)

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

