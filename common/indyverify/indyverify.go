package indyverify

import (
	"bytes"
	b64 "encoding/base64"
	"errors"

	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type verifySignatureOut struct {
	Status        string `json:"status"`
	Connection_id string `json:"connection_id"`
}

func Indyverify(ProposalBytes []byte, did []byte, signature []byte) (status bool, err error, id string) {

	hash := sha256.Sum256(ProposalBytes)
	encoded := b64.StdEncoding.EncodeToString(hash[:])
	fmt.Println()
	fmt.Println()
	fmt.Println("calculated hash", hash)
	fmt.Println("encoded hash", encoded)
	type Payload struct {
		Message   string `json:"message"`
		Did       string `json:"their_did"`
		Signature string `json:"signature"`
	}
	url := "http://10.53.17.40:8003/verify_signature"
	payload := Payload{Message: encoded, Did: string(did), Signature: string(signature)}
	payloadbytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("prepared payload", payload)
	payloadbytesstring := string(payloadbytes)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payloadbytesstring)))

	req.Header.Add("content-type", "text/plain")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Received error from Indy server", err)
		return false, errors.New("Error connecting to Indy server"), ""
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var result map[string]interface{}
	json.NewDecoder(res.Body).Decode(&result)

	fmt.Println(string(body))

	respJson := verifySignatureOut{}
	err = json.Unmarshal(body, &respJson)
	if err != nil {
		fmt.Println("error unmarshaling response from Indy", err)
		return false, errors.New("error unmarshaling response from Indy"), ""
	}

	if respJson.Status != "Signature verified" {
		return false, errors.New(respJson.Status), ""
	}
	// Verify proof code
	/*
		attrib := "app_name,app_id"
		an := "voter"
		ai := "101"
		conn_Id := respJson.Connection_id
		url = "http://10.53.17.40:8003/verify_proof"
		payload = []byte("{\"proof_attr\" : \"" + attrib + "\",\"connection_id\" : \"" + conn_Id + "\"}")

		fmt.Println("prepared payload", string(payload))
		req, _ = http.NewRequest("POST", url, bytes.NewBuffer(payload))

		req.Header.Add("content-type", "text/plain")

		res, _ = http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println(err)
		}

		defer res.Body.Close()
		body, _ = ioutil.ReadAll(res.Body)

		verifyrespJson := verifyProofOut{}
		err = json.Unmarshal(body, &verifyrespJson) //unmarshal it aka JSON.parse()
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("response received", verifyrespJson)
		fmt.Println("stringified response", string(body))
		fmt.Println(verifyrespJson.Status)

		if verifyrespJson.Status != "True" {
			return nil, nil, nil, errors.Errorf("Attributes missing !!!")
		}

		if !(verifyrespJson.Attributes.App_name == an && verifyrespJson.Attributes.App_id == ai) {
			return nil, nil, nil, errors.Errorf("Attribute values didnt match")
		}
	*/
	return true, nil, respJson.Connection_id
}
