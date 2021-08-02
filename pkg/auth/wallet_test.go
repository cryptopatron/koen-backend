package auth

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cryptopatron/koen-backend/pkg/utils"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func generateTestPayload(nonce string) Payload {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	publicAddress := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	data := []byte(nonce)
	// Prefix an ethereum-specific message and then hash
	hash := signHash(data)
	signatureBytes, err := crypto.Sign(hash, privateKey)
	if err != nil {
		log.Fatal(err)
	}
	// Need to add 27 to recovery identifier for legacy reasons
	signatureBytes[64] += 27
	fmt.Println("sign", hexutil.Encode(signatureBytes))

	return Payload{
		Nonce:               nonce,
		Signature:           hexutil.Encode(signatureBytes),
		WalletPublicAddress: publicAddress,
	}

}

func TestVerifySignature(t *testing.T) {

	t.Run("True match on a generated valid payload", func(t *testing.T) {
		payload := generateTestPayload("Hello hash")

		got, _ := verifySignature(payload)
		want := true

		if got != want {
			t.Errorf("got %t, want %t", got, want)
		}
	})

	t.Run("False match on an invalid payload", func(t *testing.T) {
		payload := generateTestPayload("Hello hash")

		payload.Signature = "bleh"

		got, _ := verifySignature(payload)
		want := false

		if got != want {
			t.Errorf("got %t, want %t", got, want)
		}
	})

	t.Run("True match on a sample valid payload", func(t *testing.T) {
		payload := Payload{
			Nonce:               "1627895001092",
			Signature:           "0xa2c97a7b4cd30a32c348c39449f4ece757513a8f6790a9a70e19b2359ba9b89b16e0b1f74bb21135f594a93215b797afe6de03fd1656aad2a0ed69e1fa4fea2a1c",
			WalletPublicAddress: "0xc2fde45f9e0a77005493930f72819fcf70210464",
		}

		got, err := verifySignature(payload)
		if err != nil {
			t.Error(err)
		}
		want := true

		if got != want {
			t.Errorf("got %t, want %t", got, want)
		}
	})
}

func TestHandleWalletAuthentication(t *testing.T) {

	t.Run("HTTP 400 on empty payload", func(t *testing.T) {
		payload := Payload{}
		// Payload to json string
		body, err := json.Marshal(payload)
		if err != nil {
			t.Error(err)
		}

		req, err := http.NewRequest("POST", "/test", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		// Create a ResponseRecorder which satisifies the interface of http.ResponseWriter
		rr := httptest.NewRecorder()

		// handler satisfies the interface of http.Handler
		// So we can use its ServeHTTP to serve the rquest to it
		HandleWalletAuthentication().ServeHTTP(rr, req)

		got := rr.Code
		want := http.StatusBadRequest

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

	t.Run("HTTP 400 on random payload", func(t *testing.T) {
		payload := Payload{
			Nonce:               "NONCE",
			Signature:           "0xsignature",
			WalletPublicAddress: "0xpublicKeyString",
		}
		// Payload to json string
		body, err := json.Marshal(payload)
		if err != nil {
			t.Error(err)
		}

		req, err := http.NewRequest("POST", "/test", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		// Create a ResponseRecorder which satisifies the interface of http.ResponseWriter
		rr := httptest.NewRecorder()

		// handler satisfies the interface of http.Handler
		// So we can use its ServeHTTP to serve the rquest to it
		HandleWalletAuthentication().ServeHTTP(rr, req)

		got := rr.Code
		want := http.StatusBadRequest

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

	t.Run("HTTP 200 and JWT on valid payload", func(t *testing.T) {
		payload := generateTestPayload("Hello hash")
		// Payload to json string
		body, err := json.Marshal(payload)
		if err != nil {
			t.Error(err)
		}

		req, err := http.NewRequest("POST", "/test", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		// Create a ResponseRecorder which satisifies the interface of http.ResponseWriter
		rr := httptest.NewRecorder()

		// handler satisfies the interface of http.Handler
		// So we can use its ServeHTTP to serve the rquest to it
		HandleWalletAuthentication().ServeHTTP(rr, req)

		got := rr.Code
		want := http.StatusOK

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
		// Buffer to JSON
		jwt := &JWT{}
		err = utils.DecodeJSON(rr.Body, jwt, false)
		if err != nil {
			t.Error(err)
		}
	})
}
