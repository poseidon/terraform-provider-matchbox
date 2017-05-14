VERSION=$(shell ./scripts/git-version)

.PHONY: all
all: build

.PHONY: build
build:
	@go install -v github.com/coreos/terraform-provider-matchbox

.PHONY: test
test:
	@./scripts/test

.PHONY: vendor
vendor:
	@glide update --strip-vendor
	@glide-vc --use-lock-file --no-tests --only-code

.PHONY: clean
clean:
	@rm -rf _output

.PHONY: release
release: \
	clean \
	_output/plugin-linux-amd64.tar.gz \
	_output/plugin-darwin-amd64.tar.gz

_output/plugin-%.tar.gz: NAME=terraform-provider-matchbox-$(VERSION)-$*
_output/plugin-%.tar.gz: DEST=_output/$(NAME)
_output/plugin-%.tar.gz: _output/%/terraform-provider-matchbox
	@mkdir -p $(DEST)
	@cp _output/$*/terraform-provider-matchbox $(DEST)
	@tar zcvf $(DEST).tar.gz -C _output $(NAME)

_output/linux-amd64/terraform-provider-matchbox: GOARGS = GOOS=linux GOARCH=amd64
_output/darwin-amd64/terraform-provider-matchbox: GOARGS = GOOS=darwin GOARCH=amd64
_output/%/terraform-provider-matchbox:
	$(GOARGS) go build -o $@ github.com/coreos/terraform-provider-matchbox
