.DEFAULT_GOAL := help

GREEN='\033[0;32m'
ORANGE='\033[0;33m'
NC='\033[0m'

build: ## Compile golang binaries
	@echo ${ORANGE}"Building go binaries..."${NC}
	go build -o itpath69 ./cmd 
	@echo ${GREEN}"done"${NC}

run: ## Run the service locally
	go run ./cmd
