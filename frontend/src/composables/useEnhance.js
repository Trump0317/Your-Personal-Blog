/**
 * HTML 增强：Mermaid 图表 + KaTeX 数学公式
 * 在页面/组件 mounted 后调用，处理 v-html 渲染的内容
 */
import { nextTick } from 'vue'

/**
 * @param {HTMLElement|string} el - 目标元素或选择器
 */
export async function enhance(el) {
  if (typeof el === 'string') el = document.querySelector(el)
  if (!el) return

  await nextTick()

  // ── Mermaid: 找到 class="language-mermaid" 的 <code> 块 ──
  const mermaidBlocks = el.querySelectorAll('code.language-mermaid')
  for (const block of mermaidBlocks) {
    const pre = block.parentElement
    const id = 'mermaid-' + Math.random().toString(36).slice(2, 8)
    try {
      const { svg } = await window.mermaid.render(id, block.textContent)
      const div = document.createElement('div')
      div.className = 'mermaid-rendered'
      div.innerHTML = svg
      pre.replaceWith(div)
    } catch {
      pre.classList.add('mermaid-error')
    }
  }

  // ── KaTeX: 渲染 $...$ 和 $$...$$ ──
  if (window.renderMathInElement) {
    try {
      window.renderMathInElement(el, {
        delimiters: [
          { left: '$$', right: '$$', display: true },
          { left: '$', right: '$', display: false },
        ],
        throwOnError: false,
      })
    } catch { /* ignore */ }
  }
}
