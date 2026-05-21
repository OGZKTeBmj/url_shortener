#Build stage
FROM golang:1.26.1-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/app/main.go

#Runtime stage
FROM alpine:3.20 as runner

RUN apk add --no-cache postgresql-client bash

WORKDIR /app

COPY --from=builder /build/app .
COPY ./wait-for-postgres.sh .

RUN chmod +x ./wait-for-postgres.sh

RUN adduser -D appuser
USER appuser

ENTRYPOINT ["./wait-for-postgres.sh"]
CMD ["./app"]
