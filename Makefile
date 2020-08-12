LISTPKG:=$(shell go list ./... | grep -v "/tests" | tr "\n" ",")

generate:
	@echo "Generating mock files"
	@go generate ./...

test:
	@echo "Testing ..."
	@go test -count=1 -coverpkg=${LISTPKG:,=} -coverprofile=coverage.out ./...

