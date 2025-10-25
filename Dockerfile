FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN go install github.com/boyter/scc/v3@v3.3.0

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/server ./main.go

FROM gcr.io/distroless/base-debian12 AS final

WORKDIR /app

COPY --from=builder /go/bin/scc /usr/local/bin/

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]

