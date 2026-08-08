package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func init() {
	globalMeta = Metadata{
		InstanceID:       "i-test12345",
		Hostname:         "test-node.internal",
		LocalIPv4:        "10.0.2.15",
		PublicIPv4:       "198.51.100.1",
		AMIID:            "flatcar-test-v1",
		AvailabilityZone: "us-west-2a",
		PublicKeys:       []string{"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAITestKey"},
		UserData:         "#cloud-config\nusers:\n  - name: testuser\n",
	}
}

func TestIMDSRootEndpoint(t *testing.T) {
	mux := setupRoutes()
	req := httptest.NewRequest("GET", "/latest/meta-data/", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "instance-id"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler response body missing expected text: got %v want substring %v", rr.Body.String(), expected)
	}
}

func TestInstanceIDEndpoint(t *testing.T) {
	mux := setupRoutes()
	req := httptest.NewRequest("GET", "/latest/meta-data/instance-id", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if rr.Body.String() != "i-test12345" {
		t.Errorf("unexpected body: got %v want %v", rr.Body.String(), "i-test12345")
	}
}

func TestUserDataEndpoint(t *testing.T) {
	mux := setupRoutes()
	req := httptest.NewRequest("GET", "/latest/user-data", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if !strings.Contains(rr.Body.String(), "#cloud-config") {
		t.Errorf("expected cloud-config user data, got: %s", rr.Body.String())
	}
}

func TestJSONEndpoint(t *testing.T) {
	mux := setupRoutes()
	req := httptest.NewRequest("GET", "/latest/meta-data.json", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("expected application/json content type, got: %s", contentType)
	}
}
