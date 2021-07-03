package auth

import (
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

func HandleGoogleAuth(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// parse the GoogleJWT that was POSTed from the front-end
		type GoogleJWT struct {
			// Make sure field name starts with capital letter
			// This makes sure its exported and visible to the JSON Decoder
			IdToken string
		}
		jwt := &GoogleJWT{}
		err := utils.DecodeJSON(w, r, jwt, true)
		if err != nil {
			utils.Respond(http.StatusBadRequest, err.Error()).ServeHTTP(w, r)
			return
		}

		// Validate the JWT
		claims, err := ValidateGoogleJWT(jwt.IdToken)
		if err != nil {
			utils.Respond(http.StatusUnauthorized, "Invalid google auth").ServeHTTP(w, r)
			return
		}
		// create a JWT for OUR app and give it back to the client for future requests
		// Stateful token authentication
		// tokenString, err := auth.MakeJWT(claims.Email, cfg.JWTSecret)
		// if err != nil {
		// 	respondWithError(w, 500, "Couldn't make authentication token")
		// 	return
		// }
		ctx := context.WithValue(r.Context(), "userData", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}
