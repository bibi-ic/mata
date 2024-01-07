FROM golang:1.21 as builder
ENV DEPLOY=PRODUCT
WORKDIR /runtime

COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/app ./cmd

FROM gcr.io/distroless/static-debian11

WORKDIR /runtime
COPY --from=builder --chown=nonroot:nonroot /go/bin/app .
COPY --chown=nonroot:nonroot config/deploy.yaml ./config/deploy.yaml
COPY --chown=nonroot:nonroot internal/db/migration ./internal/db/migration

EXPOSE 8080

USER nonroot:nonroot
CMD ["./app"]
