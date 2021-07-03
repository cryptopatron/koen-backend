FROM golang:1.14.7-alpine AS GO_BUILD
RUN apk add build-base
# Copy and download dependencies first, to cache them
WORKDIR /server
COPY ./go.* /server/
RUN go mod download
RUN ls
# Copy and run code. This is a fast step anyhow.
COPY . .
RUN go test -v ./...
RUN go build -o /go/bin/server


FROM node:12.11 AS REACT_BUILD
ADD https://api.github.com/repos/cryptopatron/web-app/git/refs/heads/master version.json
RUN git clone https://github.com/cryptopatron/web-app.git /webapp
RUN cp /webapp/package.json /tmp/
RUN cd /tmp && npm install
RUN cp -a /tmp/node_modules /webapp/

WORKDIR /webapp
RUN ls
RUN npm run build


FROM alpine:3.10
WORKDIR /app
COPY --from=REACT_BUILD /webapp/build ./webapp/build
COPY --from=GO_BUILD /go/bin/server ./
CMD ./server --servePath ./webapp/build
