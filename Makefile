install:
	@go get ./...
	@go install github.com/vektra/mockery/v2@v2.20.0

test-integration:
	@echo "Running integration tests"
	@go test -v --tags=integration ./...

test:
	@echo "Running unit tests"
	@go test -v ./...

doc:
	@swag init --pd -d cmd/api/ -o cmd/api/docs