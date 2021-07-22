package auth

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/cryptopatron/koen-backend/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// Create the JWT key used to create the signature
const jwtKey = "my_secret_key"

type Payload struct {
	Nonce           string `json:"nonce"`
	Signature       string `json:"signature"`
	WalletPublicKey string `json:"walletPublicKey"`
}

// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	WalletPublicKey string `json:"walletPublicKey"`
	jwt.StandardClaims
}

func verifySignature(payload Payload) (bool, error) {

	publicKeyBytes, err := hexutil.Decode(payload.WalletPublicKey)
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

func createJWT(c Claims) (string, error) {
	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func validateWalletJWT(token string) (Claims, error) {
	// Initialize a new instance of `Claims`
	claims := &Claims{}
	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// Finally, return the welcome message to the user, along with their
	// username given in the token
	w.Write([]byte(fmt.Sprintf("Welcome %s!", claims.Username)))
}

func HandleWalletAuthentication() http.HandlerFunc {

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
		// Declare the expiration time of the token
		expirationTime := time.Now().Add(50 * time.Minute)
		// Create the JWT claims, which includes the username and expiry time
		claims := &Claims{
			WalletPublicKey: payload.WalletPublicKey,
			StandardClaims: jwt.StandardClaims{
				// In JWT, the expiry time is expressed as unix milliseconds
				ExpiresAt: expirationTime.Unix(),
			},
		}
		jwt, err := createJWT(*claims)
		if err != nil {
			utils.Respond(http.StatusInternalServerError, "Something went wrong").ServeHTTP(w, r)
			return
		}
		utils.RespondWithJSON(map[string]string{"idToken": jwt}, 200)(w, r)

	})

}
