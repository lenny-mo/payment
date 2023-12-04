
GOPATH:=$(shell go env GOPATH)
MODIFY=proto/

.PHONY: proto
proto:
    
	protoc --micro_out=${MODIFY} --go_out=${MODIFY} proto/payment.proto
    

.PHONY: build
build: proto

	go build -o payment-service *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t payment-service:latest
