build:
	go mod tidy
	gofmt -w */*.go
	CGO_ENABLED=0 go build .
