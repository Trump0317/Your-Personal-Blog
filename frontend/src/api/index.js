const BASE = '/api'
const ADMIN = `${BASE}/admin`

async function request(url, options = {}) {
  const res = await fetch(url, {
    headers: { 'Content-Type': 'application/json', ...options.headers },
    ...options,
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({}))
    throw new Error(err.error || `HTTP ${res.status}`)
  }
  if (res.status === 204) return null
  return res.json()
}

export const posts = {
  list(params = {}) {
    const qs = new URLSearchParams(params).toString()
    return request(`${BASE}/posts?${qs}`)
  },
  get(slug) {
    return request(`${BASE}/posts/${encodeURIComponent(slug)}`)
  },
}

export const categories = {
  list() { return request(`${BASE}/categories`) },
}

export const tags = {
  list() { return request(`${BASE}/tags`) },
}

// 管理端（需要 Basic Auth，浏览器凭据自动携带）
export const admin = {
  posts: {
    list() { return request(`${ADMIN}/posts`) },
    get(id) { return request(`${ADMIN}/posts/${id}`) },
    create(body) { return request(`${ADMIN}/posts`, { method: 'POST', body: JSON.stringify(body) }) },
    update(id, body) { return request(`${ADMIN}/posts/${id}`, { method: 'PUT', body: JSON.stringify(body) }) },
    delete(id) { return request(`${ADMIN}/posts/${id}`, { method: 'DELETE' }) },
    publish(id) { return request(`${ADMIN}/posts/${id}/publish`, { method: 'PUT' }) },
    unpublish(id) { return request(`${ADMIN}/posts/${id}/unpublish`, { method: 'PUT' }) },
  },
  categories: {
    list() { return request(`${ADMIN}/categories`) },
    create(body) { return request(`${ADMIN}/categories`, { method: 'POST', body: JSON.stringify(body) }) },
    update(id, body) { return request(`${ADMIN}/categories/${id}`, { method: 'PUT', body: JSON.stringify(body) }) },
    delete(id) { return request(`${ADMIN}/categories/${id}`, { method: 'DELETE' }) },
  },
  tags: {
    list() { return request(`${ADMIN}/tags`) },
    create(body) { return request(`${ADMIN}/tags`, { method: 'POST', body: JSON.stringify(body) }) },
    delete(id) { return request(`${ADMIN}/tags/${id}`, { method: 'DELETE' }) },
  },
}
