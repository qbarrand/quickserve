all: qs

qs: $(shell find . -type f -name '*.go') go.mod go.sum VERSION
	go build -ldflags="-X main.version=$(cat VERSION)" -o $@ ./cmd/qs

clean:
	rm qs
