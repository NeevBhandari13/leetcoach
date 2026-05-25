.PHONY: dev backend frontend test

# Start both services in parallel. Ctrl-C kills both.
dev:
	@trap 'kill 0' SIGINT; \
	$(MAKE) backend & \
	$(MAKE) frontend & \
	wait

backend:
	cd backend && go run ./cmd/main.go

frontend:
	cd frontend && npm run dev

test:
	cd backend && go test ./...
