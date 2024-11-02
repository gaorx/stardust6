
# help
.PHONY: help
help:
	@echo 'make fmt|test|test-clean|upgrade-deps'

# fmt
.PHONY: fmt
fmt:
	@bash _scripts/fmt.sh

# test
.PHONY: test
test:
	@bash _scripts/test.sh

# test-clean
.PHONE: test-clean
test-clean:
	@go clean -testcache && echo "Go test cache cleared successfully."
	@bash _scripts/test.sh

# upgrade_deps
.PHONY: upgrade-deps
upgrade-deps:
	@go get -d -u -t ./...


