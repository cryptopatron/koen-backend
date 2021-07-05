package auth

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cryptopatron/koen-backend/utils"
)

func TestHandleGoogleAuth(t *testing.T) {
	handler := HandleGoogleAuth(utils.Respond(http.StatusOK, ""))

	t.Run("Bad request on empty request body", func(t *testing.T) {
		// Create a HTTP request with no JWT
		req, err := http.NewRequest("POST", "/auth/google/jwt", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Create a ResponseRecorder which satisifies the interface of http.ResponseWriter
		rr := httptest.NewRecorder()

		// handler satisfies the interface of http.Handler
		// So we can use its ServeHTTP to serve the rquest to it
		handler.ServeHTTP(rr, req)

		got := rr.Code
		want := http.StatusBadRequest

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("Bad request with random request body", func(t *testing.T) {
		// Create a HTTP request with no JWT
		body := bytes.NewBuffer([]byte("bleh"))
		req, err := http.NewRequest("POST", "/auth/google/jwt", body)
		if err != nil {
			t.Fatal(err)
		}

		// Create a ResponseRecorder which satisifies the interface of http.ResponseWriter
		rr := httptest.NewRecorder()

		// handler satisfies the interface of http.Handler
		// So we can use its ServeHTTP to serve the rquest to it
		handler.ServeHTTP(rr, req)

		got := rr.Code
		want := http.StatusBadRequest

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("Auth failed on random JWT", func(t *testing.T) {

		var jwt string = `{"idToken":"fakeasstoken"}`
		body := strings.NewReader(jwt)
		req, err := http.NewRequest("POST", "/auth/google/jwt", body)
		if err != nil {
			t.Fatal(err)
		}

		// TODO: Test to check if Content-Type is being checked for JSON
		// req.Header.Set("Content-Type", "application/json")

		// Create a ResponseRecorder which satisifies the interface of http.ResponseWriter
		rr := httptest.NewRecorder()

		// handler satisfies the interface of http.Handler
		// So we can use its ServeHTTP to serve the rquest to it
		handler.ServeHTTP(rr, req)

		got := rr.Code
		want := http.StatusUnauthorized

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	// t.Run("HTTP 200 on correct JWT but random JSON fields", func(t *testing.T) {

	// 	var jwt string = `{
	// 		"pageName": "divs",
	// 		"idToken": "eyJhbGciOiJSUzI1NiIsImtpZCI6ImI2ZjhkNTVkYTUzNGVhOTFjYjJjYjAwZTFhZjRlOGUwY2RlY2E5M2QiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJhY2NvdW50cy5nb29nbGUuY29tIiwiYXpwIjoiMTE2ODUyNDkyNTM1LTM3bjczOXM3MzJ1aTcxaGtmbTE5bjVyM2FndjZnOWM1LmFwcHMuZ29vZ2xldXNlcmNvbnRlbnQuY29tIiwiYXVkIjoiMTE2ODUyNDkyNTM1LTM3bjczOXM3MzJ1aTcxaGtmbTE5bjVyM2FndjZnOWM1LmFwcHMuZ29vZ2xldXNlcmNvbnRlbnQuY29tIiwic3ViIjoiMTA0NTM0MjY4MjEyMTk4MzA1MzkwIiwiZW1haWwiOiJkaXZnMjM5NUBnbWFpbC5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiYXRfaGFzaCI6IjYtSkVaZmozanp4ZHFRbU5Hd2RSWkEiLCJuYW1lIjoiZGl2IGd1cHRhIiwicGljdHVyZSI6Imh0dHBzOi8vbGgzLmdvb2dsZXVzZXJjb250ZW50LmNvbS9hL0FBVFhBSndjN2JHUjRMNk5nMGNEY1QyTUx0WjZwaEpyNFRSWVZ5SnFFaW9XPXM5Ni1jIiwiZ2l2ZW5fbmFtZSI6ImRpdiIsImZhbWlseV9uYW1lIjoiZ3VwdGEiLCJsb2NhbGUiOiJlbiIsImlhdCI6MTYyNTQ4MTMzMywiZXhwIjoxNjI1NDg0OTMzLCJqdGkiOiI4ZTFmMjI1ZTdmNmI3MzVhYmIyZmFlODk3YmY5NmQxMWQ0MDA0YzBhIn0.2ZSf6sOP5gljwfzRQqdzJEiFw_uTX9C0QLT9Dpgm2M7ptAFGhnPdw2aVlyJ2VxFaG4QM-IzrF235g7zPvVHYjUqe9DMq5OVDnXAnTda0yGyfJBoqKsAnJn1663jo9fv1pKNyPxbHkELNlaZF83Dhd_Q4JHA8yvRnAtbOJCe2kgDgaFkdfOzx1eF4Hat50GWAx_rtJfDnyG6o5i4ntM8QiAr-n8b-tdkMMNuAE1Sukcfiv34V6rvQOfHIEhdVyHS8ohBpecJDcFFhBHFxvnPMeeA4X2JYJosquHZ2AkAKKyGQrs5fdWXt-JCC38XLUdOz3be7_iSJUKwhRkuSjiWQZg",
	// 		"random": "random",
	// 		"metaMaskWalletPublicKey": "",
	// 		"generatedMaticWalletPublicKey": "kuhgihjygyuh"
	// 	}`
	// 	body := strings.NewReader(jwt)
	// 	req, err := http.NewRequest("POST", "/test", body)
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}

	// 	// TODO: Test to check if Content-Type is being checked for JSON
	// 	// req.Header.Set("Content-Type", "application/json")

	// 	// Create a ResponseRecorder which satisifies the interface of http.ResponseWriter
	// 	rr := httptest.NewRecorder()

	// 	// handler satisfies the interface of http.Handler
	// 	// So we can use its ServeHTTP to serve the rquest to it
	// 	handler.ServeHTTP(rr, req)

	// 	got := rr.Code
	// 	want := http.StatusOK

	// 	if got != want {
	// 		t.Errorf("got %v, want %v", got, want)
	// 	}
	// })
}
