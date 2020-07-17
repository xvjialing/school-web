FROM golang:1.12.9-alpine AS development
ENV GO111MODULE on
ENV GOPROXY https://mirrors.aliyun.com/goproxy/
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
RUN mkdir pack && tar -xzvf $APPNAME.tar.gz -C pack && cd pack && rm -rf go.mod go.sum swagger && ls

FROM alpine AS production
WORKDIR /app
ENV APPNAME school-web
COPY --from=development /go/src/$APPNAME/pack .
CMD [ "sh", "-c", "./$APPNAME" ]
