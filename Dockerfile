FROM golang:alpine

WORKDIR /black

COPY go.mod go.sum /black/

RUN go mod download

COPY . .

EXPOSE 5000

RUN go build

ENTRYPOINT ["./black"]


