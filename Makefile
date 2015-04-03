BUILD_DIR=bin
VERSION=0.1.0

all:
	go test -v ./...

install:
	go install -v ./...

deps:
	go get github.com/gobuild/gobuild3/packer

dist:
	packer --os linux  --arch amd64 --output procker-linux-amd64-$(VERSION).zip
	packer --os linux  --arch 386   --output procker-linux-386-$(VERSION).zip
	packer --os darwin --arch amd64 --output procker-mac-amd64-$(VERSION).zip
	packer --os darwin --arch 386   --output procker-mac-386-$(VERSION).zip

clean:
	rm -rf $(BUILD_DIR) *.zip
