# Stage 1: Build Gitea MCP using Go 1.24
FROM golang:1.24-bookworm AS builder

WORKDIR /app

# Install git and ca-certificates
RUN apt-get update && apt-get install -y git ca-certificates && rm -rf /var/lib/apt/lists/*

# Clone the repo
RUN git clone https://github.com/Muralidhar-br/gitea-mcp.git .

# Build the binary
RUN go build -o gitea-mcp

# Stage 2: Minimal runtime image
FROM debian:bookworm-slim

WORKDIR /app

# Add CA certificates for HTTPS support
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Copy the built binary
COPY --from=builder /app/gitea-mcp /usr/local/bin/gitea-mcp

EXPOSE 4000

ENTRYPOINT ["gitea-mcp"]
