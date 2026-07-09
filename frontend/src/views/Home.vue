<template>
  <div class="home">
    <!-- Hero -->
    <section class="hero home-hero">
      <p class="kicker">Daily Engineering Notes</p>
      <h1 class="title home-lines">
        <span>代码.</span>
        <span>匠心.</span>
        <span class="is-muted">文化.</span>
      </h1>
      <p class="desc">探索软件架构、工程美学与人类体验的交汇点。以工匠精神，雕琢每一行代码。</p>
      <div class="meta">
        <span>共 {{ siteStats.post_count || 0 }} 篇文章</span>
        <span>分类 {{ siteStats.category_count || 0 }} 个</span>
      </div>
      <div class="home-actions">
        <a href="#latest-posts" class="primary">查看最新文章</a>
        <router-link to="/categories" class="secondary">浏览主题</router-link>
      </div>
      <div class="home-trust">
        <span class="trust-badge">无广告</span>
        <span class="trust-badge">无付费软文</span>
        <span class="trust-badge">支持公开勘误</span>
      </div>
    </section>

    <!-- 精选主题 -->
    <section v-if="featuredCats.length">
      <div class="section-head">
        <div>
          <h2 class="section-title">精选主题</h2>
          <p class="section-description">核心技术的深度探索。</p>
        </div>
      </div>
      <div class="featured-grid">
        <router-link
          v-for="cat in featuredCats"
          :key="cat.id"
          :to="`/categories/${cat.slug}`"
          class="featured-card"
        >
          <span class="featured-label">{{ (cat.name || '?')[0].toUpperCase() }}</span>
          <h3>{{ cat.name }}</h3>
          <p>{{ cat.name }} 主题下的工程实践与深度文章。</p>
          <div class="featured-count">{{ cat.post_count || 0 }} 篇文章</div>
        </router-link>
      </div>
    </section>

    <!-- 信任声明 -->
    <section class="trust-grid">
      <article class="trust-card">
        <h3>广告策略</h3>
        <p>本站不投放展示广告，不嵌入联盟分成链接。</p>
      </article>
      <article class="trust-card">
        <h3>利益披露</h3>
        <p>若未来存在合作内容，会在标题区和文首显式标注。</p>
      </article>
      <article class="trust-card">
        <h3>内容勘误</h3>
        <p>发现错误可邮件反馈，确认后会修订并在文章中体现。</p>
      </article>
    </section>

    <!-- 最新文章 -->
    <section id="latest-posts">
      <div class="section-head">
        <div>
          <h2 class="section-title">最新见解</h2>
          <p class="section-description">近期的思考与工程笔记。</p>
        </div>
      </div>

      <div v-if="loading" class="card"><p class="muted">加载中...</p></div>
      <div v-else-if="error" class="card"><p class="error">{{ error }}</p></div>

      <div v-else class="list">
        <article v-for="p in posts" :key="p.id" class="card">
          <h2><router-link :to="`/post/${p.slug}`">{{ p.title }}</router-link></h2>
          <p class="excerpt">{{ p.summary || plainExcerpt(p.content, 120) }}</p>
          <div class="row">
            <span>{{ fmtDate(p.created_at) }}</span>
            <span v-if="p.category">
              <router-link :to="`/categories/${p.category.slug}`">{{ p.category.name }}</router-link>
            </span>
            <span>{{ monthLabel(p.created_at) }}</span>
          </div>
        </article>
      </div>

      <!-- 分页 -->
      <div v-if="totalPages > 1" class="pager">
        <a :class="{ disabled: page <= 1 }" @click="go(page - 1)">上一页</a>
        <a v-for="p in displayPages" :key="p" :class="{ active: p === page }" @click="go(p)">{{ p }}</a>
        <a :class="{ disabled: page >= totalPages }" @click="go(page + 1)">下一页</a>
        <span>第 {{ page }} / {{ totalPages }} 页</span>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { stats as statsApi, posts as postsApi, categories } from '../api'

const fmtDate = s => s ? new Date(s).toLocaleDateString('zh-CN') : ''
const monthLabel = s => { if (!s) return ''; const d = new Date(s); return d.getFullYear() + '-' + String(d.getMonth()+1).padStart(2,'0') }
const plainExcerpt = (t, n) => { if (!t) return ''; const p = t.replace(/[#*`>\[\]]/g,'').trim(); return p.length > n ? p.slice(0,n) + '...' : p }

const siteStats = ref({ post_count: 0, category_count: 0 })
const featuredCats = ref([])
const posts = ref([])
const loading = ref(true)
const error = ref('')
const page = ref(1)
const total = ref(0)
const size = 12

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / size)))

const displayPages = computed(() => {
  const pages = []
  const tp = totalPages.value, cp = page.value
  let s = Math.max(1, cp - 2), e = Math.min(tp, cp + 2)
  if (e - s < 4) { if (s === 1) e = Math.min(tp, 5); else s = Math.max(1, tp - 4) }
  for (let i = s; i <= e; i++) pages.push(i)
  return pages
})

async function fetchPosts(p = 1) {
  loading.value = true
  try {
    const data = await postsApi.list({ page: p, size })
    posts.value = data.posts || []
    total.value = data.total || 0
    page.value = p
  } catch (e) {
    error.value = '加载失败'
  } finally {
    loading.value = false
  }
}

function go(p) {
  if (p < 1 || p > totalPages.value) return
  fetchPosts(p)
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

onMounted(async () => {
  document.title = 'My Blog - 探索技术、设计与思考'
  await fetchPosts()
  try {
    const [s, c] = await Promise.all([statsApi.get(), categories.list()])
    siteStats.value = s
    featuredCats.value = (c || []).slice(0, 4)
  } catch { /* ignore */ }
})
</script>
