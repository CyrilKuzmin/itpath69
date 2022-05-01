.DEFAULT_GOAL := help

GREEN='\033[0;32m'
ORANGE='\033[0;33m'
NC='\033[0m'

build: ## Compile golang binaries
	@echo ${ORANGE}"Building go binaries..."${NC}
	go build -o itpath69 ./cmd 
	@echo ${GREEN}"done"${NC}

tidy: ## Format the code
	@echo ${ORANGE}"Formatting the code..."${NC}
	go mod tidy
	go fmt ./...
	goimports -w -l -local git.dev.cloud.mts.ru .
	@echo ${GREEN}"done"${NC}

stop-pf: ## Stop port-forwarding
	## brew install proctools, in case of error here
	pkill kubectl

run: ## Run the service locally
	go run ./cmd