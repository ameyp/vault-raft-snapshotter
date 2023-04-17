package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Vault struct {
	client http.Client
	token string
	addr string
}

type AuthDetails struct {
	ClientToken string `json:"client_token"`
}

type LoginRequest struct {
	Role string `json:"role"`
	Jwt string `json:"jwt"`
}

type LoginResponse struct {
	Auth AuthDetails
}

func createVaultClient(caCert []byte) http.Client {
	if caCert == nil {
		return http.Client{}
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// https client tls config
	// InsecureSkipVerify true means not validate server certificate (so no need to set RootCAs)
	tlsConfig := &tls.Config{
		RootCAs:            caCertPool,
		InsecureSkipVerify: false,
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}

	return http.Client{Transport: transport}
}

func CreateVaultClient(addr string, token string, caCert []byte) Vault {
	return Vault{
		client: createVaultClient(caCert),
		token: token,
		addr: addr,
	}
}

func (vault *Vault) Login(role string, jwtToken string) error {
	if vault.token == "" {
		log.Printf("Logging in with role %s", role)
		url := fmt.Sprintf("%s/v1/auth/kubernetes/login", vault.addr)
		data, err := json.Marshal(LoginRequest{Role: role, Jwt: jwtToken})
		if err != nil {
			return fmt.Errorf("Could not marshall LoginRequest: %w", err)
		}
		var dataReader io.Reader = bytes.NewReader(data)

		req, err := http.NewRequest("POST", url, dataReader)

		// read response
		resp, err := vault.client.Do(req)
		if err != nil {
			return fmt.Errorf("Login request failed: %w", err)
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("Could not read response of login: %w", err)
		}

		if resp.StatusCode != 200 && resp.StatusCode != 204 {
			return errors.New(fmt.Sprintf("GET response status code: %v, body: %s", resp.StatusCode, string(body)))
		}

		var response LoginResponse
		json.Unmarshal(body, &response)

		log.Println("Logged in successfully")
		vault.token = response.Auth.ClientToken
	}

	return nil
}

func (vault Vault) CreateSnapshot(filename string) error {
	// https client request
	url := fmt.Sprintf("%s/v1/sys/storage/raft/snapshot", vault.addr)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Vault-Token", vault.token)

	// read response
	resp, err := vault.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == 200 {
		outFile, err := os.Create(filename)
		if err != nil {
			return err
		}

		defer outFile.Close()

		_, err = io.Copy(outFile, resp.Body)
		if err != nil {
			return err
		}
	} else {
		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		log.Fatalf("Could not create vault snapshot. Status: %s. Body: %s", resp.Status, string(contents))
	}

	return nil
}
