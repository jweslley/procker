PROGRAM=procker
VERSION=0.1.1
LDFLAGS="-X main.programVersion=$(VERSION)"

all: test

deps:
	go get ./...

install: deps
	go install -a -v -ldflags $(LDFLAGS) ./cmd/...

test: deps
	go test -v ./...

qa:
	go vet
	golint
	go test -coverprofile=.cover~
	go tool cover -html=.cover~

dist:
	@for os in linux darwin; do \
		for arch in 386 amd64; do \
			target=$(PROGRAM)-$$os-$$arch-$(VERSION); \
			echo Building $$target; \
			GOOS=$$os GOARCH=$$arch go build -ldflags $(LDFLAGS) -o $$target/$(PROGRAM) ./cmd/... ; \
			cp ./README.md ./LICENSE $$target; \
			tar -zcf $$target.tar.gz $$target; \
			rm -rf $$target;                   \
		done                                 \
	done

clean:
	rm -rf *.tar.gz
