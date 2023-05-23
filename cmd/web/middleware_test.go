package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"gosnipit.ricci2511.dev/internal/assert"
)

func TestSecureHeaders(t *testing.T) {
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// mock the next handler in the chain
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// pass the mock handler to secureHeaders middlware and call its ServeHTTP method
	// with the recorder and dummy request
	secureHeaders(next).ServeHTTP(rr, r)

	rs := rr.Result()

	// check each of the response headers we expect
	expectedVal := "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com"
	assert.Equal(t, rs.Header.Get("Content-Security-Policy"), expectedVal)

	expectedVal = "origin-when-cross-origin"
	assert.Equal(t, rs.Header.Get("Referrer-Policy"), expectedVal)

	expectedVal = "nosniff"
	assert.Equal(t, rs.Header.Get("X-Content-Type-Options"), expectedVal)

	expectedVal = "deny"
	assert.Equal(t, rs.Header.Get("X-Frame-Options"), expectedVal)

	expectedVal = "0"
	assert.Equal(t, rs.Header.Get("X-XSS-Protection"), expectedVal)

	// check that the next handler in the chain is successfully called
	assert.Equal(t, rs.StatusCode, http.StatusOK)

	// check that the response body written by the handler equals "OK"
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")
}
