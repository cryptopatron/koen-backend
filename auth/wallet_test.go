package auth

import (
	"crypto/ecdsa"
	"log"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestVerifySignature(t *testing.T) {
	const NONCE = "Hello hash"
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

	data := []byte(NONCE)
	hash := crypto.Keccak256Hash(data)

	signatureBytes, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	signature := hexutil.Encode(signatureBytes)
	publicKeyString := hexutil.Encode(publicKeyBytes)

	t.Run("True match on a valid payload", func(t *testing.T) {
		payload := Payload{
			Nonce:           NONCE,
			Signature:       signature,
			WalletPublicKey: publicKeyString,
		}

		got, _ := verifySignature(payload)
		want := true

		if got != want {
			t.Errorf("got %t, want %t", got, want)
		}
	})

	t.Run("False match on an invalid payload", func(t *testing.T) {
		payload := Payload{
			Nonce:           NONCE,
			Signature:       "bleh",
			walletPublicKey: publicKeyString,
		}

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
			walletPublicKey: "",
		}

		got, _ := verifySignature(payload)
		want := false

		if got != want {
			t.Errorf("got %t, want %t", got, want)
		}
	})
}
