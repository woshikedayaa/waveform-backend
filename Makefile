.PHONY: deploy develop init default windows run clean docker install uninstall reinstall

.DEFAULT_GOAL=default
#
VERSION=0.1
# build
BINARY=waveform

# 这里还有好多架构不支持 例如 mips全系列 loong64
ALL_ARCH=arm arm64 386 amd64 ppc64le riscv64 s390x
BUILD_ARGS=-trimpath -ldflags="-s -w -X main.Version=$(VERSION)"
DEPLOY_TAGS="deploy"
DEVELOP_TAGS="develop"
# dir
SRC_DIR=.
DIST_DIR=./bin
INSTALL_DIR=/usr/local/bin
CONFIG_DIR=/usr/local/etc/waveform
LOG_DIR=/var/log/waveform
LIB_DIR=/var/lib/waveform

# 这个是为了发布用的 可以构建全部的架构 可以发布到 release
deploy: clean $(ALL_ARCH)
develop: clean
	@CGO_ENABLED=0 GOMIPS=softfloat go build \
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

# 这个编译不过 因为sqlite 的不支持
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
	@docker build -t $(BINARY):$(VERSION) . -f Dockerfile

install:
	mkdir -p $(LOG_DIR) $(CONFIG_DIR) $(LIB_DIR)
	cp $(SRC_DIR)/config/config_full.yaml $(CONFIG_DIR)/config.yaml
	CGO_ENABLED=0 GOMIPS=softfloat go build \
		$(BUILD_ARGS) -tags $(DEPLOY_TAGS) -o /tmp/$(BINARY) $(SRC_DIR)/*.go
	install -m 0755 -o root -g root -T /tmp/$(BINARY) $(INSTALL_DIR)/$(BINARY)
	rm -f /tmp/$(BINARY)
	@echo "Install success"

uninstall:
	rm -f $(INSTALL_DIR)/$(BINARY)
	rm -rf $(CONFIG_DIR)
	@echo "Uninstall success"


reinstall: uninstall install
