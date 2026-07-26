#!/usr/bin/env bash
#
# GoFly LiveChat 一键部署脚本（Linux + docker compose）
#
# 用法:
#   ./deploy.sh                 # 拉取/更新源码并构建启动
#   REPO_URL=xxx BRANCH=v1 ./deploy.sh   # 覆盖默认仓库地址和分支
#
set -euo pipefail

REPO_URL="${REPO_URL:-git@github.com:Lemonopf/goflylivechat.git}"
BRANCH="${BRANCH:-v1}"
DIR="${DIR:-goflylivechat}"

echo "==> 仓库: $REPO_URL  分支: $BRANCH  目录: $DIR"

# 1. 拉取或更新源码
if [ -d "$DIR/.git" ]; then
    echo "==> 更新已有源码..."
    git -C "$DIR" fetch origin
    git -C "$DIR" checkout "$BRANCH"
    git -C "$DIR" pull origin "$BRANCH"
else
    echo "==> 克隆源码..."
    git clone -b "$BRANCH" "$REPO_URL" "$DIR"
fi

cd "$DIR"

# 2. 准备配置：首次部署从模板生成 mysql.json，并提示先编辑再部署
if [ ! -f config/mysql.json ]; then
    cp config/mysql.json.demo config/mysql.json
    echo
    echo "==================================================="
    echo " 已生成 config/mysql.json ，请先编辑填入你的 MySQL 信息:"
    echo "   vi $(pwd)/config/mysql.json"
    echo " 然后重新执行本脚本: ./deploy.sh"
    echo "==================================================="
    exit 1
fi

# server.json 必须提前 touch，否则 docker 会把不存在的挂载文件创建成目录
touch config/server.json
chmod 600 config/server.json config/mysql.json

mkdir -p data/upload logs

# 3. 构建并启动
echo "==> 构建镜像并启动容器..."
docker compose up -d --build

echo
echo "==> 部署完成，访问: http://<服务器IP>:8081"
echo
echo " 首次部署需要初始化数据库(建表+导入默认数据)，执行:"
echo "   docker compose exec gofly ./gochat install"
echo " 默认客服账号: agent / agent   (登录后请立即修改密码)"
echo
echo " 常用命令:"
echo "   查看日志:  docker compose logs -f gofly"
echo "   重启:      docker compose restart gofly"
echo "   停止:      docker compose down"
