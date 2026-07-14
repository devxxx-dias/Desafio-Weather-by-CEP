package handler

import (
	"encoding/json"
	"errors"
	"log"
	"math"
	"net/http"
	"regexp"

	"github.com/mauricdias/whether-cep/internal/viacep"
)

var cepRegex = regexp.MustCompile(`^\d{8}$`)

type WeatherResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type ViaCEPClient interface {
	GetCity(cep string) (string, error)
}

type WeatherAPIClient interface {
	GetCurrentTemp(city string) (float64, error)
}

type WeatherHandler struct {
	viaCEP  ViaCEPClient
	weather WeatherAPIClient
}

func NewWeatherHandler(viaCEP ViaCEPClient, weather WeatherAPIClient) *WeatherHandler {
	return &WeatherHandler{viaCEP: viaCEP, weather: weather}
}

func (h *WeatherHandler) Handle(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")
	log.Printf("handling weather request: cep=%s", cep)

	if !cepRegex.MatchString(cep) {
		log.Printf("invalid cep format: cep=%s", cep)
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	city, err := h.viaCEP.GetCity(cep)
	log.Printf("viaCEP lookup result: cep=%s city=%q err=%v", cep, city, err)
	if err != nil {
		if errors.Is(err, viacep.ErrCEPNotFound) {
			http.Error(w, "can not find zipcode", http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	tempC, err := h.weather.GetCurrentTemp(city)
	log.Printf("weather lookup result: city=%q tempC=%f err=%v", city, tempC, err)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := WeatherResponse{
		TempC: roundToTwoDecimals(tempC),
		TempF: roundToTwoDecimals(celsiusToFahrenheit(tempC)),
		TempK: roundToTwoDecimals(celsiusToKelvin(tempC)),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func celsiusToFahrenheit(c float64) float64 {
	return c*1.8 + 32
}

func celsiusToKelvin(c float64) float64 {
	return c + 273
}

func roundToTwoDecimals(value float64) float64 {
	return math.Round(value*100) / 100
}
