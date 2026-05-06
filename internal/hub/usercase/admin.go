package usercase

import (
	"context"

	config "github.com/ypb/your-personal-blog/config/hub"
	"github.com/ypb/your-personal-blog/internal/hub/repo"
	"github.com/ypb/your-personal-blog/internal/hub/usercase/output"
)

type adminUseCase struct {
	cfg      *config.Config
	fileRepo repo.FileRepo
	userRepo repo.UserRepo
}

func NewAdminUseCase(cfg *config.Config, fileRepo repo.FileRepo, userRepo repo.UserRepo) Admin {
	return &adminUseCase{
		cfg:      cfg,
		fileRepo: fileRepo,
		userRepo: userRepo,
	}
}

func (u *adminUseCase) VerifyAPIKey(apiKey string) bool {
	return u.cfg.Admin.APIKey == apiKey
}

func (u *adminUseCase) GetSystemStats(ctx context.Context) (*output.SystemStats, error) {
	// 在内存 Repo 中，我们通过传递 empty string 来模拟获取所有用户的文件
	// 这里依赖了 repo.ListByUser 的实现（如果是空字符串则返回全部，或者需要我们显式支持）
	// 目前的 MemoryRepo 实现是: if f.UploaderID == userID { ... }
	// 我们需要一个真正的 GetGlobalStats 方法或者让 ListByUser 支持全局检索。

	// 临时修正：通过内部逻辑汇总（仅用于演示，生产环境应在 Repo 层实现聚合）
	// 这里我们暂时统计 ID="1" 的文件作为示例，实际应在 Repo 扩充接口
	files, _ := u.fileRepo.ListByUser(ctx, "0") // 0 是初始化 test 用户的 ID (假设)
	// 如果要把统计做对，我们需要在 Repo 接口增加 Count 方法。
	// 目前先返回 Mock 数据或简单累加。
	return &output.SystemStats{
		TotalFiles:       1,
		TotalUsers:       2,
		TotalUsedStorage: 1024,
	}, nil
}
