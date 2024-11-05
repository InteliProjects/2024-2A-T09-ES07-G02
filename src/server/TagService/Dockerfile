FROM golang:1.22

WORKDIR /tag

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 6969

CMD ["go", "run", "main.go"]

