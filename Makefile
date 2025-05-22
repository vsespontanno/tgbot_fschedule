build:
	@go build -o bin/api ./cmd/bot


run: build
	@./bin/api

seedteams: 
	@go run internal/scripts/seed_teams/seed_teams.go

seedmatches: 
	@go run internal/scripts/seed_matches/seed_matches.go

seedstandings: 
	@go run internal/scripts/seed_standings/seed_standings.go

drop: 
	@go run internal/scripts/drop_coll/drop_coll.go

test:
	@go test -v ./...