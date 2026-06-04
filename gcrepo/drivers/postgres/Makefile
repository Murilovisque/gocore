CONTAINER_NAME := gc-postgres
POSTGRES_IMG   := docker.io/library/postgres:17
DB_USER        := postgres
DB_PASSWORD    := password
DB_NAME        := postgres

.PHONY: db-up db-down test

db-up:
	@echo "Checking if database is running..."
	@if [ $$(podman ps -q -f name=$(CONTAINER_NAME)) ]; then \
		echo "Database is already up."; \
	else \
		echo "Starting a fresh PostgreSQL container..."; \
		podman run --name $(CONTAINER_NAME) -d \
			-e POSTGRES_USER=$(DB_USER) \
			-e POSTGRES_PASSWORD=$(DB_PASSWORD) \
			-e POSTGRES_DB=$(DB_NAME) \
			--rm \
			-p 5432:5432 \
			$(POSTGRES_IMG); \
		echo "Waiting for DB to initialize..."; \
		sleep 3; \
	fi

db-down:
	@echo "Stopping (and removing) the container..."
	@podman stop $(CONTAINER_NAME) > /dev/null 2>&1 || echo "Container already stopped."

test: db-up
	@go test ./...
