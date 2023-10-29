FROM golang:1.20-alpine AS builder

WORKDIR /builder

COPY . /builder

RUN apk add make gcc libc-dev pkgconfig opus-dev

RUN make build

FROM alpine:3.18.4

WORKDIR /app

RUN apk add ffmpeg opus

COPY --from=builder /builder/build/komputer .

COPY assets /app/assets

CMD [ "/app/komputer" ]