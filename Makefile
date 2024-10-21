.DEFAULT_GOAL := help

.PHONY: vm

test:  ## tests
	go test ./...
