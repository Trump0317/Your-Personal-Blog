<template>
  <div id="app">
    <header>
      <div class="wrap top">
        <router-link to="/" class="brand">{{ site.name }}.</router-link>
        <nav class="nav">
          <router-link to="/">首页</router-link>
          <router-link to="/archive">归档</router-link>
          <router-link to="/categories">分类</router-link>
          <router-link to="/about">关于</router-link>
          <button class="dark-toggle" @click="toggleDark" :title="isDark ? '亮色' : '暗色'">
            {{ isDark ? '☀' : '☾' }}
          </button>
        </nav>
      </div>
    </header>
    <main>
      <div class="wrap">
        <router-view v-slot="{ Component }">
          <component :is="Component" />
        </router-view>
      </div>
    </main>
    <footer>
      <div class="wrap">
        <div class="site-footer-top">
          <div>
            <div class="site-footer-title">{{ site.name }}</div>
            <p class="site-footer-description">{{ site.description }}</p>
          </div>
          <div class="site-footer-links">
            <router-link to="/">首页</router-link>
            <router-link to="/archive">归档</router-link>
            <router-link to="/about">关于</router-link>
            <a href="/rss.xml" target="_blank">RSS</a>
          </div>
        </div>
        <div class="site-footer-bottom">
          <span>&copy; {{ year }} Your Personal Blog.</span>
          <span class="site-footer-note">{{ site.footer_note }}</span>
        </div>
      </div>
    </footer>

    <!-- Toast -->
    <div class="toast-container">
      <div v-for="t in toasts" :key="t.id" :class="['toast', 'toast-' + t.type]">
        {{ t.message }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { useDark } from './composables/useDark.js'
import { useToast } from './composables/useToast.js'
import { useSite } from './composables/useSite.js'
import { provide, onMounted } from 'vue'

const { isDark, toggle: toggleDark } = useDark()
const { toasts, show } = useToast()
const { site, fetchSite } = useSite()

provide('toast', show)
onMounted(fetchSite)

const year = new Date().getFullYear()
</script>
