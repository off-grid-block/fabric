package indyverify

import (
	"bytes"
	"crypto/sha256"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type IndyResponse struct {
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
		return false, errors.New("DID size not equal to 22")
	}
	Signature := string(SignatureBytes)
	if len(Signature) != 140 {
		return false, errors.New("Signature size not equal to 140")
	}

	//Create Payload
	ProposalHash := sha256.Sum256(ProposalBytes)
	EncodedHash := b64.StdEncoding.EncodeToString(ProposalHash[:])
	type Payload struct {
		Message   string `json:"message"`
		Did       string `json:"their_did"`
		Signature string `json:"signature"`
	}
	P := &Payload{Message: EncodedHash, Did: DidValue, Signature: Signature}
	PayloadBytes, err := json.Marshal(P)
	if err != nil {
		return false, errors.New("Error creating Payload")
	}
	PayloadBytesString := string(PayloadBytes)

	//Verify Signature
	VerifyUrl := "http://10.53.17.40:8003/verify_signature"
	Request, _ := http.NewRequest("POST", VerifyUrl, bytes.NewBuffer([]byte(PayloadBytesString)))
	Request.Header.Add("content-type", "text/plain")
	Response, err := http.DefaultClient.Do(Request)
	if Response.StatusCode != 200 {
		return false, errors.New("Error connecting to Indy server")
	}
	defer Response.Body.Close()

	//Validate Response
	ResponseBody, _ := ioutil.ReadAll(Response.Body)
	var Result map[string]interface{}
	json.NewDecoder(Response.Body).Decode(&Result)
	ResponseJson := IndyResponse{}
	err = json.Unmarshal(ResponseBody, &ResponseJson)
	if err != nil {
		return false, errors.New("error unmarshaling response from Indy")
	}
	if ResponseJson.Status != "Signature verified" {
		return false, errors.New(ResponseJson.Status)
	}
	return true, nil
}
