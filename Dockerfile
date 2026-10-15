FROM golang:1.24-alpine AS builder
WORKDIR /src
COPY go.mod ./
COPY main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/wildbugs .

FROM alpine:3.22
RUN addgroup -S app && adduser -S -G app app
COPY --from=builder /out/wildbugs /usr/local/bin/wildbugs
USER app
EXPOSE 8080
CMD ["/usr/local/bin/wildbugs"]
