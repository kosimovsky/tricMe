FROM  golang:1.19.3-alpine3.15

WORKDIR /server

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN go build -o /metrics ./cmd/server/main.go

EXPOSE 8080

CMD ["/metrics"]
