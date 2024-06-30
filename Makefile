build:
	@go build -o bin/output

run: build
	@./bin/output