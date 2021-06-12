FROM golang:1.14.7-alpine AS GO_BUILD
COPY . /server
WORKDIR /server
RUN go build -o /go/bin/server


FROM alpine:3.10
RUN apk update && apk upgrade && \
    apk add --no-cache git
WORKDIR /app
RUN git clone -b development https://prampey7@bitbucket.org/cryptopatron/front-end.git
COPY --from=GO_BUILD /go/bin/server ./
RUN ls
CMD ./server --servePath /app/front-end/dist/dropcoin/
