package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestBasicHandler(t *testing.T) {
	rr := sendRootRequest(t)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	var data struct {
		Success bool     `json:"success"`
		Otps    []string `json:"otps"`
	}

	err := json.Unmarshal([]byte(rr.Body.String()), &data)
	if err != nil {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), "success")
	}
	if data.Success != true {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), "success")
	}
	t.Logf("Otps: %v", data)
}

func TestInsertHandler(t *testing.T) {
	data := map[string]string{
		"message": "This is a test data with otp 1234",
	}
	rr := sendInsertRequest(t, data)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	containsString := strings.Contains(rr.Body.String(), "Successfully")
	if !containsString {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), "Successfully")
	}

	rr = sendRootRequest(t)
	containsString = strings.Contains(rr.Body.String(), "1234")
	if !containsString {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), "1234")
	}

	rr1 := sendInsertRequest(t, nil)
	if status := rr1.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnprocessableEntity)
	}
}

func sendRootRequest(t *testing.T) *httptest.ResponseRecorder {
	req, err := http.NewRequest("GET", "/otps", nil)
	os.Setenv("testing", "1")

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	//add header to request
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logFn := func(s interface{}) {}
		otpDisplayHandler(w, r, logFn)
	})
	handler.ServeHTTP(rr, req)
	return rr
}

func sendInsertRequest(t *testing.T, data map[string]string) *httptest.ResponseRecorder {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling data:", err)
		return nil
	}
	req, err := http.NewRequest("POST", "/notifyOtp", bytes.NewBuffer(jsonData))
	if data == nil {
		req, err = http.NewRequest("POST", "/notifyOtp", nil)
	}

	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logFn := func(s interface{}) {}
		otpHandler(w, r, logFn)
	})
	handler.ServeHTTP(rr, req)
	return rr
}
