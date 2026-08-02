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

WORKDIR /workspace/gateway
RUN go build -o /bin/gateway ./cmd

FROM alpine:latest
COPY --from=builder /bin/gateway /bin/gateway
EXPOSE 3000

ENTRYPOINT [ "/bin/gateway" ]
