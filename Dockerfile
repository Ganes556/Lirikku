FROM node:20.12.1 as frontend

COPY ./frontend /app/frontend

WORKDIR /app/frontend

RUN npm install && npm run build

#
FROM golang:1.21 as builder

WORKDIR /app

COPY . .

RUN go install github.com/a-h/templ/cmd/templ@latest
RUN go clean --modcache && go mod tidy
RUN templ generate
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:3.14

#
# COPY --from=frontend ./app/static ./app/static
# COPY ./static ./app/static

COPY --from=frontend /app/static ./app/static

RUN apk --no-cache add ca-certificates

RUN apk add --no-cache ffmpeg

ARG env_file_path=./.env.production

ENV $(cat $env_file_path | xargs)

WORKDIR /app

COPY --from=builder /app/app .

ARG PORT=8000
ENV PORT=${PORT}
EXPOSE ${PORT}

CMD [ "./app" ]