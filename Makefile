build:
	@go build -o bin/api.exe


run: build
	@./bin/api.exe

seed: 
	@go run scripts/seed.go

test:
	@go test -v ./...