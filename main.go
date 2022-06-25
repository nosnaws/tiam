package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"

	"github.com/newrelic/go-agent/v3/newrelic"
)

const ServerID = "nosnaws/tiam"

type GameState struct {
	Game  Game        `json:"game"`
	Turn  int         `json:"turn"`
	Board Board       `json:"board"`
	You   Battlesnake `json:"you"`
}

type Game struct {
	ID      string  `json:"id"`
	Ruleset Ruleset `json:"ruleset"`
	Timeout int32   `json:"timeout"`
}

type Ruleset struct {
	Name     string   `json:"name"`
	Version  string   `json:"version"`
	Settings Settings `json:"settings"`
}

type Settings struct {
	FoodSpawnChance     int32  `json:"foodSpawnChance"`
	MinimumFood         int32  `json:"minimumFood"`
	HazardDamagePerTurn int32  `json:"hazardDamagePerTurn"`
	Royale              Royale `json:"royale"`
	Squad               Squad  `json:"squad"`
}

type Royale struct {
	ShrinkEveryNTurns int32 `json:"shrinkEveryNTurns"`
}

type Squad struct {
	AllowBodyCollisions bool `json:"allowBodyCollisions"`
	SharedElimination   bool `json:"sharedElimination"`
	SharedHealth        bool `json:"sharedHealth"`
	SharedLength        bool `json:"sharedLength"`
}

type Board struct {
	Height int           `json:"height"`
	Width  int           `json:"width"`
	Food   []Coord       `json:"food"`
	Snakes []Battlesnake `json:"snakes"`

	// Used in non-standard game modes
	Hazards []Coord `json:"hazards"`
}

type Customizations struct {
	Color string `json:"color"`
	Head  string `json:"head"`
	Tail  string `json:"tail"`
}

type Battlesnake struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	Health         int32          `json:"health"`
	Body           []Coord        `json:"body"`
	Head           Coord          `json:"head"`
	Length         int32          `json:"length"`
	Latency        string         `json:"latency"`
	Customizations Customizations `json:"customizations"`

	// Used in non-standard game modes
	Shout string `json:"shout"`
	Squad string `json:"squad"`
}

type Coord struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Response Structs

type BattlesnakeInfoResponse struct {
	APIVersion string `json:"apiversion"`
	Author     string `json:"author"`
	Color      string `json:"color"`
	Head       string `json:"head"`
	Tail       string `json:"tail"`
}

type BattlesnakeMoveResponse struct {
	Move  string `json:"move"`
	Shout string `json:"shout,omitempty"`
}

func recordLatency(app *newrelic.Application, state GameState) {
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
		state := GameState{}
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
		state := GameState{}
		err := json.NewDecoder(r.Body).Decode(&state)
		if err != nil {
			log.Printf("ERROR: Failed to decode move json, %s", err)
			return
		}

		txn := newrelic.FromContext(r.Context())
		getBaseAttributes(txn, state)

		recordLatency(app, state)
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
		state := GameState{}
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
		func(config *newrelic.Config) {
			config.DistributedTracer.Enabled = true
		},
	)

	http.HandleFunc(newrelic.WrapHandleFunc(app, "/", withServerID(HandleIndex(app))))
	http.HandleFunc(newrelic.WrapHandleFunc(app, "/start", withServerID(HandleStart(app))))
	http.HandleFunc(newrelic.WrapHandleFunc(app, "/move", withServerID(HandleMove(app))))
	http.HandleFunc(newrelic.WrapHandleFunc(app, "/end", withServerID(HandleEnd(app))))

	log.Printf("Starting Battlesnake Server at http://0.0.0.0:%s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
