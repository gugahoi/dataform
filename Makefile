.PHONY: install test build publish clean

dcr := docker-compose run --rm

install:
	@echo "--- Collecting ingredients :bento:"
	@${dcr} dep ensure -v

test:
	@echo "+++ Is this thing working? :hammer_and_wrench:"
	@${dcr} go test ./... -v -cover

build: clean install
	@echo "+++ Laying bricks...:building_construction:"
	@echo "--- :linux: 64-bit"
	${dcr} -e CGO_ENABLED=0 -e GOOS=linux -e GOARCH=amd64 go build -v -o dist/dfm-linux-amd64/dfm-linux-amd64
	@echo "--- :mac: 64-bit"
	${dcr} -e CGO_ENABLED=0 -e GOOS=darwin -e GOARCH=amd64 go build -v -o dist/dfm-darwin-amd64/dfm-darwin-amd64
	@echo "--- :windows: 64-bit"
	${dcr} -e CGO_ENABLED=0 -e GOOS=windows -e GOARCH=amd64 go build -v -o dist/dfm-windows-amd64/dfm-windows-amd64.exe

ifdef TRAVIS_TAG
publish: install
	@echo "+++ :octocat:"
	${dcr} goreleaser --skip-validate --rm-dist
endif

clean:
	@echo "--- C ya later"
	${dcr} go clean
	${dcr} base rm -rf dist/
