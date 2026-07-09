import { ref, watchEffect } from 'vue'

const isDark = ref(false)

export function useDark() {
  // 初始化
  const saved = localStorage.getItem('blog-dark')
  if (saved !== null) {
    isDark.value = saved === 'true'
  } else {
    isDark.value = window.matchMedia('(prefers-color-scheme: dark)').matches
  }

  watchEffect(() => {
    document.documentElement.classList.toggle('dark', isDark.value)
    localStorage.setItem('blog-dark', isDark.value)
  })

  function toggle() { isDark.value = !isDark.value }

  return { isDark, toggle }
}
