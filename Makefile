BROKER_EXECUTABLE=brokerApp
LOGGER_EXECUTABLE=loggerApp

up_build: build_broker build_logger
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
	cd broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_EXECUTABLE} ./cmd/api
	@echo "Build broker service done!"

build_logger:
	@echo "Building binary logger service ..."
	cd logger-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LOGGER_EXECUTABLE} ./cmd/api
	@echo "Build logger service done!"