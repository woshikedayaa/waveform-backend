.PHONY: build clean run windows

BINARY=wf

SRC_DIR=.
DIST_DIR=./dist

BUILD_ARCH=arm arm64 386 amd64 ppc64le riscv64 \
	mips mips64le mipsle loong64 s390x
BUILD_ARGS=-trimpath -ldflags="-s -w"

build: clean windows $(BUILD_ARCH)
$(BUILD_ARCH):
	@echo "Building Linux $@ ..."
	@mkdir -p $(DIST_DIR)/$@
	@rm -rf $(DIST_DIR)/$@/*
	@CGO_ENABLED=0 GOOS=linux GOARCH=$@ GOMIPS=softfloat go build \
		$(BUILD_ARGS) -o $(DIST_DIR)/$@/$(BINARY) $(SRC_DIR)/*.go


windows:
	@echo "Building Windows 32-bit & 64-bit ..."
	@mkdir -p $(DIST_DIR)/win32 $(DIST_DIR)/win64
	@rm -rf $(DIST_DIR)/win32/* $(DIST_DIR)/win64/*
	@CGO_ENABLED=0 GOOS=windows GOARCH=386 go build \
		$(BUILD_ARGS) -o $(DIST_DIR)/win32/$(BINARY).exe $(SRC_DIR)/*.go
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build \
		$(BUILD_ARGS) -o $(DIST_DIR)/win64/$(BINARY).exe $(SRC_DIR)/*.go

run:
	@go run $(SRC_DIR)/*.go

clean:
	@rm -rf $(DIST_DIR)/*
