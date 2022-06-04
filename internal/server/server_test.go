package server

import (
	"bytes"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testMux *mux.Router

func TestMain(m *testing.M) {
	testMux = mux.NewRouter()
	endpoints(testMux)
	m.Run()
}

func TestGetHealthz(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	req.Header.Add("content-type", "plain/text")

	rr := httptest.NewRecorder()

	testMux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected %v, got %v", http.StatusOK, rr.Code)
	}

	t.Log("good request for /healthz:", rr.Code)
}

func TestCreateUser(t *testing.T) {
	// bad request for post /user
	invalidUser := `{"name":"bob", "email":123}`

	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer([]byte(invalidUser)))
	req.Header.Add("content-type", "application/json")

	rr := httptest.NewRecorder()

	testMux.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected %v, got %v", http.StatusBadRequest, rr.Code)
	}

	t.Log("bad request for post /user:", rr.Code, rr.Body)

	// bad request for post /user if name and email are empty
	invalidUser = `{"name":"", "email": ""}`

	req = httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer([]byte(invalidUser)))
	req.Header.Add("content-type", "application/json")

	rr = httptest.NewRecorder()

	testMux.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected %v, got %v", http.StatusBadRequest, rr.Code)
	}

	t.Log("bad request for post /user if name and email are empty:", rr.Code, rr.Body)

	// good request for post /user
	validUser := `{"name":"bob","email":"123"}`

	req = httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer([]byte(validUser)))
	req.Header.Add("content-type", "application/json")

	rr = httptest.NewRecorder()

	testMux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Expected %v, got %v", http.StatusOK, rr.Code)
	}

	t.Log("good request for post /user:", rr.Code, rr.Body)
}
