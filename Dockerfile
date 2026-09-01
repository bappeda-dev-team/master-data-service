ARG GO_VERSION=1.25.5

# =========================================================
# BUILD STAGE
# =========================================================

FROM registry.docker.com/library/golang:$GO_VERSION-alpine AS builder

WORKDIR /app

# Install packages needed to build
RUN apk add --no-cache git make

# Download dependency lebih cache-friendly
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN make clean
RUN make build-docker


# =========================================================
# FINAL STAGE
# =========================================================

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/bin/app /app/app

COPY --from=builder /app/docs/ ./docs

ENTRYPOINT ["/app/app"]
