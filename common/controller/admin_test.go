package controller

import (
	"crypto/sha256"
	"testing"
	"time"
)

func TestAdminController_RegisterPublicDid(t *testing.T) {

	var seed = Seed()

	// Create an admin controller for verifying the signature
	ac, err := NewAdminController()
	if err != nil {
		t.Errorf("Error initializing controller: %v\n", err)
	}

	// Register DID of client
	t.Run("Register Did With Ledger", func(t *testing.T) {

		did, err := RegisterDidWithLedger(ac, seed)
		if err != nil {
			t.Errorf("Error occurred while registering did: %v\n", err)
			return
		}

		did, err = ac.PublicDid()
		if err != nil {
			t.Errorf("Error occurred while getting did: %v\n", err)
		}

		t.Logf("admin Public DID:  %v\n", did)
	})
}

func TestAdminController_RequestPublicDid(t *testing.T) {

	// Create an admin controller for verifying the signature
	ac, err := NewAdminController()
	if err != nil {
		t.Errorf("Error initializing controller: %v\n", err)
	}

	did, err := ac.PublicDid()
	if err != nil {
		t.Errorf("Error requesting public did: %v\n", err)
	}

	t.Logf("Public DID: %v\n", did)
}

func TestAdminController_IssueCredential(t *testing.T) {

	// Create a client controller for signing
	cc, err := NewClientController()
	if err != nil {
		t.Errorf("Error occurred while initializing client: %v\n", err)
		return
	}

	// Create an admin controller for verifying the signature
	ac, err := NewAdminController()
	if err != nil {
		t.Errorf("Error occurred while registering did: %v\n", err)
		return
	}

	// Establish connection between client and admin
	t.Run("Establish connection", func(t *testing.T) {

		// create invitation
		inv, err := CreateInvitation(ac)
		if err != nil {
			t.Errorf("Error while creating invitation: %v\n", err)
			return
		}

		t.Logf("Invitation: %+v\n", inv)

		// receive invitation
		_, err = ReceiveInvitation(cc, inv)
		if err != nil {
			t.Errorf("Receive invitation failed: %v\n", err)
		}
	})

	time.Sleep(4 * time.Second)

	// Get connection ID of connection between admin and client (FOR ADMIN)
	t.Run("Get Connection Object", func(t *testing.T) {
		_, err := ac.GetConnection()
		if err != nil {
			t.Errorf("Failed to retrieve connection id: %v\n", err)
		}
		t.Logf("Connection:     %+v\n", ac.Connection)
	})

	t.Run("Register Schema and Cred Def", func(t *testing.T) {

		schemaID, err := ac.RegisterSchema("schema")
		if err != nil {
			t.Errorf("Error occurred while registering schema: %v\n", err)
			return
		}
		t.Logf("Schema ID : %v\n", schemaID)

		credDefID, err := ac.RegisterCredentialDefinition(schemaID)
		if err != nil {
			t.Errorf("Error occurred while registering cred def: %v\n", err)
			return
		}
		t.Logf("CredDef ID : %v\n", credDefID)
	})

	t.Run("Issue Credential", func(t *testing.T) {

		err := ac.IssueCredential("voter", "101")
		if err != nil {
			t.Errorf("Error occurred while issuing credential: %v\n", err)
			return
		}
	})
}

func TestAdminController_VerifySignature(t *testing.T) {

	var signature string
	var err error
	var message = "Foo bar"
	var seed = Seed()

	// Create a client controller for signing
	cc, err := NewClientController()
	if err != nil {
		t.Errorf("Error occurred while initializing client: %v\n", err)
		return
	}

	// Create an admin controller for verifying the signature
	ac, err := NewAdminController()
	if err != nil {
		t.Errorf("Error occurred while registering did: %v\n", err)
		return
	}

	// Register DID of client
	t.Run("Register Did With Ledger", func(t *testing.T) {

		_, err := RegisterDidWithLedger(ac, seed)
		if err != nil {
			t.Errorf("Error occurred while registering did: %v\n", err)
			return
		}

		_, err = RegisterDidWithLedger(cc, seed)
		if err != nil {
			t.Errorf("Error occurred while registering did: %v\n", err)
			return
		}

		adid, err := ac.PublicDid()
		if err != nil {
			t.Errorf("Error occurred while getting did: %v\n", err)
		}

		cdid, err := ac.PublicDid()
		if err != nil {
			t.Errorf("Error occurred while getting did: %v\n", err)
		}

		t.Logf("admin Public DID:  %v\n", adid)
		t.Logf("client Public DID: %v\n", cdid)
	})

	// Register DID of client
	t.Run("Create Signing Did", func(t *testing.T) {

		err = cc.CreateSigningDid()
		if err != nil {
			t.Errorf("Error occurred while registering did: %v\n", err)
			return
		}

		t.Logf("Signing DID: %v\n", cc.SigningDid)
		t.Logf("Signing VK: %v\n", cc.SigningVk)
	})

	// Establish connection between client and admin
	t.Run("Establish connection", func(t *testing.T) {

		// create invitation
		inv, err := CreateInvitation(ac)
		if err != nil {
			t.Errorf("Error while creating invitation: %v\n", err)
			return
		}

		t.Logf("Invitation: %+v\n", inv)

		// receive invitation
		_, err = ReceiveInvitation(cc, inv)
		if err != nil {
			t.Errorf("Receive invitation failed: %v\n", err)
		}
	})

	time.Sleep(2 * time.Second)

	// Get connection ID of connection between admin and client (FOR ADMIN)
	t.Run("Get Connection Object", func(t *testing.T) {
		_, err := ac.GetConnection()
		if err != nil {
			t.Errorf("Failed to retrieve connection id: %v\n", err)
		}
		t.Logf("Connection:     %+v\n", ac.Connection)
	})

	// Get connection ID of connection between admin and client (FOR CLIENT)
	t.Run("Get Connection Object", func(t *testing.T) {
		_, err := cc.GetConnection()
		if err != nil {
			t.Errorf("Failed to retrieve connection id: %v\n", err)
		}
		t.Logf("Connection:     %+v\n", cc.Connection)
	})

	// Sign a message with using the application signing DID
	t.Run("Sign Message", func(t *testing.T) {

		messageHash := sha256.Sum256([]byte(message))

		signature, err = cc.SignMessage(messageHash[:])
		if err != nil {
			t.Errorf("Error occurred during signing: %v\n", err)
		}

		t.Logf("Signature: %s\n", signature)
	})

	time.Sleep(2 * time.Second)

	t.Run("Put Key to Ledger", func(t *testing.T) {

		err := PutKeyToLedger(cc, cc.SigningDid, cc.SigningVk)
		if err != nil {
			t.Errorf("Failed to put key to ledger: %v\n", err)
		} else {
			t.Logf("Successfully put key to ledger\n")
		}

	})

	// Verify signature
	t.Run("Verify Signature", func(t *testing.T) {

		messageHash := sha256.Sum256([]byte(message))

		verified, err := ac.VerifySignature(messageHash[:], []byte(signature), []byte(cc.SigningDid))
		if err != nil {
			t.Errorf("Error occurred while attempting to verify signature: %v\n", err)
		}

		t.Logf("Verified: %v\n", verified)

		verified, err = ac.VerifySignature(messageHash[:], []byte(signature), []byte(cc.SigningDid))
		if err != nil {
			t.Errorf("Error occurred while attempting to verify signature: %v\n", err)
		}

		t.Logf("Verified: %v\n", verified)
	})
}

