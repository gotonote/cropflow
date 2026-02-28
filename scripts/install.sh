#!/bin/bash

# CorpFlow 一键安装脚本 (非 Docker 模式)
# 支持 Ubuntu/Debian

set -e

echo "=== CorpFlow 一键安装 ==="

# 1. 安装系统依赖
echo "[1/6] 安装系统依赖..."
sudo apt update
sudo apt install -y curl git wget

# 2. 安装 Go
if ! command -v go &> /dev/null; then
    echo "[2/6] 安装 Go..."
    wget -q https://go.dev/dl/go1.21.6.linux-amd64.tar.gz -O /tmp/go.tar.gz
    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf /tmp/go.tar.gz
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    export PATH=$PATH:/usr/local/go/bin
else
    echo "[2/6] Go 已安装: $(go version)"
fi

# 3. 安装 Node.js
if ! command -v node &> /dev/null; then
    echo "[3/6] 安装 Node.js..."
    curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
    sudo apt install -y nodejs
else
    echo "[3/6] Node.js 已安装: $(node --version)"
fi

# 4. 安装 PostgreSQL
if ! command -v psql &> /dev/null; then
    echo "[4/6] 安装 PostgreSQL..."
    sudo apt install -y postgresql postgresql-contrib
    sudo systemctl start postgresql
    sudo systemctl enable postgresql
fi

# 创建数据库和用户
sudo -u postgres psql -c "CREATE USER corpflow WITH PASSWORD 'corpflow123';" 2>/dev/null || true
sudo -u postgres psql -c "CREATE DATABASE corpflow OWNER corpflow;" 2>/dev/null || true
sudo -u postgres psql -c "ALTER USER corpflow CREATEDB;" 2>/dev/null || true

# 5. 安装 Redis
if ! command -v redis-server &> /dev/null; then
    echo "[5/6] 安装 Redis..."
    sudo apt install -y redis-server
    sudo systemctl start redis-server
    sudo systemctl enable redis-server
fi

# 6. 克隆并配置项目
echo "[6/6] 配置项目..."
cd ~

if [ -d "corpflow" ]; then
    echo "corpflow 目录已存在，更新中..."
    cd corpflow
    git pull
else
    git clone https://github.com/gotonote/corpflow.git
    cd corpflow
fi

# 复制配置
cp .env.example .env

# 安装前端依赖
cd frontend
npm install
cd ..

echo ""
echo "=== 安装完成 ==="
echo ""
echo "配置步骤:"
echo "1. 编辑 .env 文件填入 API Key:"
echo "   vim ~/.corpflow/.env"
echo ""
echo "2. 启动后端 (新终端):"
echo "   cd ~/corpflow && go run cmd/server/main.go"
echo ""
echo "3. 启动前端 (新终端):"
echo "   cd ~/corpflow/frontend && npm run dev"
echo ""
echo "4. 访问 http://localhost:3000"
