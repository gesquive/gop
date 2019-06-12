#
#  Makefile
#
#  The kickoff point for all project management commands.
#

GOCC := go

# Program version
VERSION := $(shell git describe --always --tags)

# Binary name for bintray
BIN_NAME=gop

# Project owner for bintray
OWNER=gesquive

# Project name for bintray
PROJECT_NAME=gop

# Project url used for builds
# examples: github.com, bitbucket.org
REPO_HOST_URL=github.com

# Grab the current commit
GIT_COMMIT=$(shell git rev-parse HEAD)

# Check if there are uncommited changes
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)

# Use a local vendor directory for any dependencies; comment this out to
# use the global GOPATH instead
# GOPATH=$(PWD)

INSTALL_PATH=$(GOPATH)/src/${REPO_HOST_URL}/${OWNER}/${PROJECT_NAME}

PKG_LIST := ./...
PKG_TAR=cd dist/${PKG} && tar --exclude=".*" --owner=0 --group=0 -zcf ../${PROJECT_NAME}-${VERSION}-${PKG}.tar.gz *
PKG_ZIP=cd dist/${PKG} && zip --exclude .\* -qr ../${PROJECT_NAME}-${VERSION}-${PKG}.zip *
PKG_DST=cd dist && find . -mindepth 1 -maxdepth 1 -type d -exec

default: test build

.PHONY: help
help:
	@echo 'Management commands for ${PROJECT_NAME}:'
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
	 awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2}'
	@echo

.PHONY: build
build: ## Compile the project
	@echo "building ${OWNER} ${BIN_NAME} ${VERSION}"
	@echo "GOPATH=${GOPATH}"
	${GOCC} build -ldflags "-X main.version=${VERSION} -X main.dirty=${GIT_DIRTY}" -o ${BIN_NAME}

.PHONY: install
install: build ## Install the binaries on this computer
	install -d ${DESTDIR}/usr/local/bin/
	install -m 755 ./${BIN_NAME} ${DESTDIR}/usr/local/bin/${BIN_NAME}

.PHONY: deps
deps: ## Download project dependencies
	${GOCC} mod download

.PHONY: test
test: ## Run golang tests
	${GOCC} test ${PKG_LIST}

.PHONY: bench
bench: ## Run golang benchmarks
	${GOCC} test -benchmem -bench=. ${PKG_LIST}

.PHONY: cover
cover: ## Run coverage report
	${GOCC} test -v -cover ${PKG_LIST}

.PHONY: clean
clean: ## Clean the directory tree of artifacts
	${GOCC} clean
	rm -f ./${BIN_NAME}.test
	rm -f ./${BIN_NAME}
	rm -rf ./dist

.PHONY: build-dist
build-dist: gox
	gox -verbose \
	-ldflags "-X main.version=${VERSION} -X main.dirty=${GIT_DIRTY}" \
	-os="linux darwin windows" \
	-arch="amd64 386" \
	-output="dist/{{.OS}}-{{.Arch}}/{{.Dir}}" .

.PHONY: package-dist
package-dist: gop
	gop --delete \
	--os="linux darwin" \
	--arch="amd64 386" \
	--archive="tar.gz" \
	--packages="windows/amd64/zip windows/386/zip" \
	--files="LICENSE README.md" \
	--input="dist/{{.OS}}-{{.Arch}}/{{.Dir}}" \
	--output="dist/{{.Dir}}-${VERSION}-{{.OS}}-{{.Arch}}.{{.Archive}}" .

.PHONY: dist
dist: build-dist package-dist ## Cross compile and package the full distribution

.PHONY: fmt
fmt: ## Reformat the source tree with gofmt
	find . -name '*.go' -not -path './.vendor/*' -exec gofmt -w=true {} ';'

.PHONY: link
link: $(INSTALL_PATH) ## Symlink this project into the GOPATH
$(INSTALL_PATH):
	@mkdir -p `dirname $(INSTALL_PATH)`
	@ln -s $(PWD) $(INSTALL_PATH) >/dev/null 2>&1

.PHONY: gox
gox:
	@command -v gox >/dev/null 2>&1 || \
	echo "Installing gox" && ${GOCC} get -u github.com/mitchellh/gox

.PHONY: gop
gop:
	@command -v gop >/dev/null 2>&1 || \
	echo "Installing gop" && ${GOCC} get -u github.com/gesquive/gop
	@gop --version
