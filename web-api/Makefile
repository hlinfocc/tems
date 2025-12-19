export PATH := $(GOPATH)/bin:$(PATH)
export GO111MODULE=on
LDFLAGS := -s -w

all: fmt build

build: tems-web-api


fmt:
	go fmt ./...

fmt-more:
	gofumpt -l -w .

vet:
	go vet ./...

tems-web-api:
	env CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o bin/tems-web-api ./

test: gotest

gotest:
	go test -v --cover ./

	
clean:
	rm -f ./bin/tems-web-api
