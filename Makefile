# Dependencies
dependencies: dependencies-common dependencies-command-run dependencies-job dependencies-result dependencies-trigger

dependencies-command-run:
	cd service/command-run && go mod download
dependencies-job:
	cd service/job && go mod download
dependencies-result:
	cd service/result && go mod download
dependencies-trigger:
	cd service/trigger && go mod download
dependencies-common:
	cd service/common && go mod download

# Build
build: build-command-run build-job build-result build-trigger

build-command-run:
	cd service/command-run && go build -o ./../../cmd/command-run
build-job:
	cd service/job && go build -o ./../../cmd/job
build-result:
	cd service/result && go build -o ./../../cmd/result
build-trigger:
	cd service/trigger && go build -o ./../../cmd/trigger

# Generate
generate:
	sqlc generate

docker-image:
	docker build -t curl-monitor:1.3.2 .

# Migrations
migrate:
	migrate -database postgresql://user:pass@localhost:5432/job?sslmode=disable -path migration up

migrate-rollback:
	migrate -database postgresql://user:pass@localhost:5432/job?sslmode=disable -path migration down 1

# Infrastructure
infra: infra-database infra-job infra-trigger infra-result infra-command-run

infra-database:
	helm repo add bitnami https://charts.bitnami.com/bitnami
	helm upgrade --install postgres bitnami/postgresql --values ./infrastructure/database.yaml

infra-message-broker:
	helm repo add bitnami https://charts.bitnami.com/bitnami
	helm upgrade --install rabbitmq bitnami/rabbitmq --values ./infrastructure/rabbitmq.yaml --version 8.20.1

infra-job:
	helm upgrade --install curl-monitor-job ./infrastructure/base --values ./infrastructure/job.yaml

infra-trigger:
	helm upgrade --install curl-monitor-trigger ./infrastructure/base --values ./infrastructure/trigger.yaml

infra-result:
	helm upgrade --install curl-monitor-result ./infrastructure/base --values ./infrastructure/result.yaml

infra-command-run:
	helm upgrade --install curl-monitor-command-run ./infrastructure/base --values ./infrastructure/command-run.yaml
