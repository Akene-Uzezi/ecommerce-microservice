# syntax=docker/dockerfile:1
FROM golang:1.26-alpine AS builder

WORKDIR /workspace

COPY go.work go.work.sum ./
COPY ./api/ ./api
COPY ./auth/ ./auth
COPY ./gateway/ ./gateway
COPY ./orders/ ./orders
COPY ./payments/ ./payments
COPY ./shared/ ./shared
COPY ./stock/ ./stock

WORKDIR /workspace/orders
RUN --mount=type=cache,target=/root/go/pkg/mod \
  --mount=type=cache,target=/root/.cache/go-build \
  go build -o /bin/orders ./cmd

FROM alpine:latest
COPY --from=builder /bin/orders /bin/orders
EXPOSE 4444
ENTRYPOINT [ "/bin/orders" ]
