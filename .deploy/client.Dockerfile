FROM golang:1.22.1 AS builder

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN  go mod download

COPY . .

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-X 'github.com/ivan-bokov/pow-ddos/internal/version.Version=v1.0.0'" -o word-of-wisdom-cli cmd/client/main.go

FROM scratch

COPY --from=builder /build/word-of-wisdom-cli /

ENTRYPOINT ["/word-of-wisdom-cli"]