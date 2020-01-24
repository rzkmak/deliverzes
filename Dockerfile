FROM golang:1.13 AS builder
WORKDIR /go/src/aeidelos/deliverzes
COPY . .
RUN GO111MODULE=on go mod download
RUN GO111MODULE=on go mod verify
RUN GO111MODULE=on go build .

FROM alpine:latest AS deployment
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/aeidelos/deliverzes/deliverzes .
COPY --from=builder /go/src/aeidelos/deliverzes/web .
CMD ["./deliverzes"]
