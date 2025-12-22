FROM golang:1.25 AS development
WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.* ./
RUN go mod download

COPY . .

CMD ["sleep", "infinity"]

FROM golang:1.25 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/server .
CMD ["./server"]