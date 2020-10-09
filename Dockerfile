# front
# FROM node:13.8 as front
# LABEL maintainer "toshick"

# ADD . /go/src/github.com/try-csrf/
# WORKDIR /go/src/github.com/try-csrf/front-ronten

# RUN npm ci && npm test && npm run generate

# golang
FROM golang:1.14
LABEL maintainer "toshick"

ENV PATH $PATH:/usr/local/go/bin
ENV GO111MODULE on
ENV GO_ENV prd

ADD . /go/src/github.com/toshick/try-csrf/
WORKDIR /go/src/github.com/toshick/try-csrf/


# COPY --from=front /go/src/github.com/try-csrf/front-ronten/dist /go/src/github.com/try-csrf/public

CMD ["/usr/local/go/bin/go", "run", "/go/src/github.com/toshick/try-csrf/app/main.go"]
# CMD ["/bin/bash"]

EXPOSE 80