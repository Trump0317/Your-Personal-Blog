package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	baseURL    = envOr("YPB_URL", "http://localhost:8080")
	configPath = envOr("YPB_CONFIG", homeDir()+"/.ypb_config")
	token      = ""
	client     = &http.Client{Timeout: 10 * time.Second}
)

func main() {
	// 尝试加载已保存的 token
	loadToken()

	if len(os.Args) < 2 {
		usage()
		return
	}

	switch os.Args[1] {
	case "login":
		cmdLogin()
	case "new":
		cmdNew()
	case "push":
		cmdPush()
	case "pull":
		cmdPull()
	case "list":
		cmdList()
	case "publish":
		cmdPublish()
	case "delete":
		cmdDelete()
	default:
		usage()
	}
}

func usage() {
	fmt.Println(`YPB CLI — Your Personal Blog 命令行工具

用法:
  ypb login                  登录并保存 token
  ypb new <title>            创建一篇新文章（草稿）
  ypb push <file.md>         推送本地 Markdown 到博客
  ypb pull                   拉取所有文章到当前目录
  ypb list                   列出所有文章
  ypb publish <id>           发布文章
  ypb delete <id>            删除文章

环境变量:
  YPB_URL     博客地址 (默认 http://localhost:8080)
  YPB_TOKEN   Token (优先级高于配置文件)`)
}

func requireToken() string {
	if token != "" {
		return token
	}
	t := os.Getenv("YPB_TOKEN")
	if t != "" {
		return t
	}
	fmt.Fprintln(os.Stderr, "未登录，请先运行: ypb login")
	os.Exit(1)
	return ""
}

// ── login ──
func cmdLogin() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("用户名: ")
	user, _ := reader.ReadString('\n')
	user = strings.TrimSpace(user)

	fmt.Print("密码: ")
	pass, _ := reader.ReadString('\n')
	pass = strings.TrimSpace(pass)

	body, _ := json.Marshal(map[string]string{"username": user, "password": pass})
	resp, err := client.Post(baseURL+"/api/admin/login", "application/json", bytes.NewReader(body))
	if err != nil {
		fmt.Fprintf(os.Stderr, "登录失败: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Fprintln(os.Stderr, "登录失败: 用户名或密码错误")
		os.Exit(1)
	}

	var result struct{ Token string }
	json.NewDecoder(resp.Body).Decode(&result)

	if result.Token != "" {
		os.WriteFile(configPath, []byte(result.Token), 0600)
		fmt.Println("登录成功，token 已保存到", configPath)
	} else {
		fmt.Fprintln(os.Stderr, "登录失败: 无效响应")
		os.Exit(1)
	}
}

// ── new ──
func cmdNew() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "用法: ypb new <标题>")
		os.Exit(1)
	}
	title := strings.Join(os.Args[2:], " ")

	body, _ := json.Marshal(map[string]interface{}{
		"title":     title,
		"content":   "# " + title + "\n\n",
		"published": false,
	})

	resp := apiRequest("POST", "/api/admin/posts", body)
	var result struct{ ID string }
	json.NewDecoder(resp.Body).Decode(&result)
	resp.Body.Close()

	if resp.StatusCode == 201 {
		fmt.Println("草稿已创建:", result.ID)
	} else {
		fmt.Fprintln(os.Stderr, "创建失败:", resp.Status)
	}
}

// ── push ──
func cmdPush() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "用法: ypb push <file.md>")
		os.Exit(1)
	}

	filePath := os.Args[2]
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "读取文件失败: %v\n", err)
		os.Exit(1)
	}

	fm, body := parseFrontmatter(string(content))

	// 拼接请求
	payload := map[string]interface{}{
		"content":   body,
		"published": fm.published,
	}
	if fm.title != "" {
		payload["title"] = fm.title
	}
	if fm.slug != "" {
		payload["slug"] = fm.slug
	}
	if fm.summary != "" {
		payload["summary"] = fm.summary
	}
	if fm.category != "" {
		payload["category_id"] = fm.category
	}
	if len(fm.tags) > 0 {
		payload["tag_ids"] = fm.tags // 需提前创建好标签
	}

	data, _ := json.Marshal(payload)
	resp := apiRequest("POST", "/api/admin/posts", data)
	var result struct{ ID string }
	json.NewDecoder(resp.Body).Decode(&result)
	resp.Body.Close()

	if resp.StatusCode == 201 {
		fmt.Println("推送成功:", result.ID)
	} else {
		fmt.Fprintln(os.Stderr, "推送失败:", resp.Status)
	}
}

