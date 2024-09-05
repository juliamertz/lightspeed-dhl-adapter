clean:
	go mod tidy
	gofmt -w */*.go

build:
	$(MAKE) clean
	go build .
