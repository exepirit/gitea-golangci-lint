# [Build environment image]
FROM golang:alpine AS builder

WORKDIR $GOPATH/src/git.0x73.ru/exepirit/gitea-golangci-lint/
COPY . .
ENV GO111MODULE=on

RUN apk add --no-cache build-base
RUN go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.30.0
RUN go get ./...
RUN go build -v -o /go/bin/gitea-golangci-lint

# [Final image]
FROM golang:alpine
COPY --from=builder /go/bin/golangci-lint /go/bin/golangci-lint
COPY --from=builder /go/bin/gitea-golangci-lint /go/bin/gitea-integration
RUN apk add --no-cache build-base
