EXTENSION ?= 
DIST_DIR ?= dist/
GOOS ?= linux
ARCH ?= $(shell uname -m)
BUILDINFOSDET ?= 

SOFT_NAME    := sacrebleu
SOFT_VERSION := $(shell git describe --tags $(git rev-list --tags --max-count=1))
VERSION_PKG   := $(shell echo $(SOFT_VERSION) | sed 's/^v//g')
ARCH          := x86_64
LICENSE       := AGPL-3
URL           := https://github.com/outout14/sacrebleu-dns/
DESCRIPTION   := Sacrebleu is a DNS authoritative nameserver made in Go
BUILDINFOS    :=  ($(shell date +%FT%T%z)$(BUILDINFOSDET))
LDFLAGS       := '-X main.version=$(SOFT_VERSION) -X main.buildinfos=$(BUILDINFOS)'

OUTPUT_SOFT := $(DIST_DIR)sacrebleu-$(SOFT_VERSION)-$(GOOS)-$(ARCH)$(EXTENSION)

.PHONY: vet
vet:
	go vet main.go

.PHONY: prepare
prepare:
	mkdir -p $(DIST_DIR)

.PHONY: clean
clean:
	rm -rf $(DIST_DIR)

.PHONY: build
build: prepare
	go build -ldflags $(LDFLAGS) -o $(OUTPUT_SOFT)

.PHONY: package-deb
package-deb: prepare
	fpm -s dir -t deb -n $(SOFT_NAME) -v $(VERSION_PKG) \
        --description "$(DESCRIPTION)"  \
        --url "$(URL)" \
        --architecture $(ARCH) \
        --license "$(LICENSE)" \
        --package $(DIST_DIR) \
        $(OUTPUT_SOFT)=/usr/bin/sacrebleu-dns

.PHONY: package-rpm
package-rpm: prepare
	fpm -s dir -t rpm -n $(SOFT_NAME) -v $(VERSION_PKG) \
	--description "$(DESCRIPTION)" \
	--url "$(URL)" \
	--architecture $(ARCH) \
	--license "$(LICENSE) "\
	--package $(DIST_DIR) \
	$(OUTPUT_SOFT)=/usr/bin/sacrebleu-dns
