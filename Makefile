build:
	go build ./cmd/goash

release:
	go build -ldflags="-s"  ./cmd/goash
