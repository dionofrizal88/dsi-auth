.PHONY: serve

serve: ## Start grpc server
	@go run main.go

test-coverage: ## Start grpc server
	@go test -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out
