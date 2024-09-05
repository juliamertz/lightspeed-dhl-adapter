clean:
	go mod tidy
	gofmt -w */*.go

check:
	$(MAKE) clean
	go test ./...

build:
	$(MAKE) clean
	go build .
