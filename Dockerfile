# [Build environment image]
FROM golang:1.16-alpine AS builder

WORKDIR $GOPATH/src/github.com/exepirit/gitea-golangci-lint/

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN mkdir /build
RUN go build -o /build/gitea-golangci-lint

# [Final image]
FROM alpine:3
COPY --from=builder /build/gitea-golangci-lint /bin/
CMD ["/bin/gitea-golangci-lint"]
