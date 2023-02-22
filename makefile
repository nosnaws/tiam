install:
	go mod tidy

default: install

test:
	go test ./...

build:
	docker build -t tiam -f "Dockerfile.tiam" .

run:
	go run ./snakes/tiam/

compile:
	go build -o tiam ./snakes/tiam/

compile_lancer:
	go build -o lancer ./snakes/lancer/

build_random:
	docker build -t random -f "Dockerfile.random" .

build_eater:
	docker build -t eater -f "Dockerfile.eater" .

build_huey:
	docker build -t huey -f "Dockerfile.huey" .

build_mini:
	docker build -t mini -f "Dockerfile.mini" .

run_mini:
	go run ./snakes/mini/

run_lancer:
	go run ./snakes/lancer/

run_huey:
	go run ./snakes/huey/

run_mcts:
	go run ./snakes/mcts/

run_monte:
	go run ./snakes/monte_carlo/

build_all: build build_random build_eater

build_ga:
	go build -o genetic_algorithm ./gmo/evo_runner/main.go

run_ga: build_ga
	BS_EXE=../rules/battlesnake CSV_OUTPUT_DIR=evo_out ./genetic_algorithm



