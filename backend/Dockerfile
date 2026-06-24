# Stage 1: compile the Go binary
# We use the official Go image (alpine variant = smaller base) as our build environment
FROM golang:1.26-alpine AS builder

# Set the working directory inside the build container
WORKDIR /app

# Copy dependency files first and download modules
# Doing this before copying the rest of the code means Docker caches this layer —
# if only source files change (not go.mod/go.sum), the download step is skipped on rebuild
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code into the build container
COPY . .

# Compile the binary
# CGO_ENABLED=0 produces a fully static binary (no C library dependencies)
# GOOS=linux ensures the binary targets Linux even if you build on a Mac
# -o server names the output binary "server"
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/main.go


# Stage 2: the actual image that runs in production
# We start fresh from a minimal alpine image — no Go toolchain, no source code, just the binary
FROM alpine:3.20

WORKDIR /app

# Copy only the compiled binary from the build stage
# The migrations and seed data are embedded inside the binary at compile time, so nothing else is needed
COPY --from=builder /app/server .

# Document that the container listens on port 8080 (Cloud Run also reads this)
EXPOSE 8080

# Run the binary when the container starts
CMD ["./server"]