func TestAdminController_RequireProof(t *testing.T) {

	// Create a client controller for signing
	cc, err := NewClientController()
	if err != nil {
		t.Errorf("Error occurred while initializing client: %v\n", err)
		return
	}

	// Create an admin controller for verifying the signature
	ac, err := NewAdminController()
	if err != nil {
		t.Errorf("Error occurred while registering did: %v\n", err)
		return
	}

	// Establish connection between client and admin
	t.Run("Establish connection", func(t *testing.T) {

		// create invitation
		inv, err := CreateInvitation(ac)
		if err != nil {
			t.Errorf("Error while creating invitation: %v\n", err)
			return
		}

		t.Logf("Invitation: %+v\n", inv)

		// receive invitation
		_, err = ReceiveInvitation(cc, inv)
		if err != nil {
			t.Errorf("Receive invitation failed: %v\n", err)
		}
	})

	time.Sleep(4 * time.Second)

	// Get connection ID of connection between admin and client (FOR ADMIN)
	t.Run("Get Connection Object", func(t *testing.T) {
		_, err := ac.GetConnection()
		if err != nil {
			t.Errorf("Failed to retrieve connection id: %v\n", err)
		}
		t.Logf("Connection:     %+v\n", ac.Connection)
	})

	// Register DID of client
	t.Run("Register Did With Ledger", func(t *testing.T) {

		_, err := RegisterDidWithLedger(ac, Seed())
		if err != nil {
			t.Errorf("Error occurred while registering did: %v\n", err)
			return
		}

		did, err := ac.PublicDid()
		if err != nil {
			t.Errorf("Error occurred while getting did: %v\n", err)
		}

		t.Logf("Public DID : %v\n", did)
	})

	t.Run("Register Schema and Cred Def", func(t *testing.T) {

		schemaID, err := ac.RegisterSchema("schema")
		if err != nil {
			t.Errorf("Error occurred while registering schema: %v\n", err)
			return
		}
		t.Logf("Schema ID : %v\n", schemaID)

		credDefID, err := ac.RegisterCredentialDefinition(schemaID)
		if err != nil {
			t.Errorf("Error occurred while registering cred def: %v\n", err)
			return
		}
		t.Logf("CredDef ID : %v\n", credDefID)
	})

	t.Run("Issue Credential", func(t *testing.T) {

		err := ac.IssueCredential("voter", "101")
		if err != nil {
			t.Errorf("Error occurred while issuing credential: %v\n", err)
			return
		}
	})

	time.Sleep(5 * time.Second)

	t.Run("Request Proof", func(t *testing.T) {

		presExID, err := ac.RequireProof()
		if err != nil {
			t.Errorf("Error while trying to request proof: %v\n", err)
		}
		t.Logf("PresExID: %v\n", presExID)

		time.Sleep(5 * time.Second)

		verified, err := ac.CheckProofStatus(presExID)
		if err != nil {
			t.Errorf("Error while trying to check proof status: %v\n", err)
		}

		t.Logf("Verified: %v\n", verified)
	})

}

func TestAdminController_1234(t *testing.T) {

	// Create an admin controller for verifying the signature
	ac, err := NewAdminController()
	if err != nil {
		t.Errorf("Error occurred while registering did: %v\n", err)
		return
	}

	t.Run("Request Proof", func(t *testing.T) {

		presExID, err := ac.RequireProof()
		if err != nil {
			t.Errorf("Error while trying to request proof: %v\n", err)
		}
		t.Logf("PresExID: %v\n", presExID)

		verified, err := ac.CheckProofStatus(presExID)
		if err != nil {
			t.Errorf("Error while trying to check proof status: %v\n", err)
		}

		t.Logf("Verified: %v\n", verified)
	})
}