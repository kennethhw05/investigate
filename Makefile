PACKAGE   := betting-feed

SRC       := $(shell find . -type f -name "*.go")

DEV_IMAGE := us.gcr.io/estars-228814/betting-feed:dev
QA_IMAGE := us.gcr.io/estars-228814/betting-feed:qa
STAGING_IMAGE := us.gcr.io/estars-228814/betting-feed:staging
PROD_IMAGE := us.gcr.io/estars-228814/betting-feed:prod

SERVICE_PKG_DIR := ./cmd/betting-feed
MIGRATIONS_PKG_DIR := ./migrations

TEST_PATTERN?=.
TEST_OPTIONS?=-race -covermode=atomic -coverprofile=coverage.txt -failfast -p 1 # TODO migrate to different fixture lib with parallel package support

DEPLOYMENT := deployment.apps/betting-feed-app-deployment

define prepare = 
	(docker pull $(1)-latest || true)
	docker build --cache-from $(1)-latest -t $(1)-${CI_COMMIT_SHA} -t $(1)-latest --file=infrastructure/docker/$(2) .
	docker push $(1)-${CI_COMMIT_SHA}
	docker push $(1)-latest
endef

## Local docker
test-debug-docker: ## Run all the tests
	docker-compose -f docker-compose-test-debug.yml up

test-docker: ## Run all the tests
	docker-compose -f docker-compose-test.yml run --rm betting-feed-app-test make test

gqlgen-docker:
	docker-compose -f docker-compose-gen.yml run --rm app make gqlgen

go-bindata-docker:
	docker-compose -f docker-compose-gen.yml run --rm app make go-bindata

## General

setup: ## Download dependencies for the go project
	go mod download

test: ## Run all tests outside a docker container
	go test $(TEST_OPTIONS) $(shell go list ./...) -run $(TEST_PATTERN) -timeout=5m

lint: ## Run all the linters
	go vet ./...

build: ## Build go binary
	go build -o bin/${PACKAGE}-service ${SERVICE_PKG_DIR}

optimized-docker-service-build: ## Create optimized binary for x64-linux
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /bin/${PACKAGE}-service ${SERVICE_PKG_DIR}

prepare-dev-deploy: ## Build container image for dev and push
	$(call prepare,$(DEV_IMAGE),release.Dockerfile)

prepare-qa-deploy: ## Build container image for qa and push
	$(call prepare,$(QA_IMAGE),release.Dockerfile)

prepare-staging-deploy: ## Pulls previous qa image, tags as staging, and pushes to registry
	$(call prepare,$(STAGING_IMAGE),release.Dockerfile)

prepare-prod-deploy: ## Rebuilds and pushes to registry.  Temporary until prod and lowers are in same GCP project.
	$(call prepare,$(PROD_IMAGE),release.Dockerfile)

deploy-dev: ## Replaces image string in dev deployment yaml and applys it with kubectl
	kubectl --record $(DEPLOYMENT) set image $(DEPLOYMENT) betting-feed=${DEV_IMAGE}-${CI_COMMIT_SHA} --namespace=espbetting-dev

deploy-qa: ## Replaces image string in qa deployment yaml and applys it with kubectl
	kubectl --record $(DEPLOYMENT) set image $(DEPLOYMENT) betting-feed=${QA_IMAGE}-${CI_COMMIT_SHA} --namespace=espbetting-qa

deploy-staging: ## Replaces image string in staging deployment yaml and applys it with kubectl
	kubectl --record $(DEPLOYMENT) set image $(DEPLOYMENT) betting-feed=${STAGING_IMAGE}-${CI_COMMIT_SHA} --namespace=espbetting-staging

deploy-prod:
	kubectl --record $(DEPLOYMENT) set image $(DEPLOYMENT) betting-feed=${PROD_IMAGE}-${CI_COMMIT_SHA} --namespace=espbetting-prod

gqlgen:
	GO111MODULE=on go mod download
	GO111MODULE=on go run github.com/99designs/gqlgen

go-bindata:
	GO111MODULE=off go get -u github.com/go-bindata/go-bindata/...
	cd migrations && GO111MODULE=off go-bindata -ignore=migration-data.go -pkg migrations -o migration-data.go . && cd ..

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
