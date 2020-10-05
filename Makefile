.PHONY: proto
proto: protoc-gen-go-json
	PATH=$(CURDIR):$$PATH \
       protoc \
       --go_out=example \
       --go-json_out=logtostderr=true,v=10,discard_unknown:example \
       example/*.proto

protoc-gen-go-json: main.go
	go build .

.PHONY: clean
clean:
	rm protoc-gen-go-json

test:
	go test -v ./...
