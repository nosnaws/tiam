package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"runtime"

	"strconv"

	"github.com/newrelic/go-agent/v3/newrelic"
	api "github.com/nosnaws/tiam/battlesnake"
	instru "github.com/nosnaws/tiam/instrumentation"
	mctsv3 "github.com/nosnaws/tiam/mcts_v3"
)

const ServerID = "nosnaws/tiam"

func recordLatency(app *newrelic.Application, state api.GameState) {
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

func HandleStart(gc *mctsv3.GameController, app *newrelic.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gameState := api.GameState{}
		err := json.NewDecoder(r.Body).Decode(&gameState)
		if err != nil {
			log.Printf("ERROR: Failed to decode start json, %s", err)
			return
		}

		start(gc, gameState)

		// Nothing to respond with here
	}
}

func HandleMove(gc *mctsv3.GameController, app *newrelic.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gameState := api.GameState{}
		err := json.NewDecoder(r.Body).Decode(&gameState)
		if err != nil {
			log.Printf("ERROR: Failed to decode move json, %s", err)
			return
		}

		txn := newrelic.FromContext(r.Context())
		//getBaseAttributes(txn, state)

		//recordLatency(app, state)
		txn.AddAttribute("lastTurnLatency", gameState.You.Latency)

		response := move(gc, gameState, txn)

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Printf("ERROR: Failed to encode move response, %s", err)
			return
		}
	}
}

func HandleEnd(gc *mctsv3.GameController, app *newrelic.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		gameState := api.GameState{}
		err := json.NewDecoder(r.Body).Decode(&gameState)
		if err != nil {
			log.Printf("ERROR: Failed to decode end json, %s", err)
			return
		}

		txn := newrelic.FromContext(r.Context())
		instru.GetCustomAttributesEnd(txn, gameState)
		txn.AddAttribute("lastTurnLatency", gameState.You.Latency)

		recordLatency(app, gameState)

		end(gc, gameState)
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

//func parseEnv() *brain.MCTSConfig {
//exploration, err := strconv.ParseFloat(os.Getenv("EXPLORATION_CONSTANT"), 64)
//if err != nil {
//log.Println("EXPLORATION_CONSTANT not set, defaulting")
//exploration = 2.0
//}

//alpha, err := strconv.ParseFloat(os.Getenv("ALPHA_CONSTANT"), 64)
//if err != nil {
//alpha = 0.1
//log.Println("ALPHA_CONSTANT not set, defaulting", alpha)
//}

//voronoi, err := strconv.ParseFloat(os.Getenv("VORONOI_CONSTANT"), 64)
//if err != nil {
//voronoi = 0.5
//log.Println("VORONOI_CONSTANT not set, defaulting", voronoi)
//}

//foodA, err := strconv.ParseFloat(os.Getenv("FOOD_CONSTANT_A"), 64)
//if err != nil {
//foodA = 0.1
//log.Println("FOOD_CONSTANT_A not set, defaulting", foodA)
//}

//foodB, err := strconv.ParseFloat(os.Getenv("FOOD_CONSTANT_B"), 64)
//if err != nil {
//foodB = 2
//log.Println("FOOD_CONSTANT_B not set, defaulting", foodB)
//}

//bigSnake, err := strconv.ParseFloat(os.Getenv("BIG_SNAKE_CONSTANT"), 64)
//if err != nil {
//bigSnake = 2
//log.Println("BIG_SNAKE_CONSTANT not set, defaulting", bigSnake)
//}

//return &brain.MCTSConfig{
//ExplorationConstant: exploration,
//AlphaConstant:       alpha,
//VoronoiWeighting:    voronoi,
//FoodWeightA:         foodA,
//FoodWeightB:         foodB,
//BigSnakeReward:      bigSnake,
//}

//}

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

	gc := mctsv3.CreateGameController()

	http.HandleFunc(newrelic.WrapHandleFunc(app, "/", withServerID(HandleIndex(app))))
	http.HandleFunc(newrelic.WrapHandleFunc(app, "/start", withServerID(HandleStart(gc, app))))
	http.HandleFunc(newrelic.WrapHandleFunc(app, "/move", withServerID(HandleMove(gc, app))))
	http.HandleFunc(newrelic.WrapHandleFunc(app, "/end", withServerID(HandleEnd(gc, app))))

	log.Printf("Starting Battlesnake Server at http://0.0.0.0:%s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
