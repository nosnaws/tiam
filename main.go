package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"runtime"

	"strconv"

	"github.com/newrelic/go-agent/v3/newrelic"
	fastGame "github.com/nosnaws/tiam/game"
)

const ServerID = "nosnaws/tiam"

func recordLatency(app *newrelic.Application, state fastGame.GameState) {
	latency, err := strconv.ParseFloat(state.You.Latency, 64)
	if err == nil {
		app.RecordCustomMetric("lastTurnLatency", latency)
	}
}

// HTTP Handlers

func HandleIndex(app *newrelic.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := info()

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Printf("ERROR: Failed to encode info response, %s", err)
		}
	}
}

func HandleStart(app *newrelic.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state := fastGame.GameState{}
		err := json.NewDecoder(r.Body).Decode(&state)
		if err != nil {
			log.Printf("ERROR: Failed to decode start json, %s", err)
			return
		}

		start(state)

		// Nothing to respond with here
	}
}

func HandleMove(app *newrelic.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state := fastGame.GameState{}
		err := json.NewDecoder(r.Body).Decode(&state)
		if err != nil {
			log.Printf("ERROR: Failed to decode move json, %s", err)
			return
		}

		txn := newrelic.FromContext(r.Context())
		//getBaseAttributes(txn, state)

		//recordLatency(app, state)
		txn.AddAttribute("lastTurnLatency", state.You.Latency)

		response := move(state, txn)

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Printf("ERROR: Failed to encode move response, %s", err)
			return
		}
	}
}

func HandleEnd(app *newrelic.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		state := fastGame.GameState{}
		err := json.NewDecoder(r.Body).Decode(&state)
		if err != nil {
			log.Printf("ERROR: Failed to decode end json, %s", err)
			return
		}

		txn := newrelic.FromContext(r.Context())
		getCustomAttributesEnd(txn, state)
		txn.AddAttribute("lastTurnLatency", state.You.Latency)

		recordLatency(app, state)

		end(state)
		// Nothing to respond with here
	}

}

// Middleware

func withServerID(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", ServerID)
		next(w, r)
	}
}

// Main Entrypoint

func main() {
	log.Println("Version", runtime.Version())
	log.Println("NumCPU", runtime.NumCPU())
	log.Println("GOMAXPROCS", runtime.GOMAXPROCS(0))
	nrLicenseKey := os.Getenv("NEW_RELIC_LICENSE_KEY")
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	app, _ := newrelic.NewApplication(
		newrelic.ConfigAppName("Tiam"),
		newrelic.ConfigLicense(nrLicenseKey),
	)

	http.HandleFunc(newrelic.WrapHandleFunc(app, "/", withServerID(HandleIndex(app))))
	http.HandleFunc(newrelic.WrapHandleFunc(app, "/start", withServerID(HandleStart(app))))
	http.HandleFunc(newrelic.WrapHandleFunc(app, "/move", withServerID(HandleMove(app))))
	http.HandleFunc(newrelic.WrapHandleFunc(app, "/end", withServerID(HandleEnd(app))))

	log.Printf("Starting Battlesnake Server at http://0.0.0.0:%s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
