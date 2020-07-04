FROM golang:1.14-alpine AS builder

WORKDIR /app
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -o out/wemo-exporter .

FROM scratch

COPY --from=builder /app/out/wemo-exporter /wemo-exporter

EXPOSE 8080

ENTRYPOINT ["/wemo-exporter"]