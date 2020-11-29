package controller

import (
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
)

const (
	clientUrl = "http://host.docker.internal:8031"
	// clientUrl = "http://localhost:8031"
)

type ClientController struct {
	alias      string
	did        string
	SigningDid string
	SigningVk  string
	Connection Connection
	agentUrl   string
}

func NewClientController() (*ClientController, error) {
	return &ClientController{
		alias: "client",
		agentUrl: clientUrl,
	}, nil
}

func (cc *ClientController) Alias() string {
	return cc.alias
}

func (cc *ClientController) AgentUrl() string {
	return cc.agentUrl
}

func (cc *ClientController) PublicDid() (string, error) {
	return cc.did, nil
}

func (cc *ClientController) SetPublicDid(did string) {
	cc.did = did
}

func (cc *ClientController) ConnectionDid() string {
	return cc.Connection.MyDID
}

func (cc *ClientController) GetSigningDid() string {
	return cc.SigningDid
}

func (cc *ClientController) GetSigningVk() string {
	return cc.SigningVk
}

type SignMessageRequest struct {
	Message string `json:"message"`
	SigningDid string `json:"signing_did"`
}

type SignMessageResponse struct {
	Signature string `json:"signature"`
}

// signs the provided message
func (cc *ClientController) SignMessage(messageHash []byte) (string, error) {

	if cc.SigningDid == "" {
		return "", fmt.Errorf("no signing did, create a signing did before attempting to sign message")
	}

	encoded := b64.StdEncoding.EncodeToString(messageHash)

	payload := SignMessageRequest{
		Message:    encoded,
		SigningDid: cc.SigningDid,
	}

	resp, err := SendRequest_POST(clientUrl, "/connections/sign-transaction", payload)
	if err != nil {
		return "", fmt.Errorf("Failed to send post request: %v\n", err)
	}
	defer resp.Body.Close()

	var smResp SignMessageResponse
	err = json.NewDecoder(resp.Body).Decode(&smResp)
	if err != nil {
		return "", fmt.Errorf("Failed to decode json: %v\n", err)
	}

	return smResp.Signature, nil

}

type CreateSigningDidResponse struct {
	SigningDid string `json:"signing_did"`
	SigningVk string `json:"signing_vk"`
}

func (cc *ClientController) CreateSigningDid() error {
	resp, err := SendRequest_POST(clientUrl, "/connections/create-signing-did", nil)
	if err != nil {
		return fmt.Errorf("Failed to send post request: %v\n", err)
	}
	defer resp.Body.Close()

	var regDidResp CreateSigningDidResponse
	err = json.NewDecoder(resp.Body).Decode(&regDidResp)
	if err != nil {
		return fmt.Errorf("Failed to decode json: %v\n", err)
	}

	cc.SigningDid = regDidResp.SigningDid
	cc.SigningVk = regDidResp.SigningVk
	return nil
}

// Get connection ID of connection with admin agent
func (cc *ClientController) GetConnection() (GetConnectionResponse, error) {

	var getConnResp GetConnectionResponse

	resp, err := SendRequest_GET(
		clientUrl,
		"/connections",
		nil,
	)
	if err != nil {
		return getConnResp, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&getConnResp)
	if err != nil {
		return getConnResp, err
	}

	if len(getConnResp.Results) == 0 {
		return getConnResp, errors.New("no connections found")
	}

	// save the connection ID
	cc.Connection = getConnResp.Results[0]
	return getConnResp, nil
}