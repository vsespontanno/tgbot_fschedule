build:
	@go build -o bin/api


run: build
	@./bin/api

seedteams: 
	@go run scripts/seed_teams/seed_teams.go

seedmatches: 
	@go run scripts/seed_matches/seed_matches.go

seedstandings: 
	@go run scripts/seed_standings/seed_standings.go

drop: 
	@go run scripts/drop_coll/drop_coll.go

test:
	@go test -v ./...