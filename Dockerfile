FROM golang:1.22.2-alpine3.19 AS builder

WORKDIR /builder

COPY . /builder

RUN apk add make gcc libc-dev pkgconfig opus-dev

RUN make prod

FROM alpine:3.19.1

WORKDIR /app

RUN apk add ffmpeg opus

COPY --from=builder /builder/build/komputer .

VOLUME [ "/assets" ]

ENV ASSETS_DIR="/assets"

CMD [ "/app/komputer" ]