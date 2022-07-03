.PHONY: lint fmt test rtest migrate swag
lint:
	golangci-lint run

fmt:
	go fmt ./...
	gofumpt -l -w .

test:
	go test ./... -v

rtest:
	richgo test ./... -v

migrate:
	migrate -source file://db/migrations -database 'mysql://root:root@tcp(127.0.0.1:3306)/go' up

swag:
	swag init