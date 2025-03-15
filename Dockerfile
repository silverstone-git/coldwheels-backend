FROM golang:1.24-bookworm
WORKDIR /app

COPY ./.env /app/.env
COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum
RUN go mod download

COPY ./lib /app/lib
COPY ./middleware /app/middleware
COPY ./repository /app/repository
COPY ./main.go /app/main.go
RUN GOOS=linux GOARCH=amd64 go build -o ./gowheels

ENV GIN_MODE=release

EXPOSE 4054

CMD ["./gowheels"]

