package usecase

import (
	"context"
	"strconv"

	"github.com/ypb/your-personal-blog/internal/blog/model"
	"github.com/ypb/your-personal-blog/internal/blog/repo"
	"github.com/ypb/your-personal-blog/internal/blog/usecase/input"
	"github.com/ypb/your-personal-blog/internal/blog/usecase/output"
)

type postUsecase struct {
	postRepo repo.PostRepo
}

func NewPostUsecase(r repo.PostRepo) Post {
	return &postUsecase{postRepo: r}
}

func (u *postUsecase) Create(ctx context.Context, in *input.PostCreate) (int64, error) {
	post := &model.Post{
		Title:      in.Title,
		Slug:       in.Slug,
		Content:    in.Content,
		Summary:    in.Summary,
		IsTop:      in.IsTop,
		CategoryID: in.CategoryID,
		Status:     model.PostDraft,
	}
	err := u.postRepo.Create(ctx, post)
	return post.ID, err
}

func (u *postUsecase) Update(ctx context.Context, in *input.PostUpdate) error {
	post, err := u.postRepo.GetByID(ctx, in.ID)
	if err != nil {
		return err
	}

	if in.Title != nil {
		post.Title = *in.Title
	}
	if in.Slug != nil {
		post.Slug = *in.Slug
	}
	if in.Content != nil {
		post.Content = *in.Content
	}
	if in.Summary != nil {
		post.Summary = *in.Summary
	}
	if in.IsTop != nil {
		post.IsTop = *in.IsTop
	}
	if in.Status != nil {
		post.Status = *in.Status
	}
	if in.CategoryID != nil {
		post.CategoryID = *in.CategoryID
	}

	return u.postRepo.Update(ctx, post)
}

func (u *postUsecase) Delete(ctx context.Context, id int64) error {
	return u.postRepo.Delete(ctx, id)
}

func (u *postUsecase) Publish(ctx context.Context, id int64) error {
	post, err := u.postRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	post.Status = model.PostPublished
	return u.postRepo.Update(ctx, post)
}

func (u *postUsecase) Unpublish(ctx context.Context, id int64) error {
	post, err := u.postRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	post.Status = model.PostDraft
	return u.postRepo.Update(ctx, post)
}

func (u *postUsecase) Get(ctx context.Context, idOrSlug string) (*output.PostDetail, error) {
	var post *model.Post
	var err error
	if id, err := strconv.ParseInt(idOrSlug, 10, 64); err == nil {
		post, err = u.postRepo.GetByID(ctx, id)
	} else {
		post, err = u.postRepo.GetBySlug(ctx, idOrSlug)
	}

	if err != nil {
		return nil, err
	}

	return &output.PostDetail{
		ID:          post.ID,
		Title:       post.Title,
		Slug:        post.Slug,
		Content:     post.Content,
		HTMLContent: post.HTMLContent,
		Summary:     post.Summary,
		Status:      post.Status,
		IsTop:       post.IsTop,
		ViewCount:   post.ViewCount,
		CategoryID:  post.CategoryID,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
		PublishedAt: post.PublishedAt,
	}, nil
}

func (u *postUsecase) List(ctx context.Context, in *input.PostList) ([]*model.Post, int64, error) {
	offset := (in.Page - 1) * in.PageSize
	return u.postRepo.List(ctx, offset, in.PageSize, in.Status)
}
