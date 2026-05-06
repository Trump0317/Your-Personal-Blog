package model

// Admin 管理员模型 - 信息从配置文件获取，不需要 ID 和时间戳
type Admin struct {
	Username string `json:"username"`
	Password string `json:"-"`
	APIKey   string `json:"api_key"`
}
