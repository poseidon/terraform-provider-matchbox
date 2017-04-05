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
