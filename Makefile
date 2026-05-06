.PHONY: build run clean build-cli help

# 变量定义
BINARY_DIR=bin
BLOG_SERVER_OUT=$(BINARY_DIR)/blog-server
BLOG_CLI_OUT=$(BINARY_DIR)/blog-cli

# 默认构建目标
build: build-server build-cli

# 构建博客服务器
build-server:
	@echo "正在编译 blog-server..."
	@mkdir -p $(BINARY_DIR)
	go build -o $(BLOG_SERVER_OUT) cmd/blog-server/main.go

# 构建命令行工具
build-cli:
	@echo "正在编译 blog-cli..."
	@mkdir -p $(BINARY_DIR)
	go build -o $(BLOG_CLI_OUT) cmd/blog-cli/main.go

# 运行服务器
run: build-server
	@echo "正在启动服务器..."
	./$(BLOG_SERVER_OUT)

# 清理构建产物
clean:
	@echo "清理二进制文件..."
	rm -rf $(BINARY_DIR)

# 帮助信息
help:
	@echo "Makefile 使用说明:"
	@echo "  make build          - 编译所有二进制文件"
	@echo "  make run            - 编译并运行博客服务器"
	@echo "  make build-server   - 仅编译博客服务器"
	@echo "  make build-cli      - 仅编译命令行工具"
	@echo "  make clean          - 删除 bin 目录"
