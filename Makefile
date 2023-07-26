.PHONY: server
server:
		echo "  --- running nats server on port :4222"
		nats-server -js &
		echo "  --- running service server on port :8080"
		go run cmd/server/main.go

.PHONY: client
client:
		go run cmd/client/main.go

.PHONY: up
up:
		echo "  --- running postgres on port :5432"
		docker-compose up