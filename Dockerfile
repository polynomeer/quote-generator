FROM golang:1.22-alpine AS build
WORKDIR /app
COPY . .
RUN apk add --no-cache build-base && \
    go build -o /quote-generator ./cmd/quote-generator

FROM alpine:3.20
WORKDIR /app
COPY --from=build /quote-generator /app/quote-generator
COPY configs/config.yaml /app/config.yaml
EXPOSE 8080
ENV QG_CONFIG=/app/config.yaml
ENTRYPOINT ["/app/quote-generator"]