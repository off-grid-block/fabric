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
type attributes struct {
	App_name string `json:"app_name"`
	App_id   string `json:"app_id"`
}
type verifyProofOut struct {
	Status     string     `json:"status"`
	Attributes attributes `json:"attributes"`
}

func Indyverify(ProposalBytes []byte, Did []byte, Signature []byte) (status bool, err error) {

	hash := sha256.Sum256(ProposalBytes)
	encoded := b64.StdEncoding.EncodeToString(hash[:])
	fmt.Println()
	fmt.Println()
	fmt.Println("calculated hash", hash)
	fmt.Println("encoded hash", encoded)

	url := "http://10.53.17.40:8003/verify_signature"
	payload := []byte("{\"message\" : \"" + encoded + "\",\"their_did\" : \"" + string(Did) + "\",\"signature\": \"" + string(Signature) + "\"}")
	fmt.Println("prepared payload", string(payload))
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payload))

	req.Header.Add("content-type", "text/plain")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Received error from Indy server", err)
		return false, errors.New("Error connecting to Indy server")
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var result map[string]interface{}
	json.NewDecoder(res.Body).Decode(&result)

	fmt.Println(string(body))

	outJson := verifySignatureOut{}
	err = json.Unmarshal(body, &outJson)
	if err != nil {
		fmt.Println("error unmarshaling response from Indy", err)
		return false, errors.New("error unmarshaling response from Indy")
	}

	if outJson.Status != "Signature verified" {
		return false, errors.New(outJson.Status)
	}
	// Verify proof code
	/*
		attrib := "app_name,app_id"
		an := "voter"
		ai := "101"
		conn_Id := outJson.Connection_id
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

		verifyOutJson := verifyProofOut{}
		err = json.Unmarshal(body, &verifyOutJson) //unmarshal it aka JSON.parse()
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("response received", verifyOutJson)
		fmt.Println("stringified response", string(body))
		fmt.Println(verifyOutJson.Status)

		if verifyOutJson.Status != "True" {
			return nil, nil, nil, errors.Errorf("Attributes missing !!!")
		}

		if !(verifyOutJson.Attributes.App_name == an && verifyOutJson.Attributes.App_id == ai) {
			return nil, nil, nil, errors.Errorf("Attribute values didnt match")
		}
	*/
	return true, nil
}
