FROM golang:1.20-alpine AS builder

WORKDIR /builder

COPY . /builder

RUN apk add make

RUN make build

FROM alpine:3.18.4

WORKDIR /app

RUN apk add ffmpeg

COPY --from=builder /builder/build/komputer .

COPY assets /app/assets

CMD [ "/app/komputer" ]