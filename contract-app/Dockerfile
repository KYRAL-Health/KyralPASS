FROM golang:1.16-alpine as builder
RUN apk add build-base
COPY go.mod go.sum /go/src/github.com/kyral/sdk/go/
WORKDIR /go/src/github.com/kyral/sdk/go
RUN go mod download
COPY . /go/src/github.com/kyral/sdk/go
RUN GOOS=linux go build -a -installsuffix cgo -o build/bc-api github.com/kyral/sdk/go

FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/github.com/kyral/sdk/go/build/bc-api /usr/bin/bc-api
EXPOSE 1234 1234
ENTRYPOINT ["/usr/bin/bc-api"]