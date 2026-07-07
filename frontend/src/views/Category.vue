<template>
  <div>
    <div v-if="loading" class="post-card"><p class="muted">加载中...</p></div>
    <div v-else class="post-list">
      <article v-for="p in posts" :key="p.id" class="post-card">
        <h2>
          <router-link :to="`/post/${p.slug}`">{{ p.title }}</router-link>
        </h2>
        <p class="post-excerpt">{{ p.summary }}</p>
        <div class="post-row">
          <span>{{ fmtDate(p.created_at) }}</span>
          <span v-for="t in p.tags" :key="t.id">{{ t.name }}</span>
        </div>
      </article>
    </div>
    <div v-if="posts.length === 0 && !loading" class="post-card">
      <p class="muted">暂无文章。</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { posts, categories } from '../api'

const route = useRoute()
const postsList = ref([])
const loading = ref(true)

const fmtDate = (s) => s ? new Date(s).toLocaleDateString('zh-CN') : ''

onMounted(async () => {
  try {
    // slug → id 转换
    const cats = await categories.list()
    const cat = (cats || []).find(c => c.slug === route.params.slug)
    if (cat) {
      const data = await posts.list({ category: cat.id, page_size: 50 })
      postsList.value = data.posts || []
      document.title = cat.name + ' - My Blog'
    }
  } finally {
    loading.value = false
  }
})
</script>
