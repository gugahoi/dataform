.PHONY: install test

dcr := docker-compose run --rm

install:
	@echo "--- Installing Dependencies"
	@${dcr} dep ensure -v

test:
	@echo "+++ Running Unit Tests"
	@${dcr} go test ./... -v
