.PHONY: fmt fmt-check vet test

fmt:
	gofmt -w .

fmt-check:
	@test -z "$$(gofmt -l .)" || (echo "gofmt needed on:" && gofmt -l . && exit 1)

vet:
	@packages="$$(go list ./... 2>/dev/null)"; \
	if [ -n "$$packages" ]; then \
		go vet $$packages; \
	else \
		echo "no Go packages, skip go vet"; \
	fi

test:
	@packages="$$(go list ./... 2>/dev/null)"; \
	if [ -n "$$packages" ]; then \
		go test $$packages; \
	else \
		echo "no Go packages, skip go test"; \
	fi

.PHONY: up down logs backend db migrate-create migrate-up migrate-down migrate-status

up:
	docker compose up --build

down:
	docker compose down

logs:
	docker compose logs -f

backend:
	docker compose exec backend sh

db:
	docker compose exec db sh

migrate-create:
	docker compose exec backend \
	goose -dir /db/migrations create $(name) sql

migrate-up:
	docker compose exec backend sh -c 'goose -dir /db/migrations postgres "$$DATABASE_URL" up'

migrate-down:
	docker compose exec backend sh -c 'goose -dir /db/migrations postgres "$$DATABASE_URL" down'

migrate-status:
	docker compose exec backend sh -c 'goose -dir /db/migrations postgres "$$DATABASE_URL" status'