package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthCheckHandler)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestHandler1(t *testing.T) {
	u := User{ID: 76876, Balance: 93453452.66, accountID: "DFGKJ234DSFG"}
	payload, err := json.Marshal(u)

	req, err := http.NewRequest("POST", "/handler1", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handler1)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var expected User
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatal(err)
	}
	if err = json.Unmarshal(body, &expected); err != nil {
		t.Fatal(err)
	}

	if expected.ID != u.ID && expected.Balance != u.Balance {
		t.Errorf("handler returned unexpected body: got %v want %v", string(body), expected)
	}
}

func TestHandler2(t *testing.T) {
	u := User{ID: 123, Balance: 43.0, accountID: "ABCDEF123"}

	req, err := http.NewRequest("POST", "/handler2", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handler2)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var expected User
	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatal(err)
	}
	if err = json.Unmarshal(body, &expected); err != nil {
		t.Fatal(err)
	}

	if expected.ID != u.ID && expected.Balance != u.Balance {
		t.Errorf("handler returned unexpected body: got %v want %v", string(body), expected)
	}
}
