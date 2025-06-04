
FROM golang:1.24.2 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .


RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o main .


FROM alpine:latest

WORKDIR /root/

COPY --from=build /app/main .

EXPOSE 8080

CMD ["./main"]
