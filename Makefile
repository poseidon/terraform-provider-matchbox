export CGO_ENABLED:=0
export GO111MODULE=on
export GOFLAGS=-mod=vendor

VERSION=$(shell ./scripts/git-version)
GOPATH_BIN:=$(shell echo ${GOPATH} | awk 'BEGIN { FS = ":" }; { print $1 }')/bin

.PHONY: all
all: bin/terraform-provider-matchbox

.PHONY: build
bin/terraform-provider-matchbox:
	@go build -o $@ -v github.com/coreos/terraform-provider-matchbox

.PHONY: install
install: bin/terraform-provider-matchbox
	@cp $< $(GOPATH_BIN)

.PHONY: test
test:
	@./scripts/test

.PHONY: update
update:
	@GOFLAGS="" go get -u
	@go mod tidy

.PHONY: vendor
vendor:
	@go mod vendor

.PHONY: clean
clean:
	@rm -rf _output

.PHONY: release
release: \
	clean \
	_output/plugin-linux-amd64.tar.gz \
	_output/plugin-darwin-amd64.tar.gz \
	_output/plugin-windows-amd64.tar.gz

_output/plugin-%.tar.gz: NAME=terraform-provider-matchbox-$(VERSION)-$*
_output/plugin-%.tar.gz: DEST=_output/$(NAME)
_output/plugin-%.tar.gz: _output/%/terraform-provider-matchbox
	@mkdir -p $(DEST)
	@cp _output/$*/terraform-provider-matchbox $(DEST)
	@tar zcvf $(DEST).tar.gz -C _output $(NAME)

_output/linux-amd64/terraform-provider-matchbox: GOARGS = GOOS=linux GOARCH=amd64
_output/darwin-amd64/terraform-provider-matchbox: GOARGS = GOOS=darwin GOARCH=amd64
_output/windows-amd64/terraform-provider-matchbox: GOARGS = GOOS=windows GOARCH=amd64
_output/%/terraform-provider-matchbox:
	$(GOARGS) go build -o $@ github.com/coreos/terraform-provider-matchbox

.PHONY: certificates
certificates:
	@openssl req -days 3650 -nodes -x509 -config matchbox/testdata/certs.ext -extensions v3_ca -newkey rsa:4096 -keyout matchbox/testdata/ca.key -out matchbox/testdata/ca.crt -subj "/CN=fake-ca"
	@openssl req -nodes -newkey rsa:2048 -keyout matchbox/testdata/server.key -out matchbox/testdata/server.csr -subj "/CN=fake-server"
	@openssl x509 -days 365 -sha256 -extfile matchbox/testdata/certs.ext -extensions v3_server -req -in matchbox/testdata/server.csr -CA matchbox/testdata/ca.crt -CAkey matchbox/testdata/ca.key -CAcreateserial -out matchbox/testdata/server.crt
	@openssl req -nodes -newkey rsa:2048 -keyout matchbox/testdata/client.key -out matchbox/testdata/client.csr -subj "/CN=fake-client"
	@openssl x509 -days 365 -sha256 -extfile matchbox/testdata/certs.ext -extensions v3_client -req -in matchbox/testdata/client.csr -CA matchbox/testdata/ca.crt -CAkey matchbox/testdata/ca.key -CAserial matchbox/testdata/ca.srl -out matchbox/testdata/client.crt
	@rm matchbox/testdata/*.csr matchbox/testdata/ca.srl
