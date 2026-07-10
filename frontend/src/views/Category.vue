<template>
  <div>
    <div v-if="loading" class="card"><p class="muted">加载中...</p></div>
    <template v-else-if="category">
      <section class="hero">
        <p class="kicker">Category</p>
        <h1 class="title">{{ category.name }}</h1>
        <p class="desc">共 {{ total }} 篇文章。</p>
      </section>

      <div class="list">
        <article v-for="p in posts" :key="p.id" class="card">
          <h2><router-link :to="`/post/${p.slug}`">{{ p.title }}</router-link></h2>
          <p class="excerpt">{{ p.summary }}</p>
          <div class="row">
            <span>{{ fmtDate(p.created_at) }}</span>
            <span>{{ category.name }}</span>
            <span>{{ monthLabel(p.created_at) }}</span>
          </div>
        </article>
      </div>

      <div v-if="totalPages > 1" class="pager">
        <a :class="{ disabled: page <= 1 }" @click="go(page - 1)">上一页</a>
        <span>{{ page }} / {{ totalPages }} 页</span>
        <a :class="{ disabled: page >= totalPages }" @click="go(page + 1)">下一页</a>
      </div>
    </template>
    <div v-else-if="!loading" class="card"><p class="muted">分类不存在。</p></div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { posts, categories } from '../api'
import { useSite } from '../composables/useSite.js'

const { site } = useSite()
const route = useRoute()
const category = ref(null)
const postsList = ref([])
const loading = ref(true)
const page = ref(1)
const total = ref(0)
const size = 12

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / size)))

const fmtDate = s => s ? new Date(s).toLocaleDateString('zh-CN') : ''
const monthLabel = s => { if (!s) return ''; const d = new Date(s); return d.getFullYear() + '-' + String(d.getMonth()+1).padStart(2,'0') }

async function fetchPosts(catId, p = 1) {
  loading.value = true
  try {
    const data = await posts.list({ page: p, size, category_id: catId })
    postsList.value = data.posts || []
    total.value = data.total || 0
    page.value = p
  } finally {
    loading.value = false
  }
}

function go(p) { if (p >= 1 && p <= totalPages.value) fetchPosts(category.value.id, p) }

onMounted(async () => {
  try {
    const cat = await categories.get(route.params.slug)
    category.value = cat
    document.title = (cat.name || '分类') + ' - ' + site.value.name
    await fetchPosts(cat.id)
  } catch {
    category.value = null
    loading.value = false
  }
})
</script>
