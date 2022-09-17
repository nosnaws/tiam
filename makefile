install:
	go mod tidy

default: install

test:
	go test ./...

build:
	docker build -t tiam -f "Dockerfile.tiam" .

run:
	go run ./snakes/tiam/

build_random:
	docker build -t random -f "Dockerfile.random" .

build_eater:
	docker build -t eater -f "Dockerfile.eater" .

build_mini:
	docker build -t mini -f "Dockerfile.mini" .

run_mini:
	go run ./snakes/mini/

build_all: build build_random build_eater

build_ga:
	go build -o genetic_algorithm ./gmo/evo_runner/main.go

run_ga: build_ga
	BS_EXE=../rules/battlesnake CSV_OUTPUT_DIR=evo_out ./genetic_algorithm



