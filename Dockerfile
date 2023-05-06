FROM golang:1.20-alpine as builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:3.14

RUN apk --no-cache add ca-certificates

RUN apk add --no-cache ffmpeg

ARG env_file_path=./.env.production

ENV $(cat $env_file_path | xargs)

WORKDIR /usr/src/app

COPY --from=builder /usr/src/app/app .

EXPOSE 8000

CMD [ "./app" ]
