# Weather by CEP — Go Expert Challenge

A Go service that receives a Brazilian ZIP code (CEP), looks up the corresponding city via ViaCEP, and returns the current temperature in Celsius, Fahrenheit, and Kelvin via WeatherAPI.

## Live URL (Google Cloud Run)

Use this deployed endpoint:

```text
https://weather-api-cep-full-cycle-264412937559.region.run.app/weather?cep=01310100
```

---

## API Contract

### `GET /weather?cep={cep}`

**Success — 200 OK**

```json
{
  "temp_C": 15.1,
  "temp_F": 59.18,
  "temp_K": 288.1
}
```

**Failure cases**

| Condition                                           | Status | Body                   |
| --------------------------------------------------- | ------ | ---------------------- |
| CEP without 8 digits or with non-numeric characters | 422    | `invalid zipcode`      |
| Valid format but CEP not found in the database      | 404    | `can not find zipcode` |

**Conversion formulas**

- Fahrenheit: `F = C × 1.8 + 32`
- Kelvin: `K = C + 273`

---

## Prerequisites

- [Go 1.22+](https://golang.org/dl/)
- [Docker](https://docs.docker.com/get-docker/)
- A [WeatherAPI](https://www.weatherapi.com/) free-tier API key

---

## Running locally (without Docker)

```bash
export WEATHER_API_KEY=your_api_key_here
go run .
```

The server starts on port `8080`. Test it:

```bash
curl "http://localhost:8080/weather?cep=01310100"
```

---

## Running locally with Docker

Create a `.env` file in the project root with:

```env
WEATHER_API_KEY=your_api_key_here
PORT=8080
```

Then run the app with Docker Compose:

```bash
docker compose up --build
```

This starts the container on port `8080` and passes the API key into the app.

You can also run the container manually:

```powershell
# Build the image
docker build -t whether-cep .

# Run the container on Windows PowerShell
$env:WEATHER_API_KEY="your_api_key_here"
docker run -p 8080:8080 -e WEATHER_API_KEY="$env:WEATHER_API_KEY" whether-cep
```

On Windows, the quoted form above is the safest way to pass the API key.

Test the endpoints:

```bash
# Valid CEP → 200
curl "http://localhost:8080/weather?cep=01310100"

# Invalid format (not 8 digits) → 422
curl "http://localhost:8080/weather?cep=0131"

# Non-existent CEP → 404
curl "http://localhost:8080/weather?cep=00000000"
```

---

## Running tests

```bash
go test ./... -v
```

---

## Deploying to Google Cloud Run

### Option A — via the Google Cloud Console (no CLI needed)

1. Push your code to GitHub (see [Push to GitHub](#push-to-github) below).
2. Go to [console.cloud.google.com/run](https://console.cloud.google.com/run).
3. Click **"Create Service"**.
4. Choose **"Continuously deploy from a repository"** and connect your GitHub repo.
5. Set **branch** to `main` and **build type** to `Dockerfile`.
6. Under **"Variables & Secrets"**, add the environment variable:
   - Key: `WEATHER_API_KEY`
   - Value: your WeatherAPI key
7. Under **"Authentication"**, select **"Allow unauthenticated invocations"**.
8. Click **Deploy**.

After the deploy completes, Cloud Run shows your service URL. Test it:

```bash
curl "https://<your-service-url>/weather?cep=01310100"
```

---

### Option B — via the gcloud CLI

Make sure you have the [gcloud CLI](https://cloud.google.com/sdk/docs/install) installed and authenticated.

```bash
# Set your GCP project
export PROJECT_ID=your-gcp-project-id

# Build and push the image to Container Registry
gcloud builds submit --tag gcr.io/$PROJECT_ID/whether-cep

# Deploy to Cloud Run
gcloud run deploy whether-cep \
  --image gcr.io/$PROJECT_ID/whether-cep \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars WEATHER_API_KEY=your_api_key_here
```

The command prints the service URL when done. Test it:

```bash
curl "https://<your-service-url>/weather?cep=01310100"
```

---

## Push to GitHub

If you haven't pushed the project yet:

```bash
git init
git add .
git commit -m "feat: weather by CEP service"

# Create a new empty repo on github.com, then:
git remote add origin https://github.com/your-username/whether-cep.git
git push -u origin main
```

---

## Project structure

```
.
├── main.go                        # Server entry point
├── internal/
│   ├── handler/
│   │   ├── weather.go             # HTTP handler + temperature conversions
│   │   └── weather_test.go        # Handler tests (mocked clients)
│   ├── viacep/
│   │   ├── client.go              # ViaCEP API client
│   │   └── client_test.go         # Client tests (mock HTTP server)
│   └── weatherapi/
│       ├── client.go              # WeatherAPI client
│       └── client_test.go         # Client tests (mock HTTP server)
├── Dockerfile                     # Multi-stage Docker build
└── go.mod
```

---

## External APIs used

| API                                      | Purpose                     | Docs                                            |
| ---------------------------------------- | --------------------------- | ----------------------------------------------- |
| [ViaCEP](https://viacep.com.br)          | Resolves CEP to city name   | `GET https://viacep.com.br/ws/{cep}/json/`      |
| [WeatherAPI](https://www.weatherapi.com) | Current temperature by city | `GET http://api.weatherapi.com/v1/current.json` |
