# ── Stage 1: build ────────────────────────────────────────────────────────────
FROM golang:1.25-alpine AS builder

WORKDIR /build
COPY . .
RUN go mod download

ARG VERSION
RUN CGO_ENABLED=0 go build \
    # trim debug info and set version
    -ldflags="-s -w -X main.version=${VERSION:-dev}" \
    -o jail-mcp .

# ── Stage 2: runtime ──────────────────────────────────────────────────────────
FROM ubuntu:24.04

# skip prompts from apt
ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y \
    # shells
    bash \
    # core utils
    curl wget git make jq \
    # build tools
    gcc g++ build-essential pkg-config \
    # Go (for building inside the container if needed)
    golang \
    # scripting
    python3 python3-pip nodejs npm \
    # file tools
    zip unzip tar \
    # text / search
    ripgrep vim nano less \
    # network
    dnsutils iputils-ping netcat-openbsd \
    && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /var/log/jail-mcp

COPY --from=builder /build/jail-mcp /usr/local/bin/jail-mcp

CMD ["jail-mcp"]
