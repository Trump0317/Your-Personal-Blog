<template>
  <div id="app-layout">
    <header>
      <div class="wrap top">
        <router-link to="/" class="brand">My Blog.</router-link>
        <nav class="nav">
          <router-link to="/">首页</router-link>
          <router-link to="/archive">归档</router-link>
          <router-link to="/admin">管理</router-link>
        </nav>
      </div>
    </header>
    <main>
      <div class="wrap">
        <router-view />
      </div>
    </main>
    <footer>
      <div class="wrap">
        <div class="site-footer-top">
          <div>
            <div class="site-footer-title">My Blog</div>
            <p class="site-footer-description">探索技术、设计与思考的交汇。</p>
          </div>
          <div class="site-footer-links">
            <router-link to="/">首页</router-link>
            <router-link to="/archive">归档</router-link>
            <a href="/api/posts" target="_blank">API</a>
          </div>
        </div>
        <div class="site-footer-bottom">
          <span>&copy; {{ year }} My Blog. Powered by Go + SQLite + Vue.</span>
          <span class="site-footer-note">无广告 · 无付费软文</span>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup>
const year = new Date().getFullYear()
</script>
