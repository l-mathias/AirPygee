PROJECT=airpygee
VERSION=$$(git rev-parse --short=10 HEAD)

clean:
	go clean -cache

build:
	go build -v main.go

run:
	go run main.go
