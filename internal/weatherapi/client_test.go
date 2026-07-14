package weatherapi_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mauricdias/whether-cep/internal/weatherapi"
)

func TestGetCurrentTemp_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"current":{"temp_c":28.5}}`))
	}))
	defer server.Close()

	client := weatherapi.NewClientWithBaseURL(server.URL, "test-key")
	temp, err := client.GetCurrentTemp("São Paulo")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if temp != 28.5 {
		t.Errorf("expected 28.5, got %f", temp)
	}
}

func TestGetCurrentTemp_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := weatherapi.NewClientWithBaseURL(server.URL, "invalid-key")
	_, err := client.GetCurrentTemp("São Paulo")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetCurrentTemp_NormalizesAccentedLocation(t *testing.T) {
	var gotQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.Query().Get("q")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"current":{"temp_c":18.2}}`))
	}))
	defer server.Close()

	client := weatherapi.NewClientWithBaseURL(server.URL, "test-key")
	_, err := client.GetCurrentTemp("Arujá, SP")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if gotQuery != "Aruja, SP" {
		t.Fatalf("expected normalized query to be %q, got %q", "Aruja, SP", gotQuery)
	}
}
