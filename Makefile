.PHONY: all server front clean

all: server front

server:
	@echo "Starting server..."
	@cd server && go run cmd/main.go

front:
	@echo "Setting up front-end..."
	@cd front && npm run dev

clean:
	@echo "Cleaning up..."
	@rm -f server/database.db
	@cd front && npm run clean
	@cd front && npm install