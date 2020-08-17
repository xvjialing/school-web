FROM registry.cn-shenzhen.aliyuncs.com/xvjialing/school-web:base AS development
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
