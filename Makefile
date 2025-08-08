BINARY_NAME = cashd
GOBIN = $(HOME)/.local/bin
TARGET_OS = darwin linux
TARGET_ARCH = arm64 amd64

all: build

$(BINARY_NAME):

build:
	go mod tidy
	go build -o $(BINARY_NAME)

run: build
	./$(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME)
	rm -rf release

install: build
	GOBIN=$(GOBIN) go install

test:
	go test

release:
	mkdir -p release
	for os in $(TARGET_OS); do \
		for arch in $(TARGET_ARCH); do \
			echo "Building release/$(BINARY_NAME)-$$os-$$arch"; \
			GOOS=$$os GOARCH=$$arch go build -o release/$(BINARY_NAME)-$$os-$$arch; \
			tar -C release -czf release/$(BINARY_NAME)-$$os-$$arch.tar.gz $(BINARY_NAME)-$$os-$$arch; \
			rm release/$(BINARY_NAME)-$$os-$$arch; \
		done \
	done
	(cd release && shasum -a 256 *.tar.gz > checksum.txt)


.PHONY: all build run clean install fmt vet test release
