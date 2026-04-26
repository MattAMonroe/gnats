BIN=bin/gnats
.PHONY=build test run verify clean tidy

build: tidy test
	go build -o ${BIN} ./...

run:
	go run main.go

test: verify
	go test -v -count=1 ./...

clean:
	rm -f ${BIN}

tidy:
	go fmt ./...
	go mod tidy -v

verify:
	go mod verify
	go vet ./...
	GOFLAGS="-buildvcs=false" go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
