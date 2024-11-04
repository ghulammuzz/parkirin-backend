stg:
	go run cmd/main.go --env=stg
prod:
	go run cmd/main.go --env=prod
up:
	go mod tidy