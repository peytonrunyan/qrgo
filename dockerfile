# syntax=docker/dockerfile:1

FROM golang:1.16-alpine as builder
ENV GO111MODULE=on

WORKDIR /go/src/app

COPY . .
RUN go build ./cmd/qrcode-backend

FROM alpine:3.11.3
COPY --from=builder /go/src/app/qrcode-backend .
COPY .env .
EXPOSE 8082

CMD [ "./qrcode-backend" ]