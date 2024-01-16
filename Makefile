BINARY_NAME=warc-browser

start-browser:
	chromium --remote-debugging-port=9222

build-front:
	npm --prefix ./web/ install
	npm --prefix ./web/ run build

build: build-front
	go mod tidy
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME} ./cmd/cli.go

run: build
	./${BINARY_NAME}

test:
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

clean:
	go clean
	rm ${BINARY_NAME}
	rm -rf ./web/dist
	rm -rf ./web/node_modules
