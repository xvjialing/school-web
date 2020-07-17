FROM golang:1.12.9-alpine
ENV GO111MODULE on
ENV GOPROXY https://mirrors.aliyun.com/goproxy/
ENV APPNAME=go-admin
WORKDIR $GOPATH/src
RUN mkdir -p $APPNAME
WORKDIR $GOPATH/src/$APPNAME
ADD go.mod .
ADD go.sum .
RUN go mod download
ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./build/${APPNAME} -v ./
