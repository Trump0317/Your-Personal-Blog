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
        <span>共 {{ totalPosts }} 篇文章</span>
        <span>分类 {{ featuredCats.length }} 个</span>
        <span>标签 {{ totalTags }} 个</span>
      </div>
      <div class="home-actions">
        <a href="#latest-posts" class="primary">查看最新文章</a>
        <router-link to="/archive" class="secondary">浏览归档</router-link>
      </div>
      <div class="home-trust">
        <span class="trust-badge">无广告</span>
        <span class="trust-badge">无付费软文</span>
        <span class="trust-badge">支持公开勘误</span>
      </div>
    </section>

    <!-- 精选主题 -->
    <section>
      <div class="section-head">
        <div>
          <h2 class="section-title">精选主题</h2>
          <p class="section-description">核心技术的深度探索。</p>
        </div>
      </div>
      <div class="featured-grid">
        <router-link
          v-for="(cat, i) in featuredCats"
          :key="cat.id"
          :to="`/category/${cat.slug}`"
          class="featured-card"
        >
          <span class="featured-label">{{ cat.name.charAt(0).toUpperCase() }}</span>
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

      <div v-if="loading" class="post-card"><p class="muted">加载中...</p></div>
      <div v-else-if="error" class="post-card"><p class="error">{{ error }}</p></div>

      <div v-else class="post-list">
        <article v-for="p in posts" :key="p.id" class="post-card">
          <h2>
            <router-link :to="`/post/${p.slug}`">{{ p.title }}</router-link>
          </h2>
          <p class="post-excerpt">{{ p.summary || extractExcerpt(p.content, 120) }}</p>
          <div class="post-row">
            <span>{{ fmtDate(p.created_at) }}</span>
            <span v-if="p.category">{{ p.category.name }}</span>
            <span v-for="t in p.tags" :key="t.id" class="tag-badge">{{ t.name }}</span>
          </div>
        </article>
      </div>

      <!-- 分页 -->
      <div v-if="totalPages > 1" class="pager">
        <a
          :class="{ disabled: currentPage <= 1 }"
          @click="goPage(currentPage - 1)"
        >上一页</a>
        <a
          v-for="p in displayPages"
          :key="p"
          :class="{ active: p === currentPage }"
          @click="goPage(p)"
        >{{ p }}</a>
        <a
          :class="{ disabled: currentPage >= totalPages }"
          @click="goPage(currentPage + 1)"
        >下一页</a>
        <span>{{ currentPage }} / {{ totalPages }} 页</span>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { posts as postsApi, categories, tags } from '../api'

const posts = ref([])
const loading = ref(true)
const error = ref('')
const currentPage = ref(1)
const totalPosts = ref(0)
const totalTags = ref(0)
const featuredCats = ref([])
const pageSize = 12

const totalPages = computed(() => Math.max(1, Math.ceil(totalPosts.value / pageSize)))

const displayPages = computed(() => {
  const pages = []
  const tp = totalPages.value
  const cp = currentPage.value
  let start = Math.max(1, cp - 2)
  let end = Math.min(tp, cp + 2)
  if (end - start < 4) {
    if (start === 1) end = Math.min(tp, 5)
    else start = Math.max(1, tp - 4)
  }
  for (let i = start; i <= end; i++) pages.push(i)
  return pages
})

const fmtDate = (s) => s ? new Date(s).toLocaleDateString('zh-CN') : ''
const extractExcerpt = (text, len) => {
  if (!text) return ''
  const plain = text.replace(/[#*`>\[\]]/g, '').trim()
  return plain.length > len ? plain.slice(0, len) + '...' : plain
}

async function fetchPosts(page = 1) {
  loading.value = true
  try {
    const data = await postsApi.list({ page })
    posts.value = data.posts || []
    totalPosts.value = data.total || 0
    currentPage.value = page
  } catch (e) {
    error.value = '加载失败: ' + e.message
  } finally {
    loading.value = false
  }
}

function goPage(page) {
  if (page < 1 || page > totalPages.value) return
  fetchPosts(page)
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

onMounted(async () => {
  await fetchPosts()
  try {
    const [catData, tagData] = await Promise.all([
      categories.list(),
      tags.list(),
    ])
    featuredCats.value = (catData || []).slice(0, 4)
    // 统计每个分类的文章数
    const allData = await postsApi.list({ page: 1, page_size: '1000' })
    const allPosts = allData.posts || []
    featuredCats.value = (catData || []).slice(0, 4).map(c => ({
      ...c,
      post_count: allPosts.filter(p => p.category_id === c.id).length,
    }))
    totalTags.value = (tagData || []).length
  } catch { /* ignore */ }
})
</script>
