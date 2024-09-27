GREEN  = $(shell tput -Txterm setaf 2)
YELLOW = $(shell tput -Txterm setaf 3)
WHITE  = $(shell tput -Txterm setaf 7)
CYAN   = $(shell tput -Txterm setaf 6)
RESET  = $(shell tput -Txterm sgr0)

.PHONY: lint test vet tidy update_dependencies push-tag

lint: ## Lint the files
	golangci-lint run

test: ## Run unittests, '-p 1' makes tests not run in parallel;
	go test -p 1 -cover -race ./...

vet: ## Run vet
	go vet ./...

update_dependencies: ## Update Golang dependencies
	@go get -u ./...

tidy: ## Run tidy
	@go mod tidy

push-tag: ## increment the most recent tag and push it and any local commits
	bash createTag.sh

help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)