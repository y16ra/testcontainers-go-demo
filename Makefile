.PHONY: test
test:
	@echo "Running tests..."
	@cd demo && go test -v ./...
