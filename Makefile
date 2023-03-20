BINARY_NAME=deck-of-cards

server: build
	bin/${BINARY_NAME}

build:
	go build -o bin/${BINARY_NAME} main.go

clean:
	go clean
	rm bin/${BINARY_NAME}

test:
	go test ./...

test-json:
	go test -json ./...

lint:
	golangci-lint run
