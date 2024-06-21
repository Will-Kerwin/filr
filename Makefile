build:
	@go build -o bin/filr

run: build
	@./bin/filr

test:
	@go test ./... -v