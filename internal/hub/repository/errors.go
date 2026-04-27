package repository

import "errors"

// 定义 Repository 层通用的哨兵错误
// Service 层可以不依赖底层数据库驱动（如 sql.ErrNoRows）来处理业务逻辑

var (
	// ErrNotFound 表示请求的记录或文件不存在
	ErrNotFound = errors.New("repository: entity not found")

	// ErrAlreadyExists 表示尝试创建已存在的记录
	ErrAlreadyExists = errors.New("repository: entity already exists")

	// ErrInternal 表示底层存储发生的不可预知异常
	ErrInternal = errors.New("repository: internal storage error")
)
