<template>
  <div>
    <div v-if="loading" class="card"><p class="muted">加载中...</p></div>
    <template v-else-if="tag">
      <section class="hero">
        <p class="kicker">Tag</p>
        <h1 class="title">{{ tag.name }}</h1>
        <p class="desc">共 {{ total }} 篇文章。</p>
      </section>

      <div class="list">
        <article v-for="p in posts" :key="p.id" class="card">
          <h2><router-link :to="`/post/${p.slug}`">{{ p.title }}</router-link></h2>
          <p class="excerpt">{{ p.summary }}</p>
          <div class="row">
            <span>{{ fmtDate(p.created_at) }}</span>
            <span v-if="p.category">
              <router-link :to="`/categories/${p.category.slug}`">{{ p.category.name }}</router-link>
            </span>
          </div>
        </article>
      </div>

      <div v-if="totalPages > 1" class="pager">
        <a :class="{ disabled: page <= 1 }" @click="go(page - 1)">上一页</a>
        <span>{{ page }} / {{ totalPages }} 页</span>
        <a :class="{ disabled: page >= totalPages }" @click="go(page + 1)">下一页</a>
      </div>
    </template>
    <div v-else-if="!loading" class="card"><p class="muted">标签不存在。</p></div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { posts, tags } from '../api'

const route = useRoute()
const tag = ref(null)
const postsList = ref([])
const loading = ref(true)
const page = ref(1)
const total = ref(0)
const size = 12

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / size)))
const fmtDate = s => s ? new Date(s).toLocaleDateString('zh-CN') : ''

async function fetchPosts(tagId, p = 1) {
  loading.value = true
  try {
    const data = await posts.list({ page: p, size, tag_id: tagId })
    postsList.value = data.posts || []
    total.value = data.total || 0
    page.value = p
  } finally { loading.value = false }
}

function go(p) { if (p >= 1 && p <= totalPages.value) fetchPosts(tag.value.id, p) }

onMounted(async () => {
  try {
    const tagList = await tags.list()
    const t = (tagList || []).find(x => x.slug === route.params.slug)
    if (t) { tag.value = t; document.title = t.name + ' 标签 - My Blog'; await fetchPosts(t.id) }
  } catch { loading.value = false }
})
</script>
