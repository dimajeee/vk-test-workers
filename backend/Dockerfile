# syntax=docker/dockerfile:experimental

FROM --platform=$BUILDPLATFORM golang:1.24-bookworm as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    GOOS=linux GOARCH=$(echo $TARGETPLATFORM | cut -d'/' -f2) CGO_ENABLED=0 go build -o server ./cmd/server

FROM --platform=$BUILDPLATFORM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/server .

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

EXPOSE 8080

CMD ["./server"]
