package controller

import (
	"encoding/json"
	"fmt"
)

type ConnectionInvitation struct {
	ConnectionID string `json:"connection_id"`
	InvitationURL string `json:"invitation_url"`
	Invitation Invitation
}

type Invitation struct {
	ID string `json:"@id"`
	Type string `json:"@type"`
	//Did string `json:"did"`
	Label string `json:"label"`
	RecipientKeys []string `json:"recipientKeys"`
	//RoutingKeys []string `json:"routingKeys"`
	ServiceEndpoint string `json:"serviceEndpoint"`
	SigningDid string `json:"signing_did"`
}

// struct to hold get connection call results
type GetConnectionResponse struct {
	Results []Connection `json:"results"`
}

// struct representing an agent connection
type Connection struct {
	MyDID string `json:"my_did"`
	TheirDID string `json:"their_did"`
	ConnectionID string `json:"connection_id"`
	State string `json:"state"`
}

func CreateInvitation(controller Controller) (Invitation, error) {

	var inv Invitation

	resp, err := SendRequestWithParams_POST(
		controller.AgentUrl(),
		"/connections/create-invitation",
		map[string]string{"alias": "deon"},
		nil)
	if err != nil {
		return inv, fmt.Errorf("Error occurred while sending post request: %v\n", err)
	}
	defer resp.Body.Close()

	var connInv ConnectionInvitation
	err = json.NewDecoder(resp.Body).Decode(&connInv)
	if err != nil {
		return inv, fmt.Errorf("Error occurred while decoding json: %v\n", err)
	}

	return connInv.Invitation, nil
}

func ReceiveInvitation(controller Controller, invitation Invitation) (*Connection, error) {

	resp, err := SendRequest_POST(controller.AgentUrl(), "/connections/receive-invitation", invitation)
	if err != nil {
		return nil, fmt.Errorf("error occurred while sending post request: %v\n", err)
	}
	defer resp.Body.Close()

	var connection Connection
	err = json.NewDecoder(resp.Body).Decode(&connection)
	if err != nil {
		return nil, fmt.Errorf("error occurred while decoding json into connection object: %v\n", err)
	}

	return &connection, nil
}