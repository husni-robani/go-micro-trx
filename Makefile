up_build:
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