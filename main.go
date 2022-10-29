package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
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
	var destDir string
	if len(os.Args) > 1 {
		destDir = os.Args[1]
	} else {
		destDir = "."
	}
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	caCert := ReadEnvOrFile("VAULT_CACERT")
	vaultToken := RequireEnvOrFile("VAULT_TOKEN")
	vaultAddr := RequireEnvOrFile("VAULT_ADDR")

	client := CreateVaultClient(string(vaultAddr), string(vaultToken), caCert)

	snapshotPath := path.Join(destDir, "vault.snapshot")
	err := client.CreateSnapshot(snapshotPath)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Wrote vault snapshot to %s\n", snapshotPath)
}

