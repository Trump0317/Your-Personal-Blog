package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/ypb/your-personal-blog/internal/blog/model"
	"github.com/ypb/your-personal-blog/internal/blog/repo"
)

type postUsecase struct {
	postRepo     repo.PostRepo
	categoryRepo repo.CategoryRepo
	tagUC        Tag // 仅用于展示时补全标签对象
}

// NewPostUsecase 创建文章业务逻辑实例
func NewPostUsecase(pr repo.PostRepo, cr repo.CategoryRepo, tu Tag) Post {
	return &postUsecase{
		postRepo:     pr,
		categoryRepo: cr,
		tagUC:        tu,
	}
}

func (u *postUsecase) Create(ctx context.Context, in *PostCreateInput) (string, error) {
	// 1. 构造文章模型
	post := &model.Post{
		ID:         fmt.Sprintf("post_%d", time.Now().UnixNano()), // 临时生成 ID
		Title:      in.Title,
		Slug:       in.Slug,
		Summary:    in.Summary,
		Content:    in.Content,
		CategoryID: in.CategoryID,
		TagIDs:     in.TagIDs,
		Status:     in.Status,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if in.Status == model.PostPublished {
		now := time.Now()
		post.PublishedAt = &now
	}

	// 2. 保存文章本体
	if err := u.postRepo.Create(ctx, post); err != nil {
		return "", fmt.Errorf("failed to create post: %w", err)
	}

	// 3. 保存标签关联关系
	if err := u.postRepo.SetTags(ctx, post.ID, in.TagIDs); err != nil {
		return "", fmt.Errorf("failed to set post tags: %w", err)
	}

	return post.ID, nil
}

func (u *postUsecase) Update(ctx context.Context, in *PostUpdateInput) error {
	// 1. 获取原文章
	existing, err := u.postRepo.GetByID(ctx, in.ID)
	if err != nil {
		return fmt.Errorf("failed to get existing post: %w", err)
	}

	// 2. 处理标签更新
	if in.TagIDs != nil {
		if err := u.postRepo.SetTags(ctx, in.ID, *in.TagIDs); err != nil {
			return err
		}
		existing.TagIDs = *in.TagIDs
	}

	// 3. 更新基础字段
	if in.Title != nil {
		existing.Title = *in.Title
	}
	if in.Slug != nil {
		existing.Slug = *in.Slug
	}
	if in.Content != nil {
		existing.Content = *in.Content
	}
	if in.Summary != nil {
		existing.Summary = *in.Summary
	}
	if in.CategoryID != nil {
		existing.CategoryID = *in.CategoryID
	}
	if in.Status != nil {
		existing.Status = *in.Status
	}
	existing.UpdatedAt = time.Now()

	return u.postRepo.Update(ctx, existing)
}

func (u *postUsecase) Delete(ctx context.Context, id string) error {
	return u.postRepo.Delete(ctx, id)
}

func (u *postUsecase) Publish(ctx context.Context, id string) error {
	return u.postRepo.UpdateStatus(ctx, id, model.PostPublished)
}

func (u *postUsecase) Unpublish(ctx context.Context, id string) error {
	return u.postRepo.UpdateStatus(ctx, id, model.PostDraft)
}

func (u *postUsecase) Get(ctx context.Context, idOrSlug string) (*PostDetailOutput, error) {
	// 1. 获取文章本体 (尝试按 Slug, 失败则按 ID)
	post, err := u.postRepo.GetBySlug(ctx, idOrSlug)
	if err != nil {
		post, err = u.postRepo.GetByID(ctx, idOrSlug)
		if err != nil {
			return nil, err
		}
	}

	// 2. 补全分类信息
	var category *model.Category
	if post.CategoryID != "" {
		category, _ = u.categoryRepo.GetByID(ctx, post.CategoryID)
	}

	// 3. 补全标签信息
	var tags []*model.Tag
	if len(post.TagIDs) > 0 {
		tags, _ = u.tagUC.ListByIDs(ctx, post.TagIDs)
	}

	return &PostDetailOutput{
		ID:    post.ID,
		Title: post.Title, Slug: post.Slug,
		Content:     post.Content,
		HTMLContent: post.HTMLContent,
		Summary:     post.Summary,
		Status:      post.Status,
		ViewCount:   post.ViewCount,
		CategoryID:  post.CategoryID,
		Category:    category,
		Tags:        tagsToValueSlice(tags),
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
		PublishedAt: post.PublishedAt,
	}, nil
}

func (u *postUsecase) List(ctx context.Context, in *PostListInput) (*PostListOutput, error) {
	page := in.Page
	if page < 1 {
		page = 1
	}
	pageSize := in.PageSize
	if pageSize < 1 {
		pageSize = 10
	}

	filter := &repo.PostFilter{
		Offset: (page - 1) * pageSize,
		Limit:  pageSize,
	}

	// 修复：必须显式转换 model.PostStatus 指针
	if in.Status != nil {
		s := model.PostStatus(*in.Status)
		filter.Status = &s
	}

	if in.CategoryID != nil {
		filter.CategoryID = *in.CategoryID
	}
	if in.TagID != nil {
		// 优先尝试将名称解析为 ID
		tags, _ := u.tagUC.List(ctx)
		foundID := ""
		for _, t := range tags {
			if t.Name == *in.TagID || t.ID == *in.TagID {
				foundID = t.ID
				break
			}
		}
		if foundID != "" {
			filter.TagID = foundID
		} else {
			filter.TagID = *in.TagID
		}
	}
	if in.Query != "" {
		filter.Query = in.Query
	}

	posts, total, err := u.postRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	// 打印 Repos 返回的结果
	fmt.Printf("DEBUG: Repo.List returned %d items, total %d\n", len(posts), total)

	items := make([]*PostItem, 0, len(posts))
	for _, p := range posts {
		// 调试日志：fmt.Printf("DEBUG: Found post %s status %v\n", p.ID, p.Status)
		var category *model.Category
		if p.CategoryID != "" {
			category, _ = u.categoryRepo.GetByID(ctx, p.CategoryID)
		}

		var tags []*model.Tag
		if len(p.TagIDs) > 0 {
			tags, _ = u.tagUC.ListByIDs(ctx, p.TagIDs)
		}

		items = append(items, &PostItem{
			ID:          p.ID,
			Title:       p.Title,
			Slug:        p.Slug,
			Summary:     p.Summary,
			Status:      p.Status,
			ViewCount:   p.ViewCount,
			Category:    category,
			Tags:        tagsToValueSlice(tags),
			CreatedAt:   p.CreatedAt,
			PublishedAt: p.PublishedAt,
		})
	}

	return &PostListOutput{
		Posts: items,
		Total: total,
	}, nil
}

// 辅助函数：将对象切片转为模型中的非指针切片
func tagsToValueSlice(tags []*model.Tag) []model.Tag {
	res := make([]model.Tag, 0, len(tags))
	for _, t := range tags {
		if t != nil {
			res = append(res, *t)
		}
	}
	return res
}
