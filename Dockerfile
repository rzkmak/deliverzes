FROM golang:1.13
WORKDIR /go/src/aeidelos/deliverzes
COPY . .
RUN GO111MODULE=on go mod download
RUN GO111MODULE=on go mod verify
RUN GO111MODULE=on GOOS=linux go build .
CMD ["./deliverzes"]
