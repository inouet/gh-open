# parameters

BINARY_NAME=gh-open
VETPKGS = $(shell go list ./... | grep -v -e vendor)


.PHONY: build
build:
	gox --osarch "darwin/amd64 linux/amd64 windows/amd64" -output="bin/{{.OS}}_{{.Arch}}/$(BINARY_NAME)"

.PHONY: clean
clean:
	go clean
	rm -rf bin/

.PHONY: test
test:
	go test -v ./...

.PHONY: lint
lint:
	go list ./... | xargs golint -set_exit_status

.PHONY: deps
deps:
	go get -d -v .
	go mod tidy

.PHONY: vet
vet:
	go vet $(VETPKGS)
