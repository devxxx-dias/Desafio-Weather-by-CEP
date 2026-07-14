package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mauricdias/whether-cep/internal/handler"
	"github.com/mauricdias/whether-cep/internal/viacep"
)

type mockViaCEP struct {
	city string
	err  error
}

func (m *mockViaCEP) GetCity(_ string) (string, error) {
	return m.city, m.err
}

type mockWeatherAPI struct {
	tempC float64
	err   error
}

func (m *mockWeatherAPI) GetCurrentTemp(_ string) (float64, error) {
	return m.tempC, m.err
}

func TestHandle_Success(t *testing.T) {
	h := handler.NewWeatherHandler(
		&mockViaCEP{city: "São Paulo"},
		&mockWeatherAPI{tempC: 25.0},
	)

	req := httptest.NewRequest(http.MethodGet, "/?cep=01310100", nil)
	w := httptest.NewRecorder()
	h.Handle(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp handler.WeatherResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.TempC != 25.0 {
		t.Errorf("expected TempC=25.0, got %f", resp.TempC)
	}
	expectedF := 25.0*1.8 + 32
	if resp.TempF != expectedF {
		t.Errorf("expected TempF=%f, got %f", expectedF, resp.TempF)
	}
	expectedK := 25.0 + 273
	if resp.TempK != expectedK {
		t.Errorf("expected TempK=%f, got %f", expectedK, resp.TempK)
	}
}

func TestHandle_InvalidZipcode_TooShort(t *testing.T) {
	h := handler.NewWeatherHandler(&mockViaCEP{}, &mockWeatherAPI{})

	req := httptest.NewRequest(http.MethodGet, "/?cep=0131", nil)
	w := httptest.NewRecorder()
	h.Handle(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", w.Code)
	}
}

func TestHandle_InvalidZipcode_Letters(t *testing.T) {
	h := handler.NewWeatherHandler(&mockViaCEP{}, &mockWeatherAPI{})

	req := httptest.NewRequest(http.MethodGet, "/?cep=0131abc0", nil)
	w := httptest.NewRecorder()
	h.Handle(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", w.Code)
	}
}

func TestHandle_InvalidZipcode_Empty(t *testing.T) {
	h := handler.NewWeatherHandler(&mockViaCEP{}, &mockWeatherAPI{})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	h.Handle(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", w.Code)
	}
}

func TestHandle_CEPNotFound(t *testing.T) {
	h := handler.NewWeatherHandler(
		&mockViaCEP{err: viacep.ErrCEPNotFound},
		&mockWeatherAPI{},
	)

	req := httptest.NewRequest(http.MethodGet, "/?cep=99999999", nil)
	w := httptest.NewRecorder()
	h.Handle(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestTemperatureConversions(t *testing.T) {
	cases := []struct {
		tempC     float64
		expectedF float64
		expectedK float64
	}{
		{0, 32, 273},
		{100, 212, 373},
		{-40, -40, 233},
		{25, 77, 298},
	}

	for _, tc := range cases {
		h := handler.NewWeatherHandler(
			&mockViaCEP{city: "TestCity"},
			&mockWeatherAPI{tempC: tc.tempC},
		)

		req := httptest.NewRequest(http.MethodGet, "/?cep=01310100", nil)
		w := httptest.NewRecorder()
		h.Handle(w, req)

		var resp handler.WeatherResponse
		json.NewDecoder(w.Body).Decode(&resp)

		if resp.TempF != tc.expectedF {
			t.Errorf("tempC=%.1f: expected F=%.1f, got %.1f", tc.tempC, tc.expectedF, resp.TempF)
		}
		if resp.TempK != tc.expectedK {
			t.Errorf("tempC=%.1f: expected K=%.1f, got %.1f", tc.tempC, tc.expectedK, resp.TempK)
		}
	}
}
