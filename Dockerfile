FROM golang:1.23.4-alpine AS builder

WORKDIR /app
RUN apk add --no-cache curl
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
COPY template ./template
ENV GH_VERSION="2.64.0"
RUN curl -s -L "https://github.com/cli/cli/releases/download/v${GH_VERSION}/gh_${GH_VERSION}_linux_amd64.tar.gz" -o "gh_${GH_VERSION}_linux_amd64.tar.gz"
RUN tar -xvf "gh_${GH_VERSION}_linux_amd64.tar.gz"
RUN cp "gh_${GH_VERSION}_linux_amd64/bin/gh" /usr/bin/gh
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /bin/app

FROM alpine:3.21.0
COPY --from=builder /bin/app /bin/app
COPY --from=builder /usr/bin/gh /bin/gh
ENTRYPOINT ["/bin/app"]
