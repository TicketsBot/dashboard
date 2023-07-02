# Build
FROM golang:1.19-alpine

ARG branch

RUN apk update && apk upgrade && apk add git zlib-dev gcc musl-dev

RUN mkdir -p /tmp/compile
WORKDIR /tmp/compile

RUN git clone --recurse-submodules https://github.com/TicketsBot/GoPanel .
RUN cd locale && git pull origin master
RUN git checkout $branch
RUN go build -o panel cmd/panel/main.go

# Prod container
FROM alpine:latest

RUN apk update && apk upgrade && apk add curl

COPY --from=0 /tmp/compile/locale /srv/panel/locale
COPY --from=0 /tmp/compile/panel /srv/panel/panel
RUN chmod +x /srv/panel/panel

COPY --from=0 /tmp/compile/emojis.json /srv/panel/emojis.json

RUN adduser container --disabled-password --no-create-home
USER container
WORKDIR /srv/panel

CMD ["/srv/panel/panel"]