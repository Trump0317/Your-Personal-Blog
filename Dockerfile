# ── 阶段 1: 前端构建 ──
FROM node:22-alpine AS frontend-builder
WORKDIR /src
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# ── 阶段 2: 后端构建 ──
FROM golang:1.24-alpine AS backend-builder
RUN apk add --no-cache gcc musl-dev sqlite-dev
WORKDIR /src
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=1 go build -ldflags="-s -w" -o /blog-server ./cmd/blog-server/

# ── 阶段 3: 运行时 ──
FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai
ENV BLOG_PORT=:8080

WORKDIR /app
COPY --from=backend-builder /blog-server .
COPY --from=frontend-builder /src/dist ./frontend/dist
COPY config.json .

EXPOSE 8080
VOLUME ["/app/uploads"]

HEALTHCHECK --interval=30s --timeout=3s \
  CMD wget -qO- http://localhost:8080/api/stats || exit 1

CMD ["./blog-server"]
