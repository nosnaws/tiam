package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	_ "net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"strconv"

	"github.com/newrelic/go-agent/v3/newrelic"
	api "github.com/nosnaws/tiam/battlesnake"
	instru "github.com/nosnaws/tiam/instrumentation"
)

const ServerID = "nosnaws/huey"

func recordLatency(app *newrelic.Application, state api.GameState) {
	latency, err := strconv.ParseFloat(state.You.Latency, 64)
	if err == nil {
		app.RecordCustomMetric("lastTurnLatency", latency)
	}
}

// HTTP Handlers

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	response := info()

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("ERROR: Failed to encode info response, %s", err)
	}
}

func HandleStart(w http.ResponseWriter, r *http.Request) {
	state := api.GameState{}
	err := json.NewDecoder(r.Body).Decode(&state)
	if err != nil {
		log.Printf("ERROR: Failed to decode start json, %s", err)
		return
	}

	start(state)

	// Nothing to respond with here
}

func HandleMove(w http.ResponseWriter, r *http.Request) {
	duration, err := time.ParseDuration("350ms")
	if err != nil {
		panic("could not parse duration")
	}

	ctx, cancel := context.WithTimeout(r.Context(), duration)
	defer cancel()

	state := api.GameState{}
	err = json.NewDecoder(r.Body).Decode(&state)
	if err != nil {
		log.Printf("ERROR: Failed to decode move json, %s", err)
		return
	}

	txn := newrelic.FromContext(r.Context())
	//getBaseAttributes(txn, state)

	//recordLatency(app, state)
	txn.AddAttribute("lastTurnLatency", state.You.Latency)

	response := move(ctx, state)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("ERROR: Failed to encode move response, %s", err)
		return
	}
}

func HandleEnd(w http.ResponseWriter, r *http.Request) {

	state := api.GameState{}
	err := json.NewDecoder(r.Body).Decode(&state)
	if err != nil {
		log.Printf("ERROR: Failed to decode end json, %s", err)
		return
	}

	txn := newrelic.FromContext(r.Context())
	instru.GetCustomAttributesEnd(txn, state)
	txn.AddAttribute("lastTurnLatency", state.You.Latency)

	end(state)
	// Nothing to respond with here
}

// Middleware

func withServerID(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", ServerID)
		next(w, r)
	}
}

func logRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	}
}

// get the count of number of go routines in the system.
func countGoRoutines() int {
	return runtime.NumGoroutine()
}

func getGoroutines(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the count of number of go routines running.
		next(w, r)
		count := countGoRoutines()
		log.Printf("NUM GOROUTINES %d", count)
	}
}

func HandleGoRoutines(w http.ResponseWriter, r *http.Request) {
	pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
}

func HandleNumGoRoutines(w http.ResponseWriter, r *http.Request) {
	// Get the count of number of go routines running.
	count := countGoRoutines()
	log.Printf("NUM GOROUTINES %d", count)
}

// Main Entrypoint

func main() {
	log.Println("Version", runtime.Version())
	log.Println("NumCPU", runtime.NumCPU())
	log.Println("GOMAXPROCS", runtime.GOMAXPROCS(0))
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	initialize()
	http.HandleFunc("/", withServerID(HandleIndex))
	http.HandleFunc("/start", withServerID(HandleStart))
	http.HandleFunc("/move", getGoroutines(logRequest(withServerID((HandleMove)))))
	http.HandleFunc("/end", withServerID(HandleEnd))
	http.HandleFunc("/routines", withServerID(HandleNumGoRoutines))

	log.Printf("Starting Battlesnake Server at http://0.0.0.0:%s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
