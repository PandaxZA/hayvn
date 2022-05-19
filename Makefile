
install:
	go clean && \
	go get -d -v ./... && \
	go install -v ./... && \
	go get && \
	go mod tidy

build:
	go build .

start:
	go run .

test:
	go test -v -cover ./... -short

ci:
	go get -u github.com/swaggo/swag/cmd/swag@v1.7.8 && \
	swag init