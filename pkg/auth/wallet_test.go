package auth

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
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
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	// publicAddress := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	data := []byte(nonce)
	hash := crypto.Keccak256Hash(data)

	signatureBytes, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	return Payload{
		Nonce:           nonce,
		Signature:       hexutil.Encode(signatureBytes),
		WalletPublicKey: hexutil.Encode(publicKeyBytes),
	}

}

func TestVerifySignature(t *testing.T) {

	t.Run("True match on a valid payload", func(t *testing.T) {
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

	t.Run("False match on an empty payload", func(t *testing.T) {
		payload := Payload{
			Nonce:           "",
			Signature:       "",
			WalletPublicKey: "",
		}

		got, _ := verifySignature(payload)
		want := false

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
			Nonce:           "NONCE",
			Signature:       "0xsignature",
			WalletPublicKey: "0xpublicKeyString",
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
