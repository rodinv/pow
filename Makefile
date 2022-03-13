mod:
	@go mod download

server:
	@go run -race ./cmd/server/main.go

client:
	@go run -race ./cmd/client/main.go

docker_build:
	@docker build -t pow_server -f build/server.Dockerfile .
	@docker build -t pow_client -f build/client.Dockerfile .

docker_run:
	@docker-compose up -d

docker_stop:
	@docker-compose stop
