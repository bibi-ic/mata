FROM golang:1.21 as builder
ENV DEPLOY=PRODUCT
WORKDIR /runtime
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/app ./cmd

FROM gcr.io/distroless/static-debian11

WORKDIR /runtime
COPY --from=builder --chown=nonroot:nonroot /go/bin/app .
COPY --chown=nonroot:nonroot config/deploy.yaml ./config/deploy.yaml
COPY --chown=nonroot:nonroot internal/db/migration ./internal/db/migration

EXPOSE 8080

USER nonroot:nonroot
CMD ["./app"]
