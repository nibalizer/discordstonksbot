# Borrowed from https://medium.com/@treeder/multi-stage-docker-builds-for-creating-tiny-go-images-e0e1867efe5a
# build stage
FROM golang:1.21-alpine AS build-env
RUN apk --no-cache add build-base git mercurial gcc clang clang-dev
RUN mkdir /app
ADD go.mod go.sum *.go /app/
RUN cd /app && CXX=clang++ CGO_ENABLED=1 go build -o goapp

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /app/goapp /app/
ENTRYPOINT ./goapp
