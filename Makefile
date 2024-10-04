POSTGRES_CONNECTION_URL = "postgresql://postgres:postgres@127.0.0.1:5432/atc?sslmode=disable"

.PHONY: init run clean kill infra-up infra-down docker-image docker-migrate-image dockerignore-check migration migrate-up migrate-down migrate-version

default: run

# Initialize the project.
init:
	@mkdir -p static/dist
	@touch static/dist/tmp
	@echo "Installing pnpm dependencies.."
	@pnpm install
	@echo "Refreshing go packages.."
	@go mod tidy
	@echo "Initializtion complete. Run 'make run' to run the server."

# Run the server.
run:
	@echo "Killing.."
	@lsof -t -i tcp:3000 | xargs kill -9
	@lsof -t -i tcp:8080 | xargs kill -9
	@echo "Running pnpm.."
	@pnpm dev &
	@echo "Running server.."
	@air -c .air.toml

# Clean the project.
clean:
	@echo "Cleaning tmp directory.."
	@rm -rf tmp
	@echo "Cleaning static dist directory.."
	@rm -rf static/dist

# Kill the server.
kill:
	@echo "Killing.."
	@lsof -t -i tcp:3000 | xargs kill -9
	@lsof -t -i tcp:8080 | xargs kill -9

# Start the infrastructure services.
infra-up:
	@docker-compose up -d

# Stop the infrastructure services. 
infra-down:
	@docker-compose down

# Build the Docker image for the server.
docker-image:
	@echo "Building glut Docker image.."
	@docker build -t glut-debug -f Dockerfile .

# Build the Docker image for the migration tool.
docker-migrate-image:
	@docker build -t glut-migrate-debug -f db/Dockerfile .

# Check which files are included in the Docker build context.
dockerignore-check:
	@rsync -avn . /dev/shm --exclude-from .dockerignore

# Create a new migration file.
migration:
	@if [ -z "$(name)" ]; then \
		echo "migration name required e.g. 'make migration name=hello'"; \
	else \
		goose -dir db create $(name) sql; \
	fi

# Migrate database to the latest version.
migrate-up:
	@goose -dir db postgres "${POSTGRES_CONNECTION_URL}" up

# Migrate database down to the previous version.
migrate-down:
	@goose -dir db postgres "${POSTGRES_CONNECTION_URL}" down

# Report the current migration version of the database.
migrate-version:
	@goose postgres "${POSTGRES_CONNECTION_URL}" version
