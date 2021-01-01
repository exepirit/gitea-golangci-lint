# [Build environment image]
FROM golang:alpine AS builder

WORKDIR $GOPATH
RUN apk add --no-cache build-base wget
RUN wget -O - -q https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s v1.30.0

WORKDIR $GOPATH/src/github.com/exepirit/gitea-golangci-lint/
COPY . .
ENV GO111MODULE=on
RUN go build -o $GOPATH/bin/gitea-golangci-lint ./...

# [Final image]
FROM golang:alpine
COPY --from=builder /go/bin/golangci-lint /go/bin/golangci-lint
COPY --from=builder /go/bin/gitea-golangci-lint /go/bin/gitea-integration
