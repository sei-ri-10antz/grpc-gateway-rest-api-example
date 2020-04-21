# Usage:
# make proto-rpc
# make build-server
# make docker-server

GOPATH := $(shell go env GOPATH)

TARGET = $(word 1,$(subst -, ,$*))

.PONEY: proto proto-%
.PONEY: build build-%
.PONEY: docker docker-%

proto proto-%:
	@echo "Build gRPC protobuf ${TARGET}..."; \

	protoc -I. \
	  -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	  -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway \
	  --proto_path=.:${GOPATH}/src \
	  --grpc-gateway_out=paths=source_relative:. \
	  --swagger_out=logtostderr=true:. \
	  --go_out=plugins=grpc,paths=source_relative:. proto/${TARGET}/*.proto; \

	@echo ✓ compiled: ${TARGET} protobuf; \



build build-%:
	@echo "Build binary for ${TARGET}..."; \

	go build -o ./build/${TARGET} cmd/${TARGET}/main.go
#	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w' -i -o ${TARGET} cmd/${TARGET}/main.go

	@echo ✓ compiled:; \

docker docker-%:
	@echo "Build image for ${TARGET}..."; \

	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w' -i -o ${TARGET} cmd/${TARGET}/main.go \
	docker build . -t $(TAGS) \
	rm -f ${TARGET} \

	@echo ✓ compiled; \

