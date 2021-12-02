FROM golang:alpine AS builder

COPY . /application

WORKDIR /application/app
RUN go build -o /api


FROM alpine

COPY ./data /data
COPY --from=builder /api /api

CMD ["sh", "-c", "/api --port=${PORT} --gamesCode=/data/gamesCode.jsonc"]