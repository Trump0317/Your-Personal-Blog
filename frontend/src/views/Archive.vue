<template>
  <div>
    <section class="hero">
      <p class="kicker">Timeline</p>
      <h1 class="title">文章归档</h1>
      <p class="desc">沿着时间轴回看每篇文章。</p>
    </section>

    <div v-if="loading" class="card"><p class="muted">加载中...</p></div>
    <template v-else>
      <!-- 筛选模式：有 year 参数时显示该时段文章列表 -->
      <template v-if="filterYear">
        <div class="section-head">
          <div>
            <h2 class="section-title">
              {{ filterYear }} 年{{ filterMonth ? ' ' + filterMonth + ' 月' : '' }}
            </h2>
            <p class="section-description">共 {{ filteredPosts.length }} 篇文章</p>
          </div>
          <router-link to="/archive" class="back-link">← 返回归档总览</router-link>
        </div>
        <div v-if="filteredPosts.length" class="list">
          <article v-for="p in filteredPosts" :key="p.id" class="card">
            <h2><router-link :to="`/post/${p.slug}`">{{ p.title }}</router-link></h2>
            <p class="excerpt">{{ p.summary || plainExcerpt(p.content, 120) }}</p>
            <div class="row">
              <span>{{ fmtDate(p.created_at) }}</span>
              <span v-if="p.category">
                <router-link :to="`/categories/${p.category.slug}`">{{ p.category.name }}</router-link>
              </span>
            </div>
          </article>
        </div>
        <p v-else class="muted">该时段暂无文章。</p>
      </template>

      <!-- 总览模式 -->
      <template v-else>
        <section v-for="year in years" :key="year.year" class="archive-year-card">
          <h2>
            <router-link :to="`/archive?year=${year.year}`">{{ year.year }} 年</router-link>
          </h2>
          <p class="muted">{{ year.count }} 篇</p>
          <div class="grid">
            <router-link
              v-for="m in year.months"
              :key="m.month"
              :to="`/archive?year=${year.year}&month=${m.month}`"
              class="pill"
            >
              <span>{{ year.year }}-{{ String(m.month).padStart(2, '0') }}</span>
              <span class="muted">{{ m.count }}</span>
            </router-link>
          </div>
        </section>
        <p v-if="years.length === 0" class="muted">暂无归档。</p>
      </template>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { archive, posts as postsApi } from '../api'

const route = useRoute()
const years = ref([])
const filteredPosts = ref([])
const loading = ref(true)

const fmtDate = s => s ? new Date(s).toLocaleDateString('zh-CN') : ''
const plainExcerpt = (t, n) => { if (!t) return ''; const p = t.replace(/[#*`>\[\]]/g,'').trim(); return p.length > n ? p.slice(0,n) + '...' : p }

const filterYear = computed(() => {
  const y = parseInt(route.query.year)
  return y > 0 ? y : null
})
const filterMonth = computed(() => {
  const m = parseInt(route.query.month)
  return (m >= 1 && m <= 12) ? m : null
})

async function loadArchive() {
  try {
    years.value = await archive.stats()
  } catch {
    years.value = []
  }
}

async function loadFilteredPosts() {
  if (!filterYear.value) return
  try {
    const params = { year: filterYear.value, size: 50 }
    if (filterMonth.value) params.month = filterMonth.value
    const data = await postsApi.list(params)
    filteredPosts.value = data.posts || []
  } catch {
    filteredPosts.value = []
  }
}

async function load() {
  loading.value = true
  await loadArchive()
  await loadFilteredPosts()
  loading.value = false
}

onMounted(load)
watch(() => route.query, load)
</script>
