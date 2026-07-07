<template>
  <div>
    <div v-if="loading" class="post-card"><p class="muted">加载中...</p></div>
    <div v-else-if="error" class="post-card"><p class="error">{{ error }}</p></div>
    <article v-else-if="post" class="article">
      <h1>{{ post.title }}</h1>
      <div class="post-row">
        <span>发布于 {{ fmtDate(post.created_at) }}</span>
        <span v-if="post.category">{{ post.category.name }}</span>
        <span v-for="t in post.tags" :key="t.id">{{ t.name }}</span>
      </div>
      <div class="html-content" v-html="post.html_content || plainText(post.content)"></div>
    </article>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { posts } from '../api'

const route = useRoute()
const post = ref(null)
const loading = ref(true)
const error = ref('')

const fmtDate = (s) => s ? new Date(s).toLocaleDateString('zh-CN') : ''

function plainText(text) {
  return text?.replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/\n/g, '<br>') || ''
}

onMounted(async () => {
  try {
    post.value = await posts.get(route.params.slug)
    document.title = post.value.title + ' - My Blog'
  } catch (e) {
    error.value = '文章不存在'
  } finally {
    loading.value = false
  }
})
</script>
