FROM golang:1.25.1-alpine

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum
RUN go mod download

COPY ./lib /app/lib
COPY ./middleware /app/middleware
COPY ./repository /app/repository
COPY ./main.go /app/main.go

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o gowheels .

RUN adduser -D appuser

RUN chown appuser:appuser gowheels

USER appuser

ENV GIN_MODE=release

EXPOSE 4054

CMD ["./gowheels"]

