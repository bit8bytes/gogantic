# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## install: install all dev dependencies
.PHONY: install
install: confirm
	go mod download

## examples/core/pipe: to translate an input from english to spanish
.PHONY: examples/core/pipe
examples/core/pipe:
	go run ./examples/core/pipe

## examples/agents/ollama/temperature: to get the temperature using a tool (temperature hardocded to 5.54F)
.PHONY: examples/agents/ollama/temperature
examples/agents/ollama/temperature:
	go run ./examples/agents/ollama/temperature

## examples/agents/ollama/time: the agent accesses the local time of the device
.PHONY: examples/agents/ollama/time
examples/agents/ollama/time:
	go run ./examples/agents/ollama/time

## examples/agents/ollama/filesystem: the agent has access to the filesystem
.PHONY: examples/agents/ollama/filesystem
examples/agents/ollama/filesystem:
	go run ./examples/agents/ollama/filesystem