FROM golang:1.26.2-alpine AS dev

WORKDIR /app

RUN apk add --no-cache git

RUN go install github.com/air-verse/air@v1.63.0 \
  && go install github.com/pressly/goose/v3/cmd/goose@v3.24.1 \
  && go install github.com/aarondl/sqlboiler/v4@v4.19.7 \
  && go install github.com/aarondl/sqlboiler/v4/drivers/sqlboiler-psql@v4.19.7

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]