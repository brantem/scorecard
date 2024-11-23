FROM --platform=linux/amd64 golang:alpine AS builder
WORKDIR /scorecard/api

RUN apk add --no-cache make build-base
COPY . .
RUN CGO_ENABLED=1 make build

FROM alpine:latest

COPY --from=builder /scorecard/api/server .
CMD ["./server"]
