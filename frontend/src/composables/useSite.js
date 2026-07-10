import { ref } from 'vue'

const site = ref({
  name: 'Your Personal Blog',
  tagline: '探索技术、设计与思考的交汇',
  description: '探索技术、设计与思考的交汇。以极简致敬繁复。',
  footer_note: '探索技术、设计与思考的交汇',
  hero_kicker: 'Daily Engineering Notes',
  hero_title1: '代码.',
  hero_title2: '匠心.',
  hero_title3: '文化.',
  hero_desc: '探索软件架构、工程美学与人类体验的交汇点。以工匠精神，雕琢每一行代码。',
  about_title: '关于 Your Personal Blog',
  about_text: '这是一个用 Go + SQLite + Vue 构建的个人博客。静态分发优先，无运行时依赖。\n写博客的目的：记录思考，沉淀知识，分享实践。',
})

export function useSite() {
  async function fetchSite() {
    try {
      const res = await fetch('/api/site')
      if (res.ok) {
        const data = await res.json()
        site.value = data
      }
    } catch { /* use defaults */ }
  }

  return { site, fetchSite }
}
