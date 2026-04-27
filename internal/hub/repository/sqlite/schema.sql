-- Storage Hub 数据库 Schema 设计
-- 采用 SQLite 3 语法

-- 1. 客户端元数据表 (租户管理)
CREATE TABLE IF NOT EXISTS clients (
    id TEXT PRIMARY KEY,               -- 客户端唯一标识 (UUID)
    name TEXT NOT NULL,                 -- 客户端名称
    api_key TEXT UNIQUE NOT NULL,       -- API 访问密钥
    quota_limit INTEGER DEFAULT 0,      -- 存储配额限制 (字节, 0表示无限制)
    status INTEGER DEFAULT 1,           -- 状态: 0:Active, 1:Inactive, 2:Disabled
    created_at INTEGER NOT NULL,        -- 创建时间戳 (Unix int64)
    updated_at INTEGER NOT NULL         -- 更新时间戳 (Unix int64)
);

-- 2. 文件元数据表 (文件索引)
CREATE TABLE IF NOT EXISTS files (
    id TEXT PRIMARY KEY,               -- 文件唯一标识 (UUID)
    client_id TEXT NOT NULL,            -- 所属客户端 ID (关联 clients.id)
    bucket_id TEXT NOT NULL,            -- 逻辑桶 ID
    original_name TEXT NOT NULL,        -- 原始文件名
    mime_type TEXT NOT NULL DEFAULT '', -- 媒体类型 (匹配 Go string)
    file_size INTEGER NOT NULL,         -- 文件大小 (字节)
    storage_path TEXT NOT NULL,         -- 物理存储路径 (StorageEngine 使用的 key)
    status INTEGER DEFAULT 0,           -- 状态: 0:FileUploading, 1:FileActive, 2:FileDeleted
    created_at INTEGER NOT NULL,        -- 创建时间戳 (Unix int64)
    updated_at INTEGER NOT NULL,        -- 更新时间戳 (Unix int64)
    
    FOREIGN KEY (client_id) REFERENCES clients(id) ON DELETE CASCADE
);

-- 索引：加速按客户端和状态查询
CREATE INDEX IF NOT EXISTS idx_files_client_status ON files(client_id, status);
CREATE INDEX IF NOT EXISTS idx_clients_api_key ON clients(api_key);
