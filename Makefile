.PHONY: vet tidy build

# Go
tidy:
	go mod tidy

vet:
	go vet ./...

build: vet tidy
	go build ./...