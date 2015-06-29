BUILD_DIR=bin
VERSION=0.1.1

all: tests

install:
	go install -v ./...

tests: deps
	go test -v ./...

deps:
	go get ./...
	go get github.com/gobuild/gobuild3/packer

qa:
	go vet
	golint
	go test -coverprofile=.cover~
	go tool cover -html=.cover~

dist:
	packer --os linux  --arch amd64 --output procker-linux-amd64-$(VERSION).zip
	packer --os linux  --arch 386   --output procker-linux-386-$(VERSION).zip
	packer --os darwin --arch amd64 --output procker-mac-amd64-$(VERSION).zip
	packer --os darwin --arch 386   --output procker-mac-386-$(VERSION).zip

clean:
	rm -rf $(BUILD_DIR) *.zip
