<template>
  <div>
    <section class="hero">
      <p class="kicker">About</p>
      <h1 class="title">{{ site.about_title }}</h1>
      <p class="desc">{{ site.tagline }}</p>
    </section>

    <article class="article" v-html="aboutHTML"></article>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useSite } from '../composables/useSite.js'

const { site } = useSite()

const aboutHTML = computed(() => {
  return (site.value.about_text || '').split('\n').filter(l => l.trim()).map(l => `<p>${l}</p>`).join('')
})
</script>
