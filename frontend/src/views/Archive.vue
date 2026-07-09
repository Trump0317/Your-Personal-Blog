<template>
  <div>
    <section class="hero">
      <p class="kicker">Timeline</p>
      <h1 class="title">文章归档</h1>
      <p class="desc">沿着时间轴回看每篇文章。</p>
    </section>

    <div v-if="loading" class="card"><p class="muted">加载中...</p></div>
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
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { archive } from '../api'

const years = ref([])
const loading = ref(true)

onMounted(async () => {
  try {
    years.value = await archive.stats()
  } catch {
    years.value = []
  } finally {
    loading.value = false
  }
})
</script>
