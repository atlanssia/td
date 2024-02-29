# 定义全局变量和工具路径
SHELL := /bin/bash
GOPATH ?= $(shell go env GOPATH)
GOBIN  ?= $(GOPATH)/bin
REPO_PATH := $(shell git rev-parse --show-toplevel)
BUILD_DIR := build
COMMIT_HASH := $(shell git rev-parse --short HEAD)
COMMIT_TIME := $(shell git show -s --format=%cI)
BUILD_TIME := $(shell date +"%Y-%m-%dT%H:%M:%S%z")
GO_VERSION := $(shell go version)
PKG_PREFIX := github.com/atlanssia/td
GET_LATEST_TAG := $(shell git describe --tags `git rev-list --tags --max-count=1`)
VERSION := $(GET_LATEST_TAG)

# 检查go环境
.PHONY: check-go-env
check-go-env:
ifeq (,$(shell which go))
	$(error "Go is not installed or not in PATH")
endif

# 编译项目
.PHONY: build
build: check-go-env
	GOOS=linux GOARCH=amd64 go build -ldflags "-X '$(PKG_PREFIX)/version.goVersion=$(GO_VERSION)' \
		-X '$(PKG_PREFIX)/version.buildVersion=$(COMMIT_HASH)' \
		-X '$(PKG_PREFIX)/version.version=$(VERSION)' \
		-X '$(PKG_PREFIX)/version.lastCommitTime=$(COMMIT_TIME)' \
		-X '$(PKG_PREFIX)/version.buildTime=$(BUILD_TIME)' \
		-X '$(PKG_PREFIX)/version.goos=linux' \
		-X '$(PKG_PREFIX)/version.goarch=amd64' \
		-s -w" -o $(BUILD_DIR)/td-linux-amd64 cmd/td/main.go
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X '$(PKG_PREFIX)/version.goVersion=$(GO_VERSION)' \
		-X '$(PKG_PREFIX)/version.buildVersion=$(COMMIT_HASH)' \
		-X '$(PKG_PREFIX)/version.version=$(VERSION)' \
		-X '$(PKG_PREFIX)/version.lastCommitTime=$(COMMIT_TIME)' \
		-X '$(PKG_PREFIX)/version.buildTime=$(BUILD_TIME)' \
		-X '$(PKG_PREFIX)/version.goos=darwin' \
		-X '$(PKG_PREFIX)/version.goarch=amd64' \
		-s -w" -o $(BUILD_DIR)/td-darwin-amd64 cmd/td/main.go
	GOOS=windows GOARCH=amd64 go build -ldflags "-X '$(PKG_PREFIX)/version.goVersion=$(GO_VERSION)' \
		-X '$(PKG_PREFIX)/version.buildVersion=$(COMMIT_HASH)' \
		-X '$(PKG_PREFIX)/version.version=$(VERSION)' \
		-X '$(PKG_PREFIX)/version.lastCommitTime=$(COMMIT_TIME)' \
		-X '$(PKG_PREFIX)/version.buildTime=$(BUILD_TIME)' \
		-X '$(PKG_PREFIX)/version.goos=windows' \
		-X '$(PKG_PREFIX)/version.goarch=amd64' \
		-s -w" -o $(BUILD_DIR)/td-windows-amd64.exe cmd/td/main.go

.PHONY: copy-configs
copy-configs:
	mkdir -p $(BUILD_DIR)/configs
	cp -r configs/* $(BUILD_DIR)/configs
# 打包
.PHONY: package-targz
package-targz: build copy-configs
	mkdir -p $(BUILD_DIR)/dist
	tar -czf $(BUILD_DIR)/dist/td-linux-amd64-$(VERSION).tar.gz -C $(BUILD_DIR) td-linux-amd64 configs
	tar -czf $(BUILD_DIR)/dist/td-darwin-amd64-$(VERSION).tar.gz -C $(BUILD_DIR) td-darwin-amd64 configs

.PHONY: package-zip
package-zip: build copy-configs
	mkdir -p $(BUILD_DIR)/dist
	(cd $(BUILD_DIR) && zip -r dist/td-linux-amd64-$(VERSION).zip td-linux-amd64 configs)
	(cd $(BUILD_DIR) && zip -r dist/td-darwin-amd64-$(VERSION).zip td-darwin-amd64 configs)
	(cd $(BUILD_DIR) && zip -r dist/td-windows-amd64-$(VERSION).zip td-windows-amd64.exe configs)

# 同时打包为tar.gz和zip
.PHONY: package
package: package-targz package-zip

# 测试
.PHONY: test
test: check-go-env
	go test ./... -coverprofile=coverage.out

# Lint检查（假设已安装golangci-lint）
.PHONY: lint
lint: check-go-env
	golangci-lint run

# 清理
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR) coverage.out

# 获取Git Commit ID为版本号后缀
ifndef COMMIT_HASH
	$(warning "Warning: Git not found, version will not include commit hash")
endif

.DEFAULT_GOAL := build
