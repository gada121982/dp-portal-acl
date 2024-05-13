FROM golang:1.22 

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN CGO_ENABLED=1 go build -o /app/bin/server ./cmd/*.go

EXPOSE 3003

CMD ["/app/bin/server"]