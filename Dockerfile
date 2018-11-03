FROM golang:1.11 as builder

WORKDIR /src

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix "static" -o app main.go

FROM scratch

# USER 1000

COPY --from=builder /src/app /app

ENTRYPOINT ["/app"]
