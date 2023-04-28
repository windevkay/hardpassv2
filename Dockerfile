# syntax=docker/dockerfile:1

FROM golang:1.20.3

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /hardpass/web ./cmd/web

CMD ["/hardpass/web"]
