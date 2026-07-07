<template>
  <div>
    <section class="hero home-hero" style="margin-bottom:28px">
      <p class="kicker">Archive</p>
      <h1 class="title home-lines">
        <span>所有.</span>
        <span class="is-muted">文章.</span>
      </h1>
      <p class="desc">按时间线浏览所有已发布文章。</p>
      <div class="meta"><span>共 {{ totalPosts }} 篇文章</span></div>
    </section>

    <div v-if="loading" class="post-card"><p class="muted">加载中...</p></div>
    <div v-else class="post-list">
      <template v-for="group in groupedPosts" :key="group.label">
        <h3 class="archive-year">{{ group.label }}</h3>
        <article v-for="p in group.posts" :key="p.id" class="post-card" style="padding:16px 24px">
          <div class="archive-row">
            <span class="archive-date">{{ shortDate(p.created_at) }}</span>
            <router-link :to="`/post/${p.slug}`">{{ p.title }}</router-link>
          </div>
        </article>
      </template>
    </div>

    <div v-if="totalPages > 1" class="pager">
      <a :class="{ disabled: page <= 1 }" @click="fetchAll(page - 1)">上一页</a>
      <span>{{ page }} / {{ totalPages }} 页</span>
      <a :class="{ disabled: page >= totalPages }" @click="fetchAll(page + 1)">下一页</a>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { posts } from '../api'

const allPosts = ref([])
const loading = ref(true)
const totalPosts = ref(0)
const page = ref(1)
const pageSize = 50
const totalPages = computed(() => Math.max(1, Math.ceil(totalPosts.value / pageSize)))

const groupedPosts = computed(() => {
  const groups = {}
  for (const p of allPosts.value) {
    const d = new Date(p.created_at)
    const label = d.getFullYear() + ' 年 ' + (d.getMonth() + 1) + ' 月'
    if (!groups[label]) groups[label] = []
    groups[label].push(p)
  }
  return Object.entries(groups).map(([label, posts]) => ({ label, posts }))
})

const shortDate = (s) => {
  if (!s) return ''
  const d = new Date(s)
  return d.getMonth() + 1 + '-' + d.getDate()
}

async function fetchAll(p = 1) {
  loading.value = true
  try {
    const data = await posts.list({ page: p, page_size: pageSize })
    allPosts.value = data.posts || []
    totalPosts.value = data.total || 0
    page.value = p
  } finally {
    loading.value = false
    window.scrollTo({ top: 0 })
  }
}

onMounted(() => fetchAll())
</script>

<style scoped>
.archive-year {
  font-family: "Fraunces", Georgia, serif;
  font-size: 1.3rem;
  margin: 1.5rem 0 0.5rem;
  color: #94a3b8;
}
.archive-row {
  display: flex;
  gap: 1rem;
  align-items: baseline;
}
.archive-date {
  color: var(--muted);
  font-size: 0.85rem;
  min-width: 3rem;
}
.archive-row a {
  text-decoration: none;
  color: var(--text);
}
.archive-row a:hover { color: var(--accent-soft); }
</style>
