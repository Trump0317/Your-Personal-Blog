.PHONY: dev build run clean

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

# ── 清理 ──
clean:
	rm -rf bin/ frontend/dist/ backend/blog.db

help:
	@echo "make dev      开发模式"
	@echo "make build    生产构建"
	@echo "make run      构建并运行"
	@echo "make clean    清理"
