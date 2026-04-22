#!/bin/bash
# YPB 开发环境一键验证脚本

echo "=== YPB 环境检查 ==="

# 检查 Go
if command -v go &> /dev/null; then
    echo "✅ Go: $(go version)"
else
    echo "❌ Go: 未安装"
fi

# 检查 Node.js
if command -v node &> /dev/null; then
    echo "✅ Node.js: $(node --version)"
else
    echo "❌ Node.js: 未安装"
fi

# 检查 SQLite3
if command -v sqlite3 &> /dev/null; then
    echo "✅ SQLite3: $(sqlite3 --version | head -1)"
else
    echo "❌ SQLite3: 未安装"
fi

# 检查 Docker
if command -v docker &> /dev/null; then
    echo "✅ Docker: $(docker --version)"
else
    echo "❌ Docker: 未安装"
fi

# 检查 gofastdfs
if curl -s http://localhost:8080/status &> /dev/null; then
    echo "✅ gofastdfs: 运行中"
else
    echo "⚠️  gofastdfs: 未运行或未安装"
fi

echo "=== 检查完成 ==="