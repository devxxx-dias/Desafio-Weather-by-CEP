package viacep_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mauricdias/whether-cep/internal/viacep"
)

func TestGetCity_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"localidade":"São Paulo","erro":false}`))
	}))
	defer server.Close()

	client := viacep.NewClientWithBaseURL(server.URL)
	city, err := client.GetCity("01310100")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if city != "São Paulo" {
		t.Errorf("expected São Paulo, got %s", city)
	}
}

func TestGetCity_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"erro":"true"}`))
	}))
	defer server.Close()

	client := viacep.NewClientWithBaseURL(server.URL)
	_, err := client.GetCity("99999999")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err != viacep.ErrCEPNotFound {
		t.Errorf("expected ErrCEPNotFound, got %v", err)
	}
}

func TestGetCity_BadRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	client := viacep.NewClientWithBaseURL(server.URL)
	_, err := client.GetCity("00000000")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
