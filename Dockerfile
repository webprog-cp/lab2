FROM golang:1.26.2-alpine3.23@sha256:f85330846cde1e57ca9ec309382da3b8e6ae3ab943d2739500e08c86393a21b1 AS build

WORKDIR /src
ENV CGO_ENABLED=0 GOTOOLCHAIN=local

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY docs/ ./docs/
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -mod=readonly -trimpath -ldflags="-s -w" -o /out/api ./cmd/api

FROM scratch
WORKDIR /app
COPY --from=build /out/api /app/api
COPY --from=build /src/docs/openapi.yaml /app/docs/openapi.yaml
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/app/api"]
