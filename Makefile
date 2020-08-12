LISTPKG:=$(shell go list ./... | grep -v "/tests" | tr "\n" ",")

deps:
	@go get github.com/golang/mock/mockgen

generate:
	@echo -n Generating files and checking ...
	@go generate ./...
	@test $(shell git status --porcelain | wc -l) = 0 \
		|| { echo; echo "generated files are not up to date, re-generate and commit";\
		echo $(shell git status --porcelain);\
		exit 1;}\

	@echo " Everything is up to date"

test:
	@echo "Testing ..."
	@go test -count=1 -coverpkg=${LISTPKG:,=} -coverprofile=coverage.out -failfast ./...


