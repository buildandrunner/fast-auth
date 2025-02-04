FROM golang:1.23.4 AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/main.go

FROM gcr.io/distroless/base-debian12

WORKDIR /

COPY --from=builder /app/main /main
COPY --from=builder /app/templates /templates

EXPOSE 8080

ENTRYPOINT [ "/main" ]
