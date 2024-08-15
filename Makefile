build:
	go build -o ./bin/blockchain_scratch

run: build
	./bin/blockchain_scratch

test:
	go test -v ./...
