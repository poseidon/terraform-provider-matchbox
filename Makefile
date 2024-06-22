export CGO_ENABLED:=0

VERSION=$(shell git describe --tags --match=v* --always)
SEMVER=$(shell git describe --tags --match=v* --always | cut -c 2-)

.PHONY: all
all: build test vet fmt

.PHONY: build
build: clean bin/terraform-provider-matchbox

bin/terraform-provider-matchbox:
	@go build -o $@ -v github.com/poseidon/terraform-provider-matchbox

.PHONY: test
test:
	@go test ./... -cover

.PHONY: vet
vet:
	@go vet -all ./...

.PHONY: fmt
fmt:
	@test -z $$(go fmt ./...)

.PHONY: clean
clean:
	@rm -rf _output

.PHONY: release
release: \
	clean \
	_output/plugin-linux-amd64.zip \
	_output/plugin-linux-arm64.zip \
	_output/plugin-darwin-amd64.zip \
	_output/plugin-darwin-arm64.zip \
	_output/plugin-windows-amd64.zip

_output/plugin-%.zip: NAME=terraform-provider-matchbox_$(SEMVER)_$(subst -,_,$*)
_output/plugin-%.zip: DEST=_output/$(NAME)
_output/plugin-%.zip: LOCAL=$(HOME)/.terraform.d/plugins/terraform.localhost/poseidon/matchbox/$(SEMVER)
_output/plugin-%.zip: _output/%/terraform-provider-matchbox
	@mkdir -p $(DEST)
	@cp _output/$*/terraform-provider-matchbox $(DEST)/terraform-provider-matchbox_$(VERSION)
	@zip -j $(DEST).zip $(DEST)/terraform-provider-matchbox_$(VERSION)
	@mkdir -p $(LOCAL)/$(subst -,_,$*)
	@cp _output/$*/terraform-provider-matchbox $(LOCAL)/$(subst -,_,$*)/terraform-provider-matchbox_$(VERSION)

_output/linux-amd64/terraform-provider-matchbox: GOARGS = GOOS=linux GOARCH=amd64
_output/linux-arm64/terraform-provider-matchbox: GOARGS = GOOS=linux GOARCH=arm64
_output/darwin-amd64/terraform-provider-matchbox: GOARGS = GOOS=darwin GOARCH=amd64
_output/darwin-arm64/terraform-provider-matchbox: GOARGS = GOOS=darwin GOARCH=arm64
_output/windows-amd64/terraform-provider-matchbox: GOARGS = GOOS=windows GOARCH=amd64
_output/%/terraform-provider-matchbox:
	$(GOARGS) go build -o $@ github.com/poseidon/terraform-provider-matchbox

release-sign:
	cd _output; sha256sum *.zip > terraform-provider-matchbox_$(SEMVER)_SHA256SUMS
	gpg --default-key 0x8F515AD1602065C8 --detach-sign _output/terraform-provider-matchbox_$(SEMVER)_SHA256SUMS

release-verify: NAME=_output/terraform-provider-matchbox
release-verify:
	gpg --verify $(NAME)_$(SEMVER)_SHA256SUMS.sig $(NAME)_$(SEMVER)_SHA256SUMS

.PHONY: certificates
certificates:
	@openssl req -days 3650 -nodes -x509 -config matchbox/testdata/certs.ext -extensions v3_ca -newkey rsa:4096 -keyout matchbox/testdata/ca.key -out matchbox/testdata/ca.crt -subj "/CN=fake-ca"
	@openssl req -nodes -newkey rsa:2048 -keyout matchbox/testdata/server.key -out matchbox/testdata/server.csr -subj "/CN=fake-server"
	@openssl x509 -days 3650 -sha256 -extfile matchbox/testdata/certs.ext -extensions v3_server -req -in matchbox/testdata/server.csr -CA matchbox/testdata/ca.crt -CAkey matchbox/testdata/ca.key -CAcreateserial -out matchbox/testdata/server.crt
	@openssl req -nodes -newkey rsa:2048 -keyout matchbox/testdata/client.key -out matchbox/testdata/client.csr -subj "/CN=fake-client"
	@openssl x509 -days 3650 -sha256 -extfile matchbox/testdata/certs.ext -extensions v3_client -req -in matchbox/testdata/client.csr -CA matchbox/testdata/ca.crt -CAkey matchbox/testdata/ca.key -CAserial matchbox/testdata/ca.srl -out matchbox/testdata/client.crt
	@rm matchbox/testdata/*.csr matchbox/testdata/ca.srl
