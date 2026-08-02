FROM golang:1.26-alpine AS builder

WORKDIR /workdir

COPY go.work go.work.sum ./
COPY ./api/ ./api
COPY ./auth/ ./auth
COPY ./gateway/ ./gateway
COPY ./orders/ ./orders
COPY ./payments/ ./payments
COPY ./shared/ ./shared
COPY ./stock/ ./stock

WORKDIR /workspace/auth
RUN go build -o ./bin/auth ./cmd

FROM alping:latest
COPY --from=builder /bin/auth /bin/auth
EXPOSE 5555

ENTRYPOINT [ "/bin/auth" ]
