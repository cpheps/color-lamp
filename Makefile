BINARY=color_lamp

.PHONY: build clean fmt run test package
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

package:
	@rm -rf package
	mkdir package
	cp -r deliverables/* package/
	cp scripts/install.sh package/
	mkdir package/scripts
	cp scripts/wps_script.sh package/scripts/
	cp -r config package/

	