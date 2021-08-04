dependencies: dependencies-command-run dependencies-job dependencies-result dependencies-trigger

dependencies-command-run:
	cd service/command-run && go mod download
dependencies-job:
	cd service/job && go mod download
dependencies-result:
	cd service/result && go mod download
dependencies-trigger:
	cd service/trigger && go mod download

build: build-command-run build-job build-result build-trigger

build-command-run:
	cd service/command-run && go build -o ./../../cmd/command-run
build-job:
	cd service/job && go build -o ./../../cmd/job
build-result:
	cd service/result && go build -o ./../../cmd/result
build-trigger:
	cd service/trigger && go build -o ./../../cmd/trigger

generate:
	sqlc generate

migrate:
	migrate -database postgresql://user:pass@localhost:5432/job?sslmode=disable -path migration up

migrate-rollback:
	migrate -database postgresql://user:pass@localhost:5432/job?sslmode=disable -path migration down 1
