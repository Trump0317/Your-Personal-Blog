/**
 * API 接口封装
 * 所有请求基础路径 /api，与 Vite proxy 配置一致
 */
const BASE = '/api'

async function request(url, opts = {}) {
  const res = await fetch(url, opts)
  if (!res.ok) {
    const e = await res.json().catch(() => ({}))
    throw new Error(e.error || `HTTP ${res.status}`)
  }
  if (res.status === 204) return null
  return res.json()
}

// ── 统计 ──
export const stats = {
  /** GET /api/stats → { post_count, category_count, tag_count } */
  get() { return request(`${BASE}/stats`) },
}

// ── 文章 ──
export const posts = {
  /**
   * GET /api/posts?page&size&category_id&tag_id&year&month&q
   * → { posts: Post[], total: number }
   */
  list(params = {}) {
    const qs = new URLSearchParams(
      Object.entries(params).filter(([, v]) => v !== undefined && v !== '')
    ).toString()
    return request(`${BASE}/posts?${qs}`)
  },

  /**
   * GET /api/posts/:slug
   * → Post { id, title, slug, content, html_content, summary,
   *          category: { id, name, slug },
   *          tags: [{ id, name, slug }],
   *          created_at, updated_at }
   */
  get(slug) {
    return request(`${BASE}/posts/${encodeURIComponent(slug)}`)
  },
}

// ── 分类 ──
export const categories = {
  /**
   * GET /api/categories
   * → [{ id, name, slug, post_count }]
   */
  list() { return request(`${BASE}/categories`) },

  /**
   * GET /api/categories/:slug
   * → { id, name, slug, post_count }
   */
  get(slug) { return request(`${BASE}/categories/${slug}`) },
}

// ── 标签 ──
export const tags = {
  /** GET /api/tags → [{ id, name, slug, post_count }] */
  list() { return request(`${BASE}/tags`) },

  /** GET /api/tags/:slug → { id, name, slug } */
  get(slug) { return request(`${BASE}/tags/${encodeURIComponent(slug)}`) },
}

// ── 归档统计 ──
export const archive = {
  /**
   * GET /api/archive
   * → [{ year: number, count: number, months: [{ month: number, count: number }] }]
   */
  stats() { return request(`${BASE}/archive`) },
}

// ── 图片上传 ──
export async function uploadFile(file) {
  const fd = new FormData()
  fd.append('file', file)
  const res = await fetch(`${BASE}/upload`, {
    method: 'POST',
    headers: { Authorization: authHeader },
    body: fd,
  })
  if (!res.ok) throw new Error('upload failed')
  return res.json()
}

// ── 管理端（需要 Basic Auth）──
let authHeader = ''

export function setAuth(user, pass) {
  authHeader = 'Basic ' + btoa(user + ':' + pass)
}

function authRequest(url, opts = {}) {
  return request(url, {
    ...opts,
    headers: { ...opts.headers, Authorization: authHeader },
  })
}

export const admin = {
  test() { return authRequest(`${BASE}/admin/posts?size=1`) },

  posts: {
    list(params) {
      const qs = new URLSearchParams(params).toString()
      return authRequest(`${BASE}/admin/posts?${qs}`)
    },
    get(id) { return authRequest(`${BASE}/admin/posts/${id}`) },
    create(body) {
      return authRequest(`${BASE}/admin/posts`, {
        method: 'POST', headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
      })
    },
    update(id, body) {
      return authRequest(`${BASE}/admin/posts/${id}`, {
        method: 'PUT', headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
      })
    },
    delete(id) { return authRequest(`${BASE}/admin/posts/${id}`, { method: 'DELETE' }) },
    publish(id) { return authRequest(`${BASE}/admin/posts/${id}/publish`, { method: 'PUT' }) },
    unpublish(id) { return authRequest(`${BASE}/admin/posts/${id}/unpublish`, { method: 'PUT' }) },
  },

  categories: {
    list() { return authRequest(`${BASE}/admin/categories`) },
    create(body) {
      return authRequest(`${BASE}/admin/categories`, {
        method: 'POST', headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
      })
    },
    update(id, body) {
      return authRequest(`${BASE}/admin/categories/${id}`, {
        method: 'PUT', headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
      })
    },
    delete(id) { return authRequest(`${BASE}/admin/categories/${id}`, { method: 'DELETE' }) },
  },

  tags: {
    list() { return authRequest(`${BASE}/admin/tags`) },
    create(body) {
      return authRequest(`${BASE}/admin/tags`, {
        method: 'POST', headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
      })
    },
    delete(id) { return authRequest(`${BASE}/admin/tags/${id}`, { method: 'DELETE' }) },
  },
}
