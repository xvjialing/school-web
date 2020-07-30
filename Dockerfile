FROM golang:1.12.9-alpine AS development
ENV GO111MODULE on
ENV GOPROXY https://goproxy.cn
RUN go get github.com/beego/bee
ENV APPNAME=school-web
WORKDIR $GOPATH/src
RUN mkdir -p $APPNAME
WORKDIR $GOPATH/src/$APPNAME
ADD go.mod .
ADD go.sum .
RUN go mod download
ADD . .
RUN bee pack -be GOOS=linux
RUN mkdir pack && tar -xzvf $APPNAME.tar.gz -C pack && cd pack && rm -rf go.mod go.sum && ls

FROM alpine AS production
RUN apk add -U tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata \
    && date
WORKDIR /app
ENV APPNAME school-web
COPY --from=development /go/src/$APPNAME/pack .
CMD [ "sh", "-c", "date && ./$APPNAME" ]
