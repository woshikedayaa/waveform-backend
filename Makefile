.PHONY: deploy develop init default windows run clean docker

.DEFAULT_GOAL=default

VERSION=0.0.1
BINARY=wf
SRC_DIR=.
DIST_DIR=./bin
ALL_ARCH=arm arm64 386 amd64 ppc64le riscv64 \
	mips mips64le mipsle loong64 s390x
BUILD_ARGS=-trimpath -ldflags="-s -w -X main.Version=$(VERSION)"
DEPLOY_TAGS="deploy"
DEVELOP_TAGS="develop"

# 这个是为了发布用的 可以构建全部的架构 可以发布到 release
deploy: clean $(ALL_ARCH) windows
develop: clean
	@CGO_ENABLED=0 go build \
     		$(BUILD_ARGS) -tags $(DEVELOP_TAGS) -o $(DIST_DIR)/$(BINARY)_$(VERSION) $(SRC_DIR)/*.go

init:
	@mkdir -p $(DIST_DIR)

default: clean
	@echo "Building ..."
	@rm -rf $(DIST_DIR)/*
	@CGO_ENABLED=0 GOMIPS=softfloat go build \
		$(BUILD_ARGS) -tags $(DEPLOY_TAGS) -o $(DIST_DIR)/$(BINARY)_$(VERSION) $(SRC_DIR)/*.go

$(ALL_ARCH):
	@echo "Building Linux $@ ..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=$@ GOMIPS=softfloat go build \
		$(BUILD_ARGS) -tags $(DEPLOY_TAGS) -o $(DIST_DIR)/$(BINARY)_$(VERSION)_linux_$@ $(SRC_DIR)/*.go

windows:
	@echo "Building Windows 32-bit & 64-bit ..."
	@CGO_ENABLED=0 GOOS=windows GOARCH=386 go build \
		$(BUILD_ARGS) -tags $(DEPLOY_TAGS) -o $(DIST_DIR)/$(BINARY)_$(VERSION)_windows_i386.exe  $(SRC_DIR)/*.go
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build \
		$(BUILD_ARGS) -tags $(DEPLOY_TAGS) -o $(DIST_DIR)/$(BINARY)_$(VERSION)_windows_amd64.exe $(SRC_DIR)/*.go

run:
	@go run -tags $(DEVELOP_TAGS) $(SRC_DIR)/*.go

clean: init
	@rm -rf $(DIST_DIR)/*

docker:
	@docker build -t waveform-backend:$(VERSION) . -f Dockerfile

