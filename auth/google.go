package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/cryptopatron/koen-backend/utils"
	"github.com/dgrijalva/jwt-go"
)

// GoogleClaims -
type GoogleClaims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	FirstName     string `json:"given_name"`
	LastName      string `json:"family_name"`
	jwt.StandardClaims
}

// ValidateGoogleJWT -
func ValidateGoogleJWT(tokenString string) (GoogleClaims, error) {
	claimsStruct := GoogleClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) {
			pem, err := getGooglePublicKey(fmt.Sprintf("%s", token.Header["kid"]))
			if err != nil {
				return nil, err
			}
			key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pem))
			if err != nil {
				return nil, err
			}
			return key, nil
		},
	)
	if err != nil {
		return GoogleClaims{}, err
	}

	claims, ok := token.Claims.(*GoogleClaims)
	if !ok {
		return GoogleClaims{}, errors.New("Invalid Google JWT")
	}

	if claims.Issuer != "accounts.google.com" && claims.Issuer != "https://accounts.google.com" {
		return GoogleClaims{}, errors.New("iss is invalid")
	}

	if claims.Audience != "116852492535-37n739s732ui71hkfm19n5r3agv6g9c5.apps.googleusercontent.com" {
		return GoogleClaims{}, errors.New("aud is invalid")
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return GoogleClaims{}, errors.New("JWT is expired")
	}

	return *claims, nil
}

func getGooglePublicKey(keyID string) (string, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v1/certs")
	if err != nil {
		return "", err
	}
	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	myResp := map[string]string{}
	err = json.Unmarshal(dat, &myResp)
	if err != nil {
		return "", err
	}
	key, ok := myResp[keyID]
	if !ok {
		return "", errors.New("key not found")
	}
	return key, nil
}

// Middleware
func HandleGoogleAuth(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// parse the GoogleJWT that was POSTed from the front-end
		type GoogleJWT struct {
			// Make sure field name starts with capital letter
			// This makes sure its exported and visible to the JSON Decoder
			IdToken string `json:"idToken"`
		}

		if r.Body == nil {
			utils.Respond(http.StatusBadRequest, "EMopty body").ServeHTTP(w, r)
			return
		}
		//Read request body into a copy buffer
		copyBuf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			utils.Respond(http.StatusInternalServerError, err.Error()).ServeHTTP(w, r)
		}

		jwt := &GoogleJWT{}
		// Passing in copy of request body to decode
		err = utils.DecodeJSON(bytes.NewReader(copyBuf), jwt, true)
		fmt.Println("jwt", jwt)
		if err != nil {
			utils.Respond(http.StatusBadRequest, err.Error()).ServeHTTP(w, r)
			return
		}

		// Validate the JWT
		claims, err := ValidateGoogleJWT(jwt.IdToken)
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
