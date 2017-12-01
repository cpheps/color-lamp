BINARY=colorLamp

# These are the values we want to pass for Version and BuildTime
VERSION=1.0.0
BUILD_TIME=`date +%FT%T%z`

LDFLAGS=-ldflags "-X github.com/cpheps/color-lamp/main.Version=${VERSION} -X github.com/cpheps/color-lamp/main.BuildTime=${BUILD_TIME}"

.PHONY: build clean fmt run vet imports test
default: build

build: | clean
	go build ${LDFLAGS} -o ./bin/${BINARY}

clean:
	if [ -f /bin/${BINARY} ] ; then rm bin/${BINARY} ; fi

fmt: imports
	go fmt github.com/cpheps/color-lamp/...

run:
	go run main.go

vet:
	go tool vet .

imports:
	goimports -w .

test:
#Add other package tests here
	go test -v ./...
	