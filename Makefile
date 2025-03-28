.PHONY: all init server front clean

all: server front

init:
	@echo "Starting server..."
	@cd server && go run cmd/main.go populate

server:
	@echo "Starting server..."
	@cd server && /home/nic/go/bin/air

front:
	@echo "Setting up front-end..."
	@cd front && npm run dev

clean:
	@echo "Cleaning up..."
	@rm -f server/db/forum.db
	@cd front && npm run clean
	@cd front && npm install
