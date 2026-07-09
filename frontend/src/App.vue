<template>
  <div id="app">
    <header>
      <div class="wrap top">
        <router-link to="/" class="brand">My Blog.</router-link>
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
            <div class="site-footer-title">My Blog</div>
            <p class="site-footer-description">探索技术、设计与思考的交汇。<br/>以极简致敬繁复。</p>
          </div>
          <div class="site-footer-links">
            <router-link to="/">首页</router-link>
            <router-link to="/archive">归档</router-link>
            <router-link to="/about">关于</router-link>
            <a href="/rss.xml" target="_blank">RSS</a>
          </div>
        </div>
        <div class="site-footer-bottom">
          <span>&copy; {{ year }} My Blog.</span>
          <span class="site-footer-note">无广告 · 无付费软文 · 支持公开勘误</span>
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
import { provide } from 'vue'

const { isDark, toggle: toggleDark } = useDark()
const { toasts, show } = useToast()

provide('toast', show)

const year = new Date().getFullYear()
</script>
