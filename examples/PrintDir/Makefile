.DEFAULT_GOAL := buildandrun
BIN_FILE=strChain
buildandrun:
	@go build -o "${BIN_FILE}"
	./"${BIN_FILE}"

build:
	@go build -o "${BIN_FILE}"

clean:
	go clean
	rm --force "cp.out"
	rm --force nohup.out

test:
	go test

check:
	go test

cover:
	go test -coverprofile cp.out
	go tool cover -html=cp.out

run:
	./"${BIN_FILE}"

lint:
	golangci-lint run --enable-all
