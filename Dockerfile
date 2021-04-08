FROM golang:alpine AS build

RUN apk add git
RUN mkdir -p /go/src/github.com/nibalizer/discordbot
WORKDIR /go/src/github.com/nibalizer/discordbot
COPY main.go go.* /go/src/github.com/nibalizer/discordbot
RUN echo $GOPATH
RUN go get
RUN CGO_ENABLED=0 go build -o /bin/discordstonkbot

FROM alpine
RUN apk add bash lftp
COPY --from=build /bin/discordstonkbot /bin/discordstonkbot
COPY --from=build /go/src/github.com/nibalizer/stonksapi/contrib/get_stonks_db.sh /
COPY --from=build /go/src/github.com/nibalizer/stonksapi/contrib/ftpcmds.txt /
COPY stonksdata.txt /stonksdata.txt
ENTRYPOINT ["/bin/discordstonkbot"]
