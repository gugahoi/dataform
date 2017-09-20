.PHONY: install test build publish clean

dcr := docker-compose run --rm
release_bucket := s3://myob-dataform-release

install:
	@echo "--- Collecting ingredients :bento:"
	@${dcr} dep ensure -v

test:
	@echo "+++ Is this thing working? :hammer_and_wrench:"
	@${dcr} go test ./... -v

build:
	@echo "+++ Laying bricks...:building_construction:"
	@echo "--- :linux: 64-bit"
	${dcr} -e CGO_ENABLED=0 -e GOOS=linux -e GOARCH=amd64 go build -v -o build/dfm-linux-amd64
	@echo "--- :mac: 64-bit"
	${dcr} -e CGO_ENABLED=0 -e GOOS=darwin -e GOARCH=amd64 go build -v -o build/dfm-darwin-amd64
	@echo "--- :windows: 64-bit"
	${dcr} -e CGO_ENABLED=0 -e GOOS=windows -e GOARCH=amd64 go build -v -o build/dfm-windows-amd64

publish: build
	@echo "+++ Flyyyyy!"
	${dcr} aws s3 sync build/ ${release_bucket}/latest/

clean:
	@echo "--- C ya later"
	${dcr} go clean
	rm -rf build/
