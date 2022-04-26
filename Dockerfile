FROM golang:1.18.1-alpine3.15 as build

RUN apk upgrade --no-cache --force
RUN apk add --update build-base make git

WORKDIR /go/src/github.com/vibin18/shinobi_ai_bot

# Compile
COPY ./ /go/src/github.com/vibin18/shinobi_ai_bot
RUN make dependencies
RUN make test
RUN make build
RUN ./shinobi-ai-bot --help

FROM scratch AS export-stage
COPY --from=build /go/src/github.com/vibin18/shinobi_ai_bot/shinobi-ai-bot /