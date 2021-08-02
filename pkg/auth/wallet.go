package auth

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/cryptopatron/koen-backend/pkg/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// Create the JWT key used to create the signature
const jwtKey = "my_secret_key"

type Payload struct {
	Nonce               string `json:"nonce"`
	Signature           string `json:"signature"`
	WalletPublicAddress string `json:"walletPublicAddress"`
}

// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type WalletClaims struct {
	WalletPublicAddress string `json:"walletPublicAddress"`
	jwt.StandardClaims
}

// web3js prefixes a message to its messages before hashing them
// So we're doing the same
func signHash(data []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256([]byte(msg))
}

func verifySignature(payload Payload) (bool, error) {

	publicAddr := common.HexToAddress(payload.WalletPublicAddress)

	data := []byte(payload.Nonce)

	sig, err := hexutil.Decode(payload.Signature)
	if err != nil {
		return false, err
	}

	// Last byte of signature is the recovery ID (recid)
	// It has 2 possible values: 0,1; each will result in a different public address
	// web3js adds 27 to recid in its signatures
	// We are checking and removing it here
	if sig[64] != 27 && sig[64] != 28 {
		return false, errors.New("Not a valid signature. 'v' value should be 27/28")
	}
	sig[64] -= 27

	pubKey, err := crypto.SigToPub(signHash(data), sig)
	if err != nil {
		return false, err
	}

	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	fmt.Println("Public", recoveredAddr.Hex())

	return publicAddr == recoveredAddr, nil
}

func createJWT(c WalletClaims) (string, error) {
	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// Create the JWT string
	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (wc *WalletClaims) ValidateJWT(tokenString string) error {
	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tokenString, wc, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if err != nil {
		return err
	}
	if !tkn.Valid {
		return errors.New("Invalid JWT!")
	}

	fmt.Printf("Wallet %s has been validated!", wc.WalletPublicAddress)
	return nil
}

// TODO: Write tests for this!
func HandleWalletAuthentication() http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		payload := &Payload{}
		err := utils.DecodeJSON(r.Body, payload, true)
		if err != nil {
			utils.Respond(http.StatusBadRequest, err.Error()).ServeHTTP(w, r)
			return
		}

		// Cryptographically verify the given signature
		result, err := verifySignature(*payload)
		if err != nil {
			fmt.Println(err)
			utils.Respond(http.StatusBadRequest, err.Error()).ServeHTTP(w, r)
			return
		}

		if !result {
			utils.Respond(http.StatusUnauthorized, "Invalid auth").ServeHTTP(w, r)
			return
		}
		// Declare the expiration time of the token
		expirationTime := time.Now().Add(50 * time.Minute)
		// Create the JWT claims, which includes the username and expiry time
		claims := &WalletClaims{
			WalletPublicAddress: payload.WalletPublicAddress,
			StandardClaims: jwt.StandardClaims{
				// In JWT, the expiry time is expressed as unix milliseconds
				ExpiresAt: expirationTime.Unix(),
			},
		}
		jwt, err := createJWT(*claims)
		if err != nil {
			fmt.Println(err)
			utils.Respond(http.StatusInternalServerError, "Something went wrong").ServeHTTP(w, r)
			return
		}
		utils.RespondWithJSON(map[string]string{"idToken": jwt}, 200)(w, r)

	})

}
