FROM golang:1.26-alpine AS builder

WORKDIR /workspace

COPY go.work go.work.sum ./
COPY api/go.mod api/go.sum ./api/
COPY auth/go.mod ./auth/
COPY gateway/go.mod ./gateway/
COPY orders/go.mod ./orders/
COPY payments/go.mod ./payments/
COPY shared/go.mod shared/go.sum ./shared/
COPY stock/go.mod ./stock/

RUN go mod download

COPY ./api/ ./api
COPY ./auth/ ./auth
COPY ./gateway/ ./gateway
COPY ./orders/ ./orders
COPY ./payments/ ./payments
COPY ./shared/ ./shared
COPY ./stock/ ./stock


WORKDIR /workspace/auth
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/auth ./cmd

FROM alpine:latest
COPY --from=builder /bin/auth /bin/auth
EXPOSE 5555

ENTRYPOINT [ "/bin/auth" ]
