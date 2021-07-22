package auth

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cryptopatron/koen-backend/pkg/utils"
)

type JWT struct {
	// Make sure field name starts with capital letter
	// This makes sure its exported and visible to the JSON Decoder
	IdToken string `json:"idToken"`
}

type Claims struct {
	WalletClaims
	GoogleClaims
}

func (c *Claims) ValidateJWT(token string) error {
	errW := c.WalletClaims.ValidateJWT(token)
	errG := c.GoogleClaims.ValidateJWT(token)
	// If atleast one of them passed
	// JWT is legit
	if errW != nil && errG != nil {
		return errors.New("Invalid JWT!")
	}
	return nil
}

// Middleware
func HandleJWT(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Body == nil {
			utils.Respond(http.StatusBadRequest, "Empty body").ServeHTTP(w, r)
			return
		}
		//Read request body into a copy buffer
		copyBuf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			utils.Respond(http.StatusInternalServerError, err.Error()).ServeHTTP(w, r)
		}

		jwt := &JWT{}
		// Passing in copy of request body to decode
		err = utils.DecodeJSON(bytes.NewReader(copyBuf), jwt, true)
		fmt.Println("jwt", jwt)
		if err != nil {
			utils.Respond(http.StatusBadRequest, err.Error()).ServeHTTP(w, r)
			return
		}

		// Validate the JWT
		claims := &Claims{}
		err = claims.ValidateJWT(jwt.IdToken)
		if err != nil {
			fmt.Println(err)
			utils.Respond(http.StatusUnauthorized, "Invalid google auth").ServeHTTP(w, r)
			return
		}

		// Create user data context from validated JWT claims
		ctx := context.WithValue(r.Context(), "userData", claims)
		// Regenerate request body from copyBuffer
		r.Body = ioutil.NopCloser(bytes.NewBuffer(copyBuf))
		// Pass request with regenerated body and user data context to next HTTP Handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}
