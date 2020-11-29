package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
)

const (
	ledgerUrl = "http://host.docker.internal:9000"
	// ledgerUrl = "http://localhost:9000"
)

// Controller types implement basic connectivity
// to ACA-Py agents.
type Controller interface {
	Alias() string
	AgentUrl() string
	PublicDid() (string, error)
	SetPublicDid(string)
	ConnectionDid() string
}

// util that generates seeds for did registration with ledger
func Seed() string {
	seed := "my_seed_000000000000000000000000"
	randInt := rand.Intn(800000) + 100000
	seed = seed + strconv.Itoa(randInt)
	return seed[len(seed)-32:]
}

type RegisterDidRequest struct {
	Alias string `json:"alias"`
	Seed  string `json:"seed"`
	Role  string `json:"role"`
}

type RegisterDidResponse struct {
	Did string `json:"did"`
	Seed string `json:"seed"`
}

// Register agent with ledger and receive a DID
func RegisterDidWithLedger(controller Controller, seed string) (string, error) {

	if did, err := controller.PublicDid(); did != "" {
		return "", fmt.Errorf("agent already registered public DID on ledger: %v\n", err)
	} else if err != nil {
		return "", fmt.Errorf("failed while checking if public DID exists: %v\n", err)
	}

	reqBody := RegisterDidRequest{
		Alias: controller.Alias(),
		Seed:  seed,
		Role:  "TRUST_ANCHOR",
	}

	resp, err := SendRequest_POST(ledgerUrl, "/register", reqBody)
	if err != nil {
		return "", fmt.Errorf("Failed to send post request: %v\n", err)
	}
	defer resp.Body.Close()

	var didResp RegisterDidResponse
	err = json.NewDecoder(resp.Body).Decode(&didResp)
	if err != nil {
		return "", fmt.Errorf("Failed to unmarshal json: %v\n", err)
	}
	didResp.Seed = seed

	// store in wallet
	_, err = SendRequest_POST(controller.AgentUrl(), "/connections/store-public-did", didResp)
	if err != nil {
		return "", fmt.Errorf("Failed to request store_public_did: %v\n", err)
	}

	controller.SetPublicDid(didResp.Did)
	return didResp.Did, nil
}

type PutKeyToLedgerRequest struct {
	SigningDid string `json:"signing_did"`
	SigningVk string `json:"signing_vk"`
}

type PutKeyToLedgerResponse struct {
	Status string `json:"status"`
}

// put signing DID and verification key to ledger
func PutKeyToLedger(controller Controller, did, vk string) error {


	payload := PutKeyToLedgerRequest{
		SigningDid: did,
		SigningVk:  vk,
	}

	resp, err := SendRequest_POST(controller.AgentUrl(), "/connections/put-key-ledger", payload)
	if err != nil {
		return fmt.Errorf("Error occurred while sending post request: %v\n", err)
	}
	defer resp.Body.Close()

	var response PutKeyToLedgerResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return fmt.Errorf("error occurred while decoding json: %v\n", err)
	}

	if response.Status != "true" {
		return errors.New("put key to ledger failed")
	}

	return nil
}
