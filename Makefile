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

# Migrations
migrate:
	migrate -database postgresql://user:pass@localhost:5432/job?sslmode=disable -path migration up

migrate-rollback:
	migrate -database postgresql://user:pass@localhost:5432/job?sslmode=disable -path migration down 1

# Infrastructure
infra: infra-database infra-message-broker infra-job infra-trigger infra-result infra-command-run

infra-database:
	helm repo add bitnami https://charts.bitnami.com/bitnami
	helm upgrade --install --values infrastructure/database.yaml postgres bitnami/postgresql

infra-message-broker:
	helm repo add bitnami https://charts.bitnami.com/bitnami
	helm upgrade --install --values ./infrastructure/rabbitmq.yaml rabbitmq bitnami/rabbitmq --version 8.20.1

infra-job:
	helm upgrade --install --debug curl-monitor-job ./infrastructure/base -f ./infrastructure/job.yaml

infra-trigger:
	helm upgrade --install --debug curl-monitor-trigger ./infrastructure/base -f ./infrastructure/trigger.yaml

infra-result:
	helm upgrade --install --debug curl-monitor-result ./infrastructure/base -f ./infrastructure/result.yaml

infra-command-run:
	helm upgrade --install --debug curl-monitor-command-run ./infrastructure/base -f ./infrastructure/command-run.yaml
