FROM --platform=$BUILDPLATFORM golang:1.26.5 AS builder

WORKDIR /app

ARG TARGETOS
ARG TARGETARCH
ARG VERSION=rc

COPY . .

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags="-X 'main.version=${VERSION}'" -o bin/ranga-${VERSION} ./cmd
RUN chmod +x bin/ranga-${VERSION}

FROM scratch AS exporter
COPY --from=builder /app/bin/ /
