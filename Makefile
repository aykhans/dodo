lint:
	golangci-lint run

build:
	go build -ldflags "-s -w" -o "./dodo"

build-all:
	rm -rf ./binaries
	./build.sh
