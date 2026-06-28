GCP_REGION     := australia-southeast2
GCP_PROJECT    := leecoach-prod
IMAGE          := $(GCP_REGION)-docker.pkg.dev/$(GCP_PROJECT)/leetcoach/backend:latest
CLOUD_SQL_INST := leecoach-prod:$(GCP_REGION):leetcoach-db

.PHONY: dev backend frontend test install-hooks update-adrs deploy

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

# Build, push, and deploy the backend to Cloud Run
deploy:
	docker build --platform linux/amd64 -t $(IMAGE) ./backend
	docker push $(IMAGE)
	gcloud run deploy leetcoach-backend \
		--image=$(IMAGE) \
		--region=$(GCP_REGION) \
		--platform=managed \
		--service-account=leetcoach-cloud-run-sa@$(GCP_PROJECT).iam.gserviceaccount.com \
		--add-cloudsql-instances=$(CLOUD_SQL_INST) \
		--env-vars-file=backend/cloudrun-env.yaml \
		--set-secrets="ANTHROPIC_API_KEY=ANTHROPIC_API_KEY:latest,DB_PASSWORD=DB_PASSWORD:latest" \
		--allow-unauthenticated \
		--project=$(GCP_PROJECT)
