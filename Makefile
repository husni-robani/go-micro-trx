BROKER_EXECUTABLE=brokerApp
LOGGER_EXECUTABLE=loggerApp
TASK_EXECUTABLE=taskApp
TRANSACTION_EXECUTABLE=transactionApp
MAIL_EXECUTABLE=mailApp

up_build: build_broker build_logger build_task build_transaction build_mail
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
	@echo "Building binary broker-service ..."
	cd broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_EXECUTABLE} ./cmd
	@echo "Build broker service done!"

build_logger:
	@echo "Building binary logger-service ..."
	cd logger-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LOGGER_EXECUTABLE} ./cmd
	@echo "Build logger service done!"

build_task:
	@echo "building binary task-service ..."
	cd task-service && env GOOS=linux CGO_ENABLED=0 go build -o ${TASK_EXECUTABLE} ./cmd
	@echo "Build task service done!"

build_transaction:
	@echo "building binary transaction-service ..."
	cd transaction-service && env GOOS=linux CGO_ENABLED=0 go build -o ${TRANSACTION_EXECUTABLE} ./cmd
	@echo "Build transaction service done"

build_mail:
	@echo "building binary mail-service ..."
	cd mail-service && env GOOS=linux CGO_ENABLED=0 go build -o ${MAIL_EXECUTABLE} ./cmd
	@echo "Build mail service done"