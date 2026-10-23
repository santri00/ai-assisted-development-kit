# syntax=docker/dockerfile:1
# Go multi-stage template. Targets: test (checks, fails build on failure) and runtime (small, non-root).
# Replace SERVICE and the main package path. Pin base images by digest for production.
ARG GO_VERSION=1.24

FROM golang:${GO_VERSION} AS base
WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download
COPY . .

FROM base AS test
RUN go vet ./... && go test -race ./...

FROM base AS build
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/SERVICE ./cmd/SERVICE

FROM alpine:3.20 AS runtime
RUN addgroup -S app && adduser -S -G app app
COPY --from=build /out/SERVICE /usr/local/bin/SERVICE
USER app
EXPOSE 8080
HEALTHCHECK --interval=15s --timeout=3s --start-period=5s --retries=3 \
  CMD wget -q -O /dev/null http://127.0.0.1:8080/healthz || exit 1
ENTRYPOINT ["/usr/local/bin/SERVICE"]
