.PHONY: lint
lint:
	golangci-lint run

.PHONY: fmt
fmt:
	go fmt ./...
	gofumpt -l -w .

test:
	go test ./... -v

migrate:
	migrate -source file://db/migrations -database 'mysql://root:root@tcp(127.0.0.1:3306)/go' up