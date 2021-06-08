# [Build environment image]
FROM golang:1.16-alpine AS builder

WORKDIR $GOPATH/src/github.com/exepirit/gitea-golangci-lint/
COPY . .
RUN go build -o $GOPATH/bin/gitea-golangci-lint ./...

# [Final image]
FROM alpine:3
COPY --from=builder /go/bin/gitea-golangci-lint /bin/
CMD ["/bin/gitea-golangci-lint"]
