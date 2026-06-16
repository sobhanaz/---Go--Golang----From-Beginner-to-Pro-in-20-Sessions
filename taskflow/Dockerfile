# ---- Build stage ----
# Compile a fully static binary. Because TaskFlow uses the pure-Go SQLite
# driver (modernc.org/sqlite), we can set CGO_ENABLED=0 — no C toolchain needed.
FROM golang:1.25 AS build

WORKDIR /src

# Download dependencies first (this layer is cached unless go.mod/go.sum change).
COPY go.mod go.sum ./
RUN go mod download

# Copy the source and build.
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /taskflow .

# ---- Runtime stage ----
# A minimal, secure base: no shell, no package manager, just CA certs.
FROM gcr.io/distroless/static-debian12

WORKDIR /app
COPY --from=build /taskflow /app/taskflow

# Default config; override at runtime with -e or a compose file.
ENV TASKFLOW_ADDR=":8080" \
    TASKFLOW_DB="/app/data/taskflow.db"

EXPOSE 8080

# distroless has no shell, so use the exec form.
ENTRYPOINT ["/app/taskflow"]
