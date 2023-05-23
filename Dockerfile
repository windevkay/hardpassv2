FROM golang:1.20.3 as build

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o web ./cmd/web

FROM scratch as image

COPY --from=build /app/web .

EXPOSE 4000

CMD ["/web"]
