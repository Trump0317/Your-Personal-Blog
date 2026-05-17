package repo

import (
	"context"
	"errors"
	"strings"
	"sync"

	"github.com/ypb/your-personal-blog/internal/blog/model"
)

// --- PostRepo Implementation ---

type memoryPostRepo struct {
	mu    sync.RWMutex
	posts map[string]*model.Post
}

func NewMemoryPostRepo() PostRepo {
	return &memoryPostRepo{
		posts: make(map[string]*model.Post),
	}
}

func (r *memoryPostRepo) Create(ctx context.Context, post *model.Post) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.posts[post.ID]; ok {
		return errors.New("post already exists")
	}
	r.posts[post.ID] = post
	return nil
}

func (r *memoryPostRepo) GetByID(ctx context.Context, id string) (*model.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	post, ok := r.posts[id]
	if !ok {
		return nil, errors.New("post not found")
	}
	return post, nil
}

func (r *memoryPostRepo) GetBySlug(ctx context.Context, slug string) (*model.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.posts {
		if p.Slug == slug {
			return p, nil
		}
	}
	return nil, errors.New("post not found")
}

func (r *memoryPostRepo) Update(ctx context.Context, post *model.Post) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.posts[post.ID]; !ok {
		return errors.New("post not found")
	}
	r.posts[post.ID] = post
	return nil
}

func (r *memoryPostRepo) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.posts, id)
	return nil
}

func (r *memoryPostRepo) UpdateStatus(ctx context.Context, id string, status model.PostStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	post, ok := r.posts[id]
	if !ok {
		return errors.New("post not found")
	}
	post.Status = status
	if status == model.PostPublished && post.PublishedAt == nil {
		now := post.UpdatedAt
		post.PublishedAt = &now
	}
	return nil
}

func (r *memoryPostRepo) SetTags(ctx context.Context, postID string, tagIDs []string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	post, ok := r.posts[postID]
	if !ok {
		return errors.New("post not found")
	}
	post.TagIDs = tagIDs
	return nil
}

func (r *memoryPostRepo) List(ctx context.Context, filter *PostFilter) ([]*model.Post, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var allPosts []*model.Post
	for _, p := range r.posts {
		allPosts = append(allPosts, p)
	}

	// 调试日志：打印内存中所有文章的状态
	// for _, p := range allPosts {
	// 	fmt.Printf("Memory Repo List: ID=%s, Status=%v, PublishedAt=%v\n", p.ID, p.Status, p.PublishedAt)
	// }

	var filtered []*model.Post
	for _, p := range allPosts {
		if filter != nil {
			if filter.Status != nil && p.Status != *filter.Status {
				continue
			}
			if filter.CategoryID != "" && p.CategoryID != filter.CategoryID {
				continue
			}
			if filter.TagID != "" {
				found := false
				for _, tid := range p.TagIDs {
					if tid == filter.TagID {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}
			if filter.Query != "" {
				if !strings.Contains(p.Title, filter.Query) && !strings.Contains(p.Content, filter.Query) {
					continue
				}
			}
		}
		filtered = append(filtered, p)
	}

	total := int64(len(filtered))

	// Apply Offset and Limit
	start := 0
	end := len(filtered)
	if filter != nil {
		if filter.Offset > 0 {
			start = filter.Offset
		}
		if filter.Limit > 0 {
			end = start + filter.Limit
		}
	}

	if start > len(filtered) {
		return []*model.Post{}, total, nil
	}
	if end > len(filtered) {
		end = len(filtered)
	}

	return filtered[start:end], total, nil
}

// --- CategoryRepo Implementation ---

type memoryCategoryRepo struct {
	mu         sync.RWMutex
	categories map[string]*model.Category
}

func NewMemoryCategoryRepo() CategoryRepo {
	return &memoryCategoryRepo{
		categories: make(map[string]*model.Category),
	}
}

func (r *memoryCategoryRepo) Create(ctx context.Context, category *model.Category) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.categories[category.ID]; ok {
		return errors.New("category already exists")
	}
	r.categories[category.ID] = category
	return nil
}

func (r *memoryCategoryRepo) GetByID(ctx context.Context, id string) (*model.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cat, ok := r.categories[id]
	if !ok {
		return nil, errors.New("category not found")
	}
	return cat, nil
}

func (r *memoryCategoryRepo) GetBySlug(ctx context.Context, slug string) (*model.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, cat := range r.categories {
		if cat.Slug == slug {
			return cat, nil
		}
	}
	return nil, errors.New("category not found")
}

func (r *memoryCategoryRepo) List(ctx context.Context) ([]*model.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*model.Category
	for _, cat := range r.categories {
		list = append(list, cat)
	}
	return list, nil
}

func (r *memoryCategoryRepo) Update(ctx context.Context, category *model.Category) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.categories[category.ID]; !ok {
		return errors.New("category not found")
	}
	r.categories[category.ID] = category
	return nil
}

func (r *memoryCategoryRepo) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.categories, id)
	return nil
}

// --- TagRepo Implementation ---

type memoryTagRepo struct {
	mu   sync.RWMutex
	tags map[string]*model.Tag
}

func NewMemoryTagRepo() TagRepo {
	return &memoryTagRepo{
		tags: make(map[string]*model.Tag),
	}
}

func (r *memoryTagRepo) Save(ctx context.Context, tag *model.Tag) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tags[tag.ID] = tag
	return nil
}

func (r *memoryTagRepo) GetByID(ctx context.Context, id string) (*model.Tag, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	tag, ok := r.tags[id]
	if !ok {
		return nil, errors.New("tag not found")
	}
	return tag, nil
}

func (r *memoryTagRepo) GetByName(ctx context.Context, name string) (*model.Tag, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, tag := range r.tags {
		if tag.Name == name {
			return tag, nil
		}
	}
	return nil, errors.New("tag not found")
}

func (r *memoryTagRepo) List(ctx context.Context) ([]*model.Tag, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*model.Tag
	for _, tag := range r.tags {
		list = append(list, tag)
	}
	return list, nil
}

func (r *memoryTagRepo) BatchGetByIDs(ctx context.Context, ids []string) ([]*model.Tag, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*model.Tag
	for _, id := range ids {
		if tag, ok := r.tags[id]; ok {
			list = append(list, tag)
		}
	}
	return list, nil
}

func (r *memoryTagRepo) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.tags, id)
	return nil
}
