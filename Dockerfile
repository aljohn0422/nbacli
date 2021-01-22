FROM golang:alpine

RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates
RUN go get github.com/manifoldco/promptui
WORKDIR /app
COPY . /app
RUN go build -o nba .
CMD ./nba

FROM alpine
WORKDIR /app
COPY --from=0 /app /app/
CMD ./nba