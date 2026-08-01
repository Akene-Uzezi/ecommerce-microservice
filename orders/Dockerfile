FROM golang:1.26-alpine AS builder

WORKDIR /workspace

COPY go.work go.work.sum ./
COPY ./api/ ./api
COPY ./shared ./shared
COPY ./orders ./orders

WORKDIR /workspace/orders
RUN go build -o /bin/orders ./cmd

FROM alpine:latest
COPY --from=builder /bin/orders /bin/orders
EXPOSE 4444
ENTRYPOINT [ "/bin/orders" ]
