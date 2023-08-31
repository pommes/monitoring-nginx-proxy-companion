.PHONY: .assert_colima_runs .assert_mybuilder_exists .assert_golang_exists .create_git_tag

DOCKERX_BUILDER ?= mybuilder
ARCH_TARGET ?= linux/amd64
IMAGE ?= registry.komm.link/base/docker/nginx-proxy-metrics
TAG ?= latest
NEXT_VERSION ?= 2023-6

GHCR_IMAGE ?= ghcr.io/pommes/nginx-proxy-metrics
GHCR_NEXT_VERSION ?= v1.0.0

# RUN TARGETS ##############

build_go: .assert_golang_exists
	@echo "=== Building main"
	go mod tidy
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./.build/main ./cmd

build_image_latest: build_go .assert_colima_runs .assert_mybuilder_exists
	@echo "=== Building image: $(IMAGE):$(TAG) ($(ARCH_TARGET))"
	export DOCKER_CLI_EXPERIMENTAL=enabled
	@docker buildx version
	docker buildx build --platform $(ARCH_TARGET) -t $(IMAGE):$(TAG) . --load
	docker push ${IMAGE}:${TAG}

build_image_latest_ghcr:
	export IMAGE=$(GHCR_IMAGE) && \
	export TAG=latest && \
	$(MAKE) build_image_latest

build_image_tag: .create_git_tag
	$(MAKE) build_image_latest && \
	export TAG=$(NEXT_VERSION) && \
	$(MAKE) build_image_latest

release_ghcr:
	export NEXT_VERSION=$(GHCR_NEXT_VERSION) && \
	$(MAKE) .create_git_tag && \
	git push origin $(GHCR_NEXT_VERSION)


# LIB TARGETS ##############

.create_git_tag:
	@if git tag -a $(NEXT_VERSION) -m "tag: $(NEXT_VERSION)"; then \
		echo "--- git tag '$(NEXT_VERSION)' not found. Creating tag ..."; \
	else \
		echo "Could not create git tag '$(NEXT_VERSION)'"; \
		exit 1; \
	fi

.assert_golang_exists:
	@if which go 1>/dev/null 2>&1; then \
		echo "--- go found."; \
	else \
		echo "!!! go not found. Intalling via homebrew..."; \
		brew install golang; \
	fi

.assert_colima_runs:
	@if colima status 1>/dev/null 2>&1; then \
		echo "--- colima is already running."; \
	else \
		echo "!!! colima is not running. Starting colima..."; \
		colima start; \
	fi

.assert_mybuilder_exists: .assert_colima_runs
	@if docker buildx inspect ${DOCKERX_BUILDER} 1>/dev/null 2>&1; then \
		echo "--- buildx builder '${DOCKERX_BUILDER}' found."; \
    else \
    	echo "!!! buildx builder '${DOCKERX_BUILDER}' not found. Creating builder..."; \
    	docker buildx create --use --name ${DOCKERX_BUILDER}; \
		docker buildx inspect ${DOCKERX_BUILDER} --bootstrap; \
    fi


