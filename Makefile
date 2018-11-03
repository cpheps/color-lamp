BINARY=colorLamp

.PHONY: build clean fmt run test
default: build

build: | clean
	go build -o ./bin/${BINARY}

clean:
	if [ -f /bin/${BINARY} ] ; then rm bin/${BINARY} ; fi

fmt: 
	@echo Formatting
	@goimports -w .
	@gofmt -s -w .

run:
	go run main.go
test:
#Add other package tests here
	go test -v ./...
	