FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY *.go ./
RUN apk --no-cache add --repository=http://dl-cdn.alpinelinux.org/alpine/edge/community github-cli=2.32.0-r1
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /bin/app

FROM gcr.io/distroless/static
COPY --from=builder /bin/app /bin/app
COPY --from=builder /usr/bin/gh /bin/gh
ENTRYPOINT ["/bin/app"]
