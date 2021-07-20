package auth

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/cryptopatron/koen-backend/utils"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type Payload struct {
	Nonce                   string `json:"nonce"`
	Signature               string `json:"signature"`
	MetaMaskWalletPublicKey string `json:"metaMaskWalletPublicKey"`
}

func verifySignature(payload Payload) (bool, error) {

	publicKeyBytes, err := hexutil.Decode(payload.MetaMaskWalletPublicKey)
	if err != nil {
		fmt.Print(err)
		return false, err
	}

	// Generate hash of Nonce
	data := []byte(payload.Nonce)
	hash := crypto.Keccak256Hash(data)

	signature, err := hexutil.Decode(payload.Signature)
	if err != nil {
		return false, err
	}

	sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), signature)
	if err != nil {
		return false, err
	}

	match := bytes.Equal(sigPublicKey, publicKeyBytes)

	return match, nil
}

func HandleMetaMaskAuthentication() http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		payload := &Payload{}
		err := utils.DecodeJSON(r.Body, payload, true)
		fmt.Println("Payload", payload)
		if err != nil {
			utils.Respond(http.StatusBadRequest, err.Error()).ServeHTTP(w, r)
			return
		}

		// Cryptographically verify the given signature
		result, err := verifySignature(*payload)
		if err != nil {
			fmt.Println(err)
			// Bad request here!
		}

		if !result {
			utils.Respond(http.StatusUnauthorized, "Invalid auth").ServeHTTP(w, r)
			return
		}

		// Make and return a new JWT here

	})

}
