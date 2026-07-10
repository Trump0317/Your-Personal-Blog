<template>
  <div>
    <div v-if="loading" class="card"><p class="muted">加载中...</p></div>
    <div v-else-if="error" class="card"><p class="error">{{ error }}</p></div>
    <template v-else-if="post">
      <section class="hero">
        <p class="kicker">Article</p>
        <h1 class="title">{{ post.title }}</h1>
        <p class="desc">{{ post.summary }}</p>
        <div class="meta">
          <span>{{ fmtDate(post.created_at) }}</span>
          <span v-if="post.category">
            <router-link :to="`/categories/${post.category.slug}`">{{ post.category.name }}</router-link>
          </span>
          <span v-for="t in post.tags" :key="t.id">
            <router-link :to="`/tags/${t.slug}`">#{{ t.name }}</router-link>
          </span>
        </div>
      </section>
      <article class="article" v-html="post.html_content || plainText(post.content)"></article>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { posts } from '../api'
import { useSite } from '../composables/useSite.js'

const { site } = useSite()
const route = useRoute()
const post = ref(null)
const loading = ref(true)
const error = ref('')
const fmtDate = s => s ? new Date(s).toLocaleDateString('zh-CN') : ''
const plainText = t => t ? t.replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/\n/g,'<br>') : ''

onMounted(async () => {
  try {
    post.value = await posts.get(route.params.slug)
    document.title = post.value.title + ' - ' + site.value.name
  } catch (e) {
    error.value = '文章不存在'
  } finally {
    loading.value = false
  }
})
</script>
