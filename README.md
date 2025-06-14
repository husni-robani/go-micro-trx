This is a simple microservices project built in Go, simulating a basic transaction system. It consists of four main services — **Broker**, **Task**, **Transaction**, and **Logger** — along with additional components like **Mail** and **RabbitMQ** integration. The purpose of this project is to explore service communication, logging, orchestration, and database operations within a distributed architecture.

### Architecture Overview

![image](https://github.com/user-attachments/assets/7c4891a5-e209-4529-bc18-5eb29f72bcc1)

### Key Components

- **Broker Service**: Entry point for REST requests, handles routing to other services.
- **Task Service**: Handle task creation and modification for a transaction.
- **Transaction Service**: Processes transactions and publishes messages through RabbitMQ.
- **Logger Service**: Listens to logs via RabbitMQ and stores them in MongoDB.
- **Mail Service**: Listens for events via RabbitMQ and handles email sending.
- **Postgres**: Main relational database.
- **MongoDB**: Stores logs for auditing or debugging purposes.
- **RabbitMQ**: Message broker for event-driven communication between services.

---
