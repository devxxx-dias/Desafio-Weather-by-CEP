package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mauricdias/whether-cep/internal/handler"
	"github.com/mauricdias/whether-cep/internal/viacep"
	"github.com/mauricdias/whether-cep/internal/weatherapi"
)

func main() {

	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		log.Fatal("WEATHER_API_KEY environment variable is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	viaCEPClient := viacep.NewClient()
	weatherClient := weatherapi.NewClient(apiKey)
	weatherHandler := handler.NewWeatherHandler(viaCEPClient, weatherClient)

	mux := http.NewServeMux()
	mux.HandleFunc("/weather", weatherHandler.Handle)

	addr := fmt.Sprintf(":%s", port)
	log.Printf("server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
