FROM golang:1.16-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /bin/app

FROM gcr.io/distroless/static
COPY --from=builder /bin/app /bin/app
ENTRYPOINT ["/bin/app"]