// ── pull ──
func cmdPull() {
	resp := apiRequest("GET", "/api/admin/posts?size=50", nil)
	var result struct {
		Posts []struct {
			ID        string `json:"id"`
			Title     string `json:"title"`
			Slug      string `json:"slug"`
			Content   string `json:"content"`
			Summary   string `json:"summary"`
			Published bool   `json:"published"`
			CreatedAt string `json:"created_at"`
			UpdatedAt string `json:"updated_at"`
		} `json:"posts"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	resp.Body.Close()

	dir := "_posts"
	os.MkdirAll(dir, 0755)

	for _, p := range result.Posts {
		filename := filepath.Join(dir, p.Slug+".md")
		fm := fmt.Sprintf("---\ntitle: %s\nslug: %s\npublished: %t\ncreated: %s\nupdated: %s\n---\n\n",
			p.Title, p.Slug, p.Published, p.CreatedAt, p.UpdatedAt)
		os.WriteFile(filename, []byte(fm+p.Content), 0644)
		fmt.Println("  →", filename)
	}
	fmt.Printf("拉取完成: %d 篇文章 → %s/\n", len(result.Posts), dir)
}

// ── list ──
func cmdList() {
	resp := apiRequest("GET", "/api/admin/posts?size=50", nil)
	var result struct {
		Posts []struct {
			ID        string `json:"id"`
			Title     string `json:"title"`
			Published bool   `json:"published"`
			CreatedAt string `json:"created_at"`
		} `json:"posts"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	resp.Body.Close()

	fmt.Printf("%-36s  %-6s  %-12s  %s\n", "ID", "状态", "日期", "标题")
	fmt.Println(strings.Repeat("-", 90))
	for _, p := range result.Posts {
		status := "草稿"
		if p.Published {
			status = "已发布"
		}
		date := p.CreatedAt[:10]
		fmt.Printf("%-36s  %-6s  %-12s  %s\n", p.ID[:8]+"...", status, date, truncate(p.Title, 30))
	}
}

// ── publish ──
func cmdPublish() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "用法: ypb publish <id>")
		os.Exit(1)
	}
	id := os.Args[2]
	resp := apiRequest("PUT", "/api/admin/posts/"+id+"/publish", nil)
	resp.Body.Close()
	if resp.StatusCode == 204 {
		fmt.Println("已发布:", id)
	} else {
		fmt.Fprintln(os.Stderr, "发布失败:", resp.Status)
	}
}

// ── delete ──
func cmdDelete() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "用法: ypb delete <id>")
		os.Exit(1)
	}
	id := os.Args[2]
	resp := apiRequest("DELETE", "/api/admin/posts/"+id, nil)
	resp.Body.Close()
	if resp.StatusCode == 204 {
		fmt.Println("已删除:", id)
	} else {
		fmt.Fprintln(os.Stderr, "删除失败:", resp.Status)
	}
}

// ── 工具函数 ──

func apiRequest(method, path string, body []byte) *http.Response {
	tok := requireToken()
	var r io.Reader
	if body != nil {
		r = bytes.NewReader(body)
	}
	req, _ := http.NewRequest(method, baseURL+path, r)
	req.Header.Set("Authorization", "Bearer "+tok)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "请求失败: %v\n", err)
		os.Exit(1)
	}
	return resp
}

func loadToken() {
	if t := os.Getenv("YPB_TOKEN"); t != "" {
		token = t
		return
	}
	data, err := os.ReadFile(configPath)
	if err == nil {
		token = strings.TrimSpace(string(data))
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return "."
}

func truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) > n {
		return string(runes[:n]) + "..."
	}
	return s
}

type frontmatter struct {
	title     string
	slug      string
	summary   string
	category  string
	tags      []string
	published bool
}

func parseFrontmatter(text string) (frontmatter, string) {
	var fm frontmatter
	text = strings.TrimSpace(text)
	if !strings.HasPrefix(text, "---") {
		return fm, text
	}
	idx := strings.Index(text[3:], "---")
	if idx < 0 {
		return fm, text
	}
	fmBlock := text[3 : idx+3]
	body := strings.TrimSpace(text[idx+6:])

	for _, line := range strings.Split(fmBlock, "\n") {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		k := strings.TrimSpace(parts[0])
		v := strings.TrimSpace(parts[1])
		switch k {
		case "title":
			fm.title = v
		case "slug":
			fm.slug = v
		case "summary":
			fm.summary = v
		case "category":
			fm.category = v
		case "tags":
			// 简单逗号分隔
			for _, t := range strings.Split(v, ",") {
				t = strings.TrimSpace(t)
				if t != "" {
					fm.tags = append(fm.tags, t)
				}
			}
		case "published":
			fm.published, _ = strconv.ParseBool(v)
		}
	}
	return fm, body
}
