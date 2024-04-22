# 使用 Go 官方镜像作为编译阶段的基础镜像
FROM golang:1.22 AS builder

WORKDIR /app

COPY . .

RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod tidy

RUN go build -o jbactive cmd/main.go

RUN chmod +x jbactive

# 清理编译时生成的临时文件和缓存
RUN go clean -cache -testcache

# 删除源代码（保留编译后的可执行文件）
RUN rm -rf cmd internal

# 使用轻量级基础镜像作为运行阶段
FROM alpine:latest AS runtime

WORKDIR /app

# 将编译阶段生成的可执行文件以及静态文件复制到运行时镜像
COPY --from=builder /app/jbactive .
COPY --from=builder /app/static static

RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
#暴露接口
EXPOSE 10800

# 启动应用
CMD ["./jbactive", "-addr", ":10800"]
