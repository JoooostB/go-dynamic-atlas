FROM golang:1.23-alpine AS builder

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app/
COPY . .

RUN apk add --no-cache ca-certificates
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-w -s" -o go-dynamic-atlas *.go

FROM scratch
WORKDIR /app/
COPY --from=builder /app/go-dynamic-atlas /app/go-dynamic-atlas
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENTRYPOINT [ "/app/go-dynamic-atlas" ]
