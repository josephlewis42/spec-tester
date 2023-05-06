.PHONY: build
build: generate vendor bin
	go build -o bin/spec-tester main.go

.PHONY: generate
generate:
	go generate ./...


.PHONY:  vendor
vendor: generate
	go mod vendor

bin:
	mkdir bin

.PHONY: clean
clean:
	rm -rf bin