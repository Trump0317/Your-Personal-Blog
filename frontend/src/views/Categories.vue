<template>
  <div>
    <section class="hero">
      <p class="kicker">Taxonomy</p>
      <h1 class="title">全部分类</h1>
      <p class="desc">共 {{ cats.length }} 个技术主题。</p>
    </section>

    <div v-if="loading" class="card"><p class="muted">加载中...</p></div>
    <div v-else class="grid two">
      <router-link
        v-for="c in cats"
        :key="c.id"
        :to="`/categories/${c.slug}`"
        class="pill"
      >
        <span>{{ c.name }}</span>
        <span class="muted">{{ c.post_count || 0 }}</span>
      </router-link>
    </div>
    <p v-if="cats.length === 0 && !loading" class="muted">暂无分类。</p>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { categories } from '../api'

const cats = ref([])
const loading = ref(true)

onMounted(async () => {
  try {
    cats.value = await categories.list()
  } catch {
    cats.value = []
  } finally {
    loading.value = false
  }
})
</script>
