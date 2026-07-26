# ---------- 构建阶段 ----------
FROM golang:1.22-bookworm AS builder

WORKDIR /build

# 先拷贝依赖清单，利用构建缓存
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o gochat .

# ---------- 运行阶段 ----------
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /build/gochat .
# 运行期需要的文件：页面模板、静态资源、GeoIP 库、建表 SQL、配置模板
COPY --from=builder /build/static ./static
COPY --from=builder /build/config/GeoLite2-City.mmdb ./config/GeoLite2-City.mmdb
COPY --from=builder /build/config/mysql.json.demo ./config/mysql.json.demo
COPY --from=builder /build/import.sql .

EXPOSE 8081

CMD ["./gochat", "server", "-p", "8081"]
