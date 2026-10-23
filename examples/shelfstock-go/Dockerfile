# syntax=docker/dockerfile:1
# Multi-stage build. Targets:
#   test     -> go vet + race tests (used by `task docker:check`; fails the build if checks fail)
#   runtime  -> small, non-root production image (default final stage)

ARG GO_VERSION=1.24

FROM golang:${GO_VERSION} AS base
WORKDIR /src
COPY go.mod ./
# go.sum does not exist yet (stdlib only). Add `COPY go.sum ./` when dependencies are introduced.
RUN go mod download
COPY . .

FROM base AS test
RUN go vet ./... && go test -race ./...

FROM base AS build
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/shelfstock ./cmd/shelfstock

FROM alpine:3.20 AS runtime
RUN addgroup -S app && adduser -S -G app app
COPY --from=build /out/shelfstock /usr/local/bin/shelfstock
USER app
ENV SHELFSTOCK_ADDR=:8080
EXPOSE 8080
HEALTHCHECK --interval=15s --timeout=3s --start-period=5s --retries=3 \
  CMD wget -q -O /dev/null http://127.0.0.1:8080/healthz || exit 1
ENTRYPOINT ["/usr/local/bin/shelfstock"]
