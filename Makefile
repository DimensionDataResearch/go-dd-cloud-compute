default: test

fmt:
	go fmt ./compute/...

test: fmt
	go test -v ./compute/...
