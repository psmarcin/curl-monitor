generate:
	sqlc generate

migrate:
	migrate -database postgresql://user:pass@localhost:5432/job?sslmode=disable -path migration up

migrate-rollback:
	migrate -database postgresql://user:pass@localhost:5432/job?sslmode=disable -path migration down 1
