build:
	@go build -o bin/api.exe


run: build
	@./bin/api.exe

seedteams: 
	@go run scripts/seed_teams.go

seedmatches: 
	@go run scripts/seed_matches.go

test:
	@go test -v ./...