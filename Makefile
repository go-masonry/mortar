LISTPKG:=$(shell go list ./... | grep -vE "/tests|/mock|/mortar/http/server/health|/mortar/http/server/proto" | tr "\n" ",")

deps:
	@go install github.com/golang/mock/mockgen

generate:
	@echo -n Generating files and checking ...
	@go generate ./...
	@test $(shell git status --porcelain | wc -l) = 0 \
		|| { echo; echo "generated files are not up to date, re-generate and commit";\
		echo $(shell git status --porcelain);\
		exit 1;}\

	@echo " Everything is up to date"

cover-report:
	@go tool cover -html=coverage.out -o coverage.html
	@go tool cover -func=coverage.out -o coverage-summary.txt

test:
	@echo "Testing ..."
	@go test -count=1 -coverpkg=${LISTPKG:,=} -coverprofile=coverage.out -failfast ./...


