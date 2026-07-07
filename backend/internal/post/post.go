package post

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// ── 模型 ──

type Post struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Content     string     `json:"content"`
	HTMLContent string     `json:"html_content"`
	Summary     string     `json:"summary"`
	CategoryID  string     `json:"category_id"`
	TagIDs      []string   `json:"tag_ids"`
	Published   bool       `json:"published"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	// 展示时填充
	Category *CategoryRef `json:"category,omitempty"`
	Tags     []TagRef     `json:"tags,omitempty"`
}

type CategoryRef struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type TagRef struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// ── 存储接口 ──

type Store interface {
	Create(ctx context.Context, p *Post) error
	Get(ctx context.Context, idOrSlug string) (*Post, error)
	List(ctx context.Context, f Filter) ([]*Post, int, error)
	Update(ctx context.Context, p *Post) error
	Delete(ctx context.Context, id string) error
	SetTags(ctx context.Context, postID string, tagIDs []string) error
}

type Filter struct {
	CategoryID string
	TagID      string
	Query      string
	Published  *bool
	Offset     int
	Limit      int
}

// ── 依赖接口（由 category/tag 包实现）──

type CategoryStore interface {
	Get(ctx context.Context, id string) (*CategoryRef, error)
	List(ctx context.Context) ([]*CategoryRef, error)
}

type TagStore interface {
	GetByIDs(ctx context.Context, ids []string) ([]TagRef, error)
	List(ctx context.Context) ([]TagRef, error)
}

// ── 业务逻辑 ──

type Service struct {
	store    Store
	catStore CategoryStore
	tagStore TagStore
}

func NewService(store Store, catStore CategoryStore, tagStore TagStore) *Service {
	return &Service{store: store, catStore: catStore, tagStore: tagStore}
}

func (s *Service) Create(ctx context.Context, in CreateInput) (string, error) {
	if strings.TrimSpace(in.Title) == "" {
		return "", errors.New("title required")
	}

	now := time.Now()
	slug := in.Slug
	if slug == "" {
		slug = toSlug(in.Title)
		slug = uniqueSlug(slug, now)
	}

	p := &Post{
		ID:         uuid.New().String(),
		Title:      strings.TrimSpace(in.Title),
		Slug:       slug,
		Content:    in.Content,
		Summary:    in.Summary,
		CategoryID: in.CategoryID,
		TagIDs:     in.TagIDs,
		Published:  in.Published,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if p.Content != "" {
		p.HTMLContent = renderMarkdown(p.Content)
	}
	if p.Summary == "" && p.Content != "" {
		p.Summary = extractSummary(p.Content)
	}

	if err := s.store.Create(ctx, p); err != nil {
		return "", err
	}
	if len(p.TagIDs) > 0 {
		_ = s.store.SetTags(ctx, p.ID, p.TagIDs)
	}
	return p.ID, nil
}

// CreateInput 创建文章输入
type CreateInput struct {
	Title      string
	Slug       string
	Content    string
	Summary    string
	CategoryID string
	TagIDs     []string
	Published  bool
}

func (s *Service) Get(ctx context.Context, idOrSlug string) (*Post, error) {
	p, err := s.store.Get(ctx, idOrSlug)
	if err != nil {
		return nil, err
	}
	return s.fillRefs(ctx, p)
}

func (s *Service) List(ctx context.Context, f Filter) ([]*Post, int, error) {
	posts, total, err := s.store.List(ctx, f)
	if err != nil {
		return nil, 0, err
	}
	result := make([]*Post, len(posts))
	for i, p := range posts {
		filled, _ := s.fillRefs(ctx, p)
		result[i] = filled
	}
	return result, total, nil
}

func (s *Service) Update(ctx context.Context, in UpdateInput) error {
	p, err := s.store.Get(ctx, in.ID)
	if err != nil {
		return err
	}

	if in.Title != nil {
		p.Title = *in.Title
	}
	if in.Content != nil {
		p.Content = *in.Content
		p.HTMLContent = renderMarkdown(p.Content)
		if in.Summary == nil {
			p.Summary = extractSummary(p.Content)
		}
	}
	if in.Summary != nil {
		p.Summary = *in.Summary
	}
	if in.CategoryID != nil {
		p.CategoryID = *in.CategoryID
	}
	if in.TagIDs != nil {
		p.TagIDs = in.TagIDs
	}
	if in.Published != nil {
		p.Published = *in.Published
	}
	if in.Slug != nil {
		p.Slug = *in.Slug
	}
	p.UpdatedAt = time.Now()

	if err := s.store.Update(ctx, p); err != nil {
		return err
	}
	return s.store.SetTags(ctx, p.ID, p.TagIDs)
}

// UpdateInput 更新文章输入
type UpdateInput struct {
	ID         string
	Title      *string
	Slug       *string
	Content    *string
	Summary    *string
	CategoryID *string
	TagIDs     []string
	Published  *bool
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.store.Delete(ctx, id)
}

func (s *Service) Publish(ctx context.Context, id string) error {
	p, err := s.store.Get(ctx, id)
	if err != nil {
		return err
	}
	p.Published = true
	p.UpdatedAt = time.Now()
	return s.store.Update(ctx, p)
}

func (s *Service) Unpublish(ctx context.Context, id string) error {
	p, err := s.store.Get(ctx, id)
	if err != nil {
		return err
	}
	p.Published = false
	p.UpdatedAt = time.Now()
	return s.store.Update(ctx, p)
}

// fillRefs 填充关联的分类和标签
func (s *Service) fillRefs(ctx context.Context, p *Post) (*Post, error) {
	if p.CategoryID != "" && s.catStore != nil {
		cat, err := s.catStore.Get(ctx, p.CategoryID)
		if err == nil {
			p.Category = cat
		}
	}
	if len(p.TagIDs) > 0 && s.tagStore != nil {
		tags, err := s.tagStore.GetByIDs(ctx, p.TagIDs)
		if err == nil {
			p.Tags = tags
		}
	}
	if p.Tags == nil {
		p.Tags = []TagRef{}
	}
	return p, nil
}

// ── 辅助 ──

func toSlug(s string) string {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, " ", "-")
	return s
}

func uniqueSlug(base string, t time.Time) string {
	return fmt.Sprintf("%s-%d", base, t.UnixMilli()%100000)
}

func extractSummary(content string) string {
	// 去除 markdown 标记，取前 150 个字符
	plain := content
	plain = strings.ReplaceAll(plain, "#", "")
	plain = strings.ReplaceAll(plain, "*", "")
	plain = strings.ReplaceAll(plain, "`", "")
	plain = strings.TrimSpace(plain)
	if len(plain) > 150 {
		return plain[:150] + "..."
	}
	return plain
}

func renderMarkdown(md string) string {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(md))

	opts := html.RendererOptions{Flags: html.CommonFlags}
	renderer := html.NewRenderer(opts)

	return string(markdown.Render(doc, renderer))
}
