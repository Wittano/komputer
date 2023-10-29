FROM golang:1.20-alpine AS builder

WORKDIR /builder

COPY . /builder

RUN apk add make

RUN make build

FROM golang:1.20-alpine

WORKDIR /app

RUN apk add ffmpeg

COPY --from=builder /builder/build/komputer .

COPY assets /app/assets

CMD [ "/app/komputer" ]