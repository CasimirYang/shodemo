## 引入最新的golan ，不设置版本即为最新版本
FROM golang:1.14

RUN mkdir -p /www/webapp
## 设置工作目录
WORKDIR /www/webapp

COPY src /www/webapp

RUN go mod download
RUN go build -o gdp .
EXPOSE 50051
# 执行可执行文件
CMD ["./gdp"]