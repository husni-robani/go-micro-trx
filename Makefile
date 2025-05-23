BROKER_BINARY=brokerApp

up_build: build_broker
	@echo "Stoping docker image (if running) ..."
	docker compose down
	@echo "Building and starting docker images ..."
	docker compose up --build -d
	@echo "Docker images built and started!"

down:
	@echo "Stoping docker images ..."
	docker compose down

up:
	@echo "Starting docker images ..."
	docker compose up -d
	@echo "Docker images started!"

build_broker:
	@echo "Building binary broker service ..."
	cd broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
	@echo "Done!"