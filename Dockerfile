FROM golang:1.14.7-alpine AS GO_BUILD
# Copy and download dependencies first, to cache them
WORKDIR /server
# dependecneis for testing
RUN apk add build-base
COPY ./go.* /server/
RUN go mod download
# Copy and run code. This is a fast step anyhow.
COPY pkg pkg
COPY server.go .
RUN go test -v /server/pkg/...
RUN go build -o /go/bin/server


FROM node:14.15.3 AS REACT_BUILD
ARG GITHUB_CRED
ARG BRANCH
RUN curl -u ${GITHUB_CRED} https://api.github.com/repos/cryptopatron/web-app/git/refs/heads/${BRANCH} > /tmp/version.json
RUN cat /tmp/version.json
RUN git clone -b ${BRANCH} https://${GITHUB_CRED}@github.com/cryptopatron/web-app.git /webapp
RUN cp /webapp/package.json /tmp/
RUN cd /tmp && npm install && npm install yarn && yarn
RUN cp -a /tmp/node_modules /webapp/

WORKDIR /webapp
RUN npm run build


FROM alpine:3.10
WORKDIR /app
ARG VERSION
ENV VERSION=${VERSION}
RUN echo $VERSION
COPY --from=REACT_BUILD /webapp/build ./webapp/build
COPY --from=GO_BUILD /go/bin/server ./
CMD ./server --servePath ./webapp/build
