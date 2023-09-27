FROM golang:1.19

WORKDIR /dockerapp

COPY . .

RUN go mod tidy

RUN go mod download
RUN go mod verify

RUN go build -o /main ./cmd/main.go


CMD ["/main"]