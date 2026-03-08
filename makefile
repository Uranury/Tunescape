.PHONY: up down build restart logs migrate

up:
	docker compose up -d

down:
	docker compose down

build:
	docker compose up -d --build

restart:
	docker compose restart api

logs:
	docker compose logs -f api

migrate:
	docker compose exec api ./app migrate

fmt:
	goimports -w .