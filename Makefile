LISTPKG:=$(shell go list ./... | grep -vE "/tests|/mock|/mortar/http/server/health|/mortar/http/server/proto" | tr "\n" ",")

deps:
	@go install github.com/golang/mock/mockgen golang.org/x/tools/cmd/goimports golang.org/x/lint/golint

generate:
	@echo -n Generating files...
	@go generate ./...
	@echo " done"
	
generate-test: deps generate
	@echo -n Checking ...
	@go generate ./...
	@test $(shell git status --porcelain | wc -l) = 0 \
		|| { echo; echo "generated files are not up to date, re-generate and commit";\
		echo $(shell git status --porcelain);\
		exit 1;}\

	@echo " everything is up to date"

go-fmt:
	@echo -n Checking format...
	@test $(shell goimports -l ./ | grep -v mock | wc -l) = 0 \
		|| { echo; echo "some files are not properly formatted";\
		echo $(shell goimports -l ./ | grep -v mock);\
		exit 1;}\

	@echo " everything formatted properly"	

go-lint:
	@echo -n Checking with linter...
	@test $(shell golint ./... | wc -l) = 0 \
		|| { echo; echo "some files are not properly linted";\
		echo $(shell golint ./...);\
		exit 1;}\

	@echo " everything linted properly"	

cover-report:
	@go tool cover -html=coverage.out -o coverage.html
	@go tool cover -func=coverage.out -o coverage-summary.txt
	@tail -n1 coverage-summary.txt

test:
	@echo "Testing ..."
	@go test -count=1 -coverpkg=${LISTPKG:,=} -coverprofile=coverage.out -failfast ./...

monitor-race-test:
	@echo "Monitoring race test..."
	@go test -race ./monitoring/...

test-with-report: test monitor-race-test cover-report

code-up-to-date: generate-test go-fmt go-lint

health-proto:
	 @protoc -I. \
			-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
			--go_out=plugins=grpc:. \
          	--grpc-gateway_out=logtostderr=true,paths=source_relative:. ./http/server/health/health.proto

service-test-proto:
	 @protoc -I. \
			-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
			--go_out=plugins=grpc:. \
          	--grpc-gateway_out=logtostderr=true,paths=source_relative:. ./http/server/proto/service.proto

proto: health-proto service-test-proto

all: code-up-to-date test-with-report 
