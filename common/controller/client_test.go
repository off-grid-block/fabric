package controller

import (
	"testing"
)

func TestClientController_SignMessage(t *testing.T) {

	cc, err := NewClientController()
	if err != nil {
		t.Errorf("Error while initializing client: %v\n", err)
		return
	}

	t.Run("Register Did", func(t *testing.T) {

		err := cc.CreateSigningDid()
		if err != nil {
			t.Errorf("Error occurred while registering did: %v\n", err)
			return
		}

		t.Logf("Signing DID: %v\n", cc.SigningDid)
	})

	t.Run("Sign Message", func(t *testing.T) {

		signature, err := cc.SignMessage([]byte("foo bar"))
		if err != nil {
			t.Errorf("Error occurred during signing: %v\n", err)
		}

		t.Logf("Signature: %s\n", signature)
	})
}
