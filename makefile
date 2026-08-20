include .env
export $(shell sed 's/=.*//' .env)

.PHONY: grpc_gen model_gen docker_down docker_up docker_rebuild migrate_create

# gRPC 代码生成
grpc_gen:
	goctl rpc protoc $(proto) \
		--go_out=. \
		--go-grpc_out=. \
		--zrpc_out=. \
		--module=go-zero-flash-sale

# MySQL Model 代码生成
model_gen:
	goctl model mysql ddl -src ./deploy/sql/$(name) -dir ./deploy/models -cache

# Docker 相关
docker_down:
	docker compose down -v
docker_up:
	docker compose --env-file .env up -d
docker_rebuild:
	docker compose down -v
	docker compose --env-file .env up -d --build

# Migration 相关
migrate_create:
	migrate create -ext sql -dir ./deploy/migration -seq $(name)
migrate_up:
	migrate -path ./deploy/migration -database "$(DB_URL)" -verbose up
