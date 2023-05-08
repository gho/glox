run: generate
	go run *.go test

build: generate
	go build

generate:
	go generate

clean:
	rm -f glox *_string.go
