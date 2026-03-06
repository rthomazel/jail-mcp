# ── Stage 1: build ────────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
RUN CGO_ENABLED=0 go build \
    -ldflags="-s -w -X main.version=$(date +%Y%m%d)" \
    -o jail-mcp .

# ── Stage 2: runtime ──────────────────────────────────────────────────────────
# Ubuntu 24.04 for a generous, real toolchain.
# Alpine's busybox is too minimal for real dev work.
FROM ubuntu:24.04

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

# Log directory — bind-mounted at runtime so logs survive container restarts.
RUN mkdir -p /var/log/jail-mcp

COPY --from=builder /build/jail-mcp /usr/local/bin/jail-mcp

# MCP runs over stdio — no port needed.
CMD ["jail-mcp"]
