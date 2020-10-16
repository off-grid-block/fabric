/*
Copyright TCS All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package indyverify

import (
	"bytes"
	"crypto/sha256"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"fmt"
)

type indyResponse struct {
	Status       string `json:"status"`
	ConnectionID string `json:"connection_id"`
}

// Indyverify - this function receives Proposal, DID and Signature as bytes.
func Indyverify(ProposalBytes []byte, DidBytes []byte, SignatureBytes []byte) (Status bool, err error) {

	//Validate Inputs
	if len(ProposalBytes) == 0 {
		return false, errors.New("Empty proposal bytes received while verifying Indy signature")
	}
	if len(DidBytes) == 0 {
		return false, errors.New("Empty DID received while verifying Indy signature")
	}
	if len(SignatureBytes) == 0 {
		return false, errors.New("Empty signature received while verifying Indy signature")
	}
	DidValue := string(DidBytes)
	if len(DidValue) != 22 {
		return false, errors.New("DID size not equal to 22 (the DID is): " + DidValue)
	}
	Signature := string(SignatureBytes)

	fmt.Println("inside Indyverify")

	fmt.Printf("ProposalBytes: %v\n", ProposalBytes)
	fmt.Printf("DidBytes: %v\n", DidBytes)
	fmt.Printf("Signature: %v\n", SignatureBytes)

	//Create Payload
	ProposalHash := sha256.Sum256(ProposalBytes)
	EncodedHash := b64.StdEncoding.EncodeToString(ProposalHash[:])
	type Payload struct {
		Message   string `json:"message"`
		Did       string `json:"their_did"`
		Signature string `json:"signature"`
	}
	P := &Payload{Message: EncodedHash, Did: DidValue, Signature: Signature}

	fmt.Printf("ProposalBytes (Payload): %s\n", P.Message)
	fmt.Printf("DidBytes (Payload): %s\n", P.Did)
	fmt.Printf("Signature (Payload): %s\n", P.Signature)

	PayloadBytes, err := json.Marshal(P)
	if err != nil {
		return false, errors.New("Error creating Payload")
	}
	PayloadBytesString := string(PayloadBytes)

	//Verify Signature
	VerifyURL := "http://ci_msp.example.com:7997/verify_signature"
	Request, _ := http.NewRequest("POST", VerifyURL, bytes.NewBuffer([]byte(PayloadBytesString)))
	Request.Header.Add("content-type", "text/plain")
	Response, err := http.DefaultClient.Do(Request)
	if err != nil {
		return false, fmt.Errorf("Error sending response to Indy server: %v", err)
	}
	if Response == nil {
		return false, errors.New("No response from Indy server")
	}

	defer Response.Body.Close()

	//Validate Response
	ResponseBody, _ := ioutil.ReadAll(Response.Body)
	var Result map[string]interface{}
	json.NewDecoder(Response.Body).Decode(&Result)
	ResponseJSON := indyResponse{}
	err = json.Unmarshal(ResponseBody, &ResponseJSON)
	if err != nil {
		return false, errors.New("error unmarshaling response from Indy")
	}
	if ResponseJSON.Status != "Signature verified" {
		return false, errors.New("Response from Indy:" + ResponseJSON.Status)
	}
	return true, nil
}
