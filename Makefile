.PHONY: install test build publish clean

dcr := docker-compose run --rm
release_bucket := s3://myob-dataform-release

install:
	@echo "--- Collecting ingredients :bento:"
	@${dcr} dep ensure -v

test:
	@echo "+++ Is this thing working? :hammer_and_wrench:"
	@${dcr} go test ./... -v -cover

build: install
	@echo "+++ Laying bricks...:building_construction:"
	@echo "--- :linux: 64-bit"
	${dcr} -e CGO_ENABLED=0 -e GOOS=linux -e GOARCH=amd64 go build -v -o dist/dfm-linux-amd64
	@echo "--- :mac: 64-bit"
	${dcr} -e CGO_ENABLED=0 -e GOOS=darwin -e GOARCH=amd64 go build -v -o dist/dfm-darwin-amd64
	@echo "--- :windows: 64-bit"
	${dcr} -e CGO_ENABLED=0 -e GOOS=windows -e GOARCH=amd64 go build -v -o dist/dfm-windows-amd64

ifdef BUILDKITE_TAG
publish: install
	@# GITHUB_TOKEN can be found in buildkite secrets bucket
	@echo "+++ :octocat:"
	${dcr} goreleaser --skip-validate --rm-dist
	@echo "+++ :s3:"
	${dcr} aws s3 sync dist/ ${release_bucket}/${BUILDKITE_TAG}/
else
publish: build
	@echo "+++ :s3:"
	${dcr} aws s3 sync dist/ ${release_bucket}/latest/
endif

clean:
	@echo "--- C ya later"
	${dcr} go clean
	rm -rf dist/
