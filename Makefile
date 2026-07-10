.PHONY: dev build run clean docker-build docker-up docker-down test

# ── 开发 ──
dev:
	@echo "Starting: backend :8080 | frontend :3000"
	@trap 'kill 0' EXIT; \
		(cd backend && go run ./cmd/blog-server/ &) \
		(cd frontend && npm run dev &) \
		wait

# ── 生产构建 ──
build:
	@echo "Building frontend..."
	cd frontend && npm run build
	@echo "Building backend..."
	cd backend && go build -o ../bin/blog-server ./cmd/blog-server/
	@echo "Done → bin/blog-server"

# ── 生产运行 ──
run: build
	./bin/blog-server

# ── Docker ──
docker-build:
	docker build -t ypb:latest .

docker-up:
	docker compose up -d

docker-down:
	docker compose down

# ── 测试 ──
test:
	cd backend && go test ./internal/...

# ── 清理 ──
clean:
	rm -rf bin/ frontend/dist/ backend/blog.db

help:
	@echo "make dev           开发模式"
	@echo "make build         生产构建"
	@echo "make run           构建并运行"
	@echo "make docker-build  Docker 构建"
	@echo "make docker-up     Docker 启动"
	@echo "make docker-down   Docker 停止"
	@echo "make test          运行测试"
	@echo "make clean         清理"
