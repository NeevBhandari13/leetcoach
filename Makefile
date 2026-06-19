.PHONY: dev backend frontend test install-hooks update-adrs

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

# Install the post-commit git hook (run once after cloning)
install-hooks:
	@printf '#!/usr/bin/env bash\n"$$(git rev-parse --show-toplevel)/scripts/update-adrs.sh"\n' > .git/hooks/post-commit
	@chmod +x .git/hooks/post-commit
	@echo "Hook installed: .git/hooks/post-commit"

# Run ADR analysis on the last commit manually
update-adrs:
	@bash scripts/update-adrs.sh
