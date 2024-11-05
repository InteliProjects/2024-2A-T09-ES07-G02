FROM golang:1.23

WORKDIR /speech

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 7070

CMD ["go", "run", "main.go"]

