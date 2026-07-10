<template>
  <div>
    <!-- ── 登录 ── -->
    <div v-if="!authed" class="login-box">
      <h2>管理员登录</h2>
      <form @submit.prevent="login">
        <label>用户名 <input v-model="authUser" required></label>
        <label>密码 <input v-model="authPass" type="password" required></label>
        <button type="submit">登录</button>
      </form>
      <p v-if="loginError" class="error">{{ loginError }}</p>
    </div>

    <div v-else>
      <!-- ── 文章列表工具栏 ── -->
      <div class="admin-header">
        <h2 class="section-title">文章管理</h2>
        <button @click="newPost">+ 新建文章</button>
      </div>

      <div v-if="!editing" class="card">
        <!-- 筛选栏 -->
        <div class="filter-bar">
          <input v-model="searchQuery" placeholder="搜索文章..." class="filter-input">
          <select v-model="statusFilter" class="filter-select">
            <option value="">全部状态</option>
            <option value="published">已发布</option>
            <option value="draft">草稿</option>
          </select>
          <select v-model="categoryFilter" class="filter-select">
            <option value="">全部分类</option>
            <option v-for="c in categories" :key="c.id" :value="c.id">{{ c.name }}</option>
          </select>
          <span class="muted" style="margin-left:auto">{{ filteredPosts.length }} 篇</span>
        </div>

        <div v-if="filteredPosts.length === 0" class="muted" style="padding:1.5rem 0;text-align:center">
          暂无匹配文章。
        </div>
        <div v-for="p in filteredPosts" :key="p.id" class="admin-post-item">
          <div class="post-info">
            <strong>{{ p.title }}</strong>
            <span class="post-meta">
              <span :class="p.published ? 'pub' : 'draft'">{{ p.published ? '已发布' : '草稿' }}</span>
              <span v-if="p.category" class="muted">{{ p.category.name }}</span>
              <span class="muted">{{ fmtDate(p.created_at) }}</span>
            </span>
          </div>
          <span class="post-actions">
            <button @click="editPost(p)">编辑</button>
            <button @click="togglePublish(p)" class="secondary">{{ p.published ? '取消发布' : '发布' }}</button>
            <button @click="deletePost(p.id)" class="danger">删除</button>
          </span>
        </div>
      </div>

      <!-- ── 编辑器 ── -->
      <div v-else class="editor-container">
        <div class="card editor-meta">
          <h3>{{ formId ? '编辑文章' : '新建文章' }}
            <span v-if="draftSaved" class="auto-save">✓ 草稿已保存</span>
          </h3>
          <div class="meta-row">
            <label class="meta-field flex-2">
              标题 <input v-model="form.title" required @input="scheduleDraft">
            </label>
            <label class="meta-field">
              Slug <input v-model="form.slug" placeholder="留空自动生成" @input="scheduleDraft">
            </label>
          </div>
          <div class="meta-row">
            <label class="meta-field">
              分类 <select v-model="form.category_id" @change="scheduleDraft">
                <option value="">无分类</option>
                <option v-for="c in categories" :key="c.id" :value="c.id">{{ c.name }}</option>
              </select>
            </label>
            <label class="meta-field">
              摘要 <input v-model="form.summary" placeholder="留空自动提取" @input="scheduleDraft">
            </label>
          </div>
          <div class="tag-checkboxes">
            <label v-for="t in tagsList" :key="t.id" class="tag-label">
              <input type="checkbox" :value="t.id" v-model="form.tag_ids" @change="scheduleDraft"> {{ t.name }}
            </label>
          </div>
        </div>

        <!-- 工具栏 -->
        <div class="toolbar">
          <button type="button" class="tool-btn" title="粗体 (Ctrl+B)" @click="insertMark('**', '**')"><b>B</b></button>
          <button type="button" class="tool-btn" title="斜体 (Ctrl+I)" @click="insertMark('*', '*')"><i>I</i></button>
          <button type="button" class="tool-btn" title="标题" @click="insertLine('## ')">H</button>
          <button type="button" class="tool-btn" title="链接" @click="insertLink">🔗</button>
          <button type="button" class="tool-btn" title="行内代码" @click="insertMark('`', '`')">&lt;/&gt;</button>
          <button type="button" class="tool-btn" title="代码块" @click="insertBlock('```')">```</button>
          <button type="button" class="tool-btn" title="引用" @click="insertLine('> ')">❝</button>
          <button type="button" class="tool-btn" title="列表" @click="insertLine('- ')">≡</button>
          <span class="toolbar-spacer"></span>
          <input type="file" accept="image/*" @change="uploadImage" ref="fileInput" style="display:none">
          <button type="button" class="tool-btn" title="上传图片" @click="fileInput.click()">🖼</button>
          <span v-if="uploading" class="muted" style="font-size:0.8rem">上传中...</span>
        </div>

        <!-- 分栏编辑 -->
        <div class="split-pane">
          <div class="pane pane-edit">
            <textarea
              ref="editorArea"
              v-model="form.content"
              placeholder="在此编写 Markdown..."
              @input="onContentInput"
              @keydown="onKeydown"
            ></textarea>
          </div>
          <div class="pane pane-preview">
            <div class="preview-content" v-html="previewHTML"></div>
          </div>
        </div>

        <div class="card editor-footer">
          <label class="publish-check">
            <input type="checkbox" v-model="form.published"> 发布
          </label>
          <div class="form-actions">
            <button type="button" @click="save" class="primary">保存</button>
            <button type="button" @click="cancelEdit" class="secondary">取消</button>
          </div>
        </div>
      </div>

      <!-- ── 分类 & 标签管理 ── -->
      <div v-if="!editing" class="meta-panels">
        <div class="meta-panel">
          <h3>分类</h3>
          <form @submit.prevent="addCategory">
            <input v-model="newCat" placeholder="新分类名称" required>
            <button type="submit" class="sm">添加</button>
          </form>
          <div v-for="c in categories" :key="c.id" class="meta-item">
            <template v-if="renamingCat === c.id">
              <input v-model="renameCatName" @keyup.enter="saveRenameCat(c.id)" @keyup.escape="renamingCat = ''" class="rename-input">
              <button @click="saveRenameCat(c.id)" class="sm">✓</button>
            </template>
            <template v-else>
              <span>{{ c.name }}</span>
              <span>
                <button @click="startRenameCat(c)" class="sm secondary" title="重命名">✎</button>
                <button @click="delCategory(c.id)" class="danger sm">×</button>
              </span>
            </template>
          </div>
        </div>
        <div class="meta-panel">
          <h3>标签</h3>
          <form @submit.prevent="addTag">
            <input v-model="newTag" placeholder="新标签名称" required>
            <button type="submit" class="sm">添加</button>
          </form>
          <div v-for="t in tagsList" :key="t.id" class="meta-item">
            <template v-if="renamingTag === t.id">
              <input v-model="renameTagName" @keyup.enter="saveRenameTag(t.id)" @keyup.escape="renamingTag = ''" class="rename-input">
              <button @click="saveRenameTag(t.id)" class="sm">✓</button>
            </template>
            <template v-else>
              <span>{{ t.name }}</span>
              <span>
                <button @click="startRenameTag(t)" class="sm secondary" title="重命名">✎</button>
                <button @click="delTag(t.id)" class="danger sm">×</button>
              </span>
            </template>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, inject, onBeforeUnmount, nextTick } from 'vue'
import { setAuth, admin, uploadFile } from '../api'

const toast = inject('toast')

// ── 认证 ──
const authed = ref(false)
const authUser = ref('')
const authPass = ref('')
const loginError = ref('')

function login() {
  setAuth(authUser.value, authPass.value)
  admin.test().then(() => {
    authed.value = true; loginError.value = ''
    localStorage.setItem('blog_auth', btoa(authUser.value + ':' + authPass.value))
    loadAll()
  }).catch(() => loginError.value = '登录失败')
}

// ── 数据 ──
const postsList = ref([])
const categories = ref([])
const tagsList = ref([])
const editing = ref(false)
const formId = ref('')
const newCat = ref('')
const newTag = ref('')
const uploading = ref(false)
const fileInput = ref(null)
const editorArea = ref(null)
const draftSaved = ref(false)

const form = reactive({ title: '', slug: '', content: '', summary: '', category_id: '', tag_ids: [], published: false })

// ── 筛选 ──
const searchQuery = ref('')
const statusFilter = ref('')
const categoryFilter = ref('')

const filteredPosts = computed(() => {
  let list = postsList.value
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    list = list.filter(p => p.title.toLowerCase().includes(q))
  }
  if (statusFilter.value === 'published') list = list.filter(p => p.published)
  if (statusFilter.value === 'draft') list = list.filter(p => !p.published)
  if (categoryFilter.value) list = list.filter(p => p.category_id === categoryFilter.value)
  return list
})

function fmtDate(s) { return s ? new Date(s).toLocaleDateString('zh-CN') : '' }

// ── 数据加载 ──
async function loadAll() {
  const [p, c, t] = await Promise.all([admin.posts.list({ size: 50 }), admin.categories.list(), admin.tags.list()])
  postsList.value = p.posts || []
  // 补充 category 名称用于列表展示
  for (const post of postsList.value) {
    if (post.category_id && !post.category) {
      post.category = (c || []).find(x => x.id === post.category_id)
    }
  }
  categories.value = c || []; tagsList.value = t || []
}

// ── 编辑器 ──
function newPost() {
  clearDraft()
  formId.value = ''; editing.value = true
  Object.assign(form, { title: '', slug: '', content: '', summary: '', category_id: '', tag_ids: [], published: false })
  draftSaved.value = false
  nextTick(() => editorArea.value?.focus())
}

function editPost(p) {
  clearDraft()
  formId.value = p.id; editing.value = true
  Object.assign(form, { title: p.title, slug: p.slug, content: p.content, summary: p.summary || '', category_id: p.category_id || '', tag_ids: p.tag_ids || [], published: p.published })
  draftSaved.value = false
}

function cancelEdit() {
  if (form.content && !formId.value) scheduleDraft() // 新建时取消也保存草稿
  editing.value = false
}

async function save() {
  const body = { title: form.title, slug: form.slug || undefined, content: form.content, summary: form.summary, category_id: form.category_id, tag_ids: form.tag_ids, published: form.published }
  try {
    if (formId.value) { await admin.posts.update(formId.value, body) } else { await admin.posts.create(body) }
    toast('保存成功', 'success')
    clearDraft()
    editing.value = false; await loadAll()
  } catch (e) { toast('保存失败: ' + e.message, 'error') }
}

async function togglePublish(p) {
  if (p.published) await admin.posts.unpublish(p.id); else await admin.posts.publish(p.id)
  toast(p.published ? '已取消发布' : '已发布', 'success')
  await loadAll()
}

async function deletePost(id) {
  if (!confirm('确认删除？删除后不可恢复。')) return
  await admin.posts.delete(id); await loadAll()
  toast('已删除', 'info')
}

// ── Markdown 预览（简易渲染）──
const previewHTML = computed(() => {
  if (!form.content) return '<p class="muted">预览将显示在这里...</p>'
  return renderPreview(form.content)
})

function renderPreview(md) {
  if (!md) return ''
  let html = md
    .replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
  // 代码块（优先处理，防止内部标记被转义）
  html = html.replace(/```(\w*)\n([\s\S]*?)```/g, (_, lang, code) =>
    `<pre><code class="language-${lang}">${code.replace(/&amp;/g,'&').replace(/&lt;/g,'<').replace(/&gt;/g,'>')}</code></pre>`)
  // 行内代码
  html = html.replace(/`([^`]+)`/g, '<code>$1</code>')
  // 标题
  html = html.replace(/^### (.+)$/gm, '<h3>$1</h3>')
  html = html.replace(/^## (.+)$/gm, '<h2>$1</h2>')
  html = html.replace(/^# (.+)$/gm, '<h1>$1</h1>')
  // 粗斜体
  html = html.replace(/\*\*\*(.+?)\*\*\*/g, '<strong><em>$1</em></strong>')
  html = html.replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>')
  html = html.replace(/\*(.+?)\*/g, '<em>$1</em>')
  // 图片
  html = html.replace(/!\[([^\]]*)\]\(([^)]+)\)/g, '<img src="$2" alt="$1">')
  // 链接
  html = html.replace(/\[([^\]]+)\]\(([^)]+)\)/g, '<a href="$2" target="_blank">$1</a>')
  // 引用
  html = html.replace(/^&gt; (.+)$/gm, '<blockquote>$1</blockquote>')
  // 无序列表
  html = html.replace(/^- (.+)$/gm, '<li>$1</li>')
  html = html.replace(/(<li>.*<\/li>\n?)+/g, '<ul>$&</ul>')
  // 水平线
  html = html.replace(/^---$/gm, '<hr>')
  // 段落
  html = html.replace(/^(?!<[a-z/])(.+)$/gm, '<p>$1</p>')
  // 合并连续块
  html = html.replace(/<\/blockquote>\n<blockquote>/g, '\n')
  html = html.replace(/<\/pre>\n?<pre>/g, '\n')
  return html
}

// ── 工具栏 ──
function insertMark(before, after) {
  const ta = editorArea.value; if (!ta) return
  const { selectionStart: s, selectionEnd: e, value } = ta
  const sel = value.slice(s, e)
  form.content = value.slice(0, s) + before + sel + after + value.slice(e)
  nextTick(() => { ta.focus(); ta.setSelectionRange(s + before.length, s + before.length + sel.length) })
}

function insertLine(prefix) {
  const ta = editorArea.value; if (!ta) return
  const { selectionStart: s, value } = ta
  const lineStart = value.lastIndexOf('\n', s - 1) + 1
  form.content = value.slice(0, lineStart) + prefix + value.slice(lineStart)
  nextTick(() => { ta.focus(); ta.setSelectionRange(lineStart + prefix.length, lineStart + prefix.length) })
}

function insertBlock(fence) {
  const ta = editorArea.value; if (!ta) return
  const { selectionStart: s, selectionEnd: e, value } = ta
  const sel = value.slice(s, e)
  form.content = value.slice(0, s) + fence + '\n' + sel + '\n' + fence + value.slice(e)
  nextTick(() => { ta.focus(); ta.setSelectionRange(s + fence.length + 1, s + fence.length + 1 + sel.length) })
}

async function insertLink() {
  const url = prompt('输入链接 URL:', 'https://')
  if (!url) return
  const ta = editorArea.value; if (!ta) return
  const { selectionStart: s, selectionEnd: e, value } = ta
  const sel = value.slice(s, e) || 'link'
  form.content = value.slice(0, s) + `[${sel}](${url})` + value.slice(e)
  nextTick(() => ta.focus())
}

function onKeydown(e) {
  if ((e.ctrlKey || e.metaKey) && e.key === 'b') { e.preventDefault(); insertMark('**', '**') }
  if ((e.ctrlKey || e.metaKey) && e.key === 'i') { e.preventDefault(); insertMark('*', '*') }
  if (e.key === 'Tab') { e.preventDefault(); insertMark('  ', '') }
}
function onContentInput() { scheduleDraft() }

// ── 图片上传 ──
async function uploadImage(e) {
  const file = e.target.files?.[0]; if (!file) return
  uploading.value = true
  try {
    const res = await uploadFile(file)
    form.content += `\n![](${res.url})\n`
    toast('图片已上传', 'success')
  } catch (err) { toast('上传失败', 'error') }
  finally { uploading.value = false; e.target.value = '' }
}

// ── 自动保存草稿 ──
const DRAFT_KEY = 'blog_editor_draft'
let draftTimer = null

function scheduleDraft() {
  draftSaved.value = false
  clearTimeout(draftTimer)
  draftTimer = setTimeout(saveDraft, 3000)
}

function saveDraft() {
  if (!form.title && !form.content) return
  localStorage.setItem(DRAFT_KEY, JSON.stringify({
    title: form.title, slug: form.slug, content: form.content,
    summary: form.summary, category_id: form.category_id,
    tag_ids: form.tag_ids, published: form.published,
    savedAt: new Date().toISOString()
  }))
  draftSaved.value = true
}

function clearDraft() {
  localStorage.removeItem(DRAFT_KEY)
  draftSaved.value = false
  clearTimeout(draftTimer)
}

// 恢复草稿
function restoreDraft() {
  try {
    const raw = localStorage.getItem(DRAFT_KEY)
    if (!raw) return
    const d = JSON.parse(raw)
    if (d.title || d.content) {
      Object.assign(form, d)
      formId.value = ''
      editing.value = true
      toast('已恢复未保存的草稿', 'info')
    }
  } catch { /* ignore */ }
}

// ── 分类管理 ──
const renamingCat = ref('')
const renameCatName = ref('')

async function addCategory() {
  if (!newCat.value.trim()) return
  await admin.categories.create({ name: newCat.value.trim() }); newCat.value = ''
  categories.value = await admin.categories.list(); toast('分类已添加', 'success')
}
function startRenameCat(c) { renamingCat.value = c.id; renameCatName.value = c.name }
async function saveRenameCat(id) {
  if (!renameCatName.value.trim()) { renamingCat.value = ''; return }
  await admin.categories.update(id, { name: renameCatName.value.trim() })
  categories.value = await admin.categories.list(); renamingCat.value = ''
  toast('分类已重命名', 'success')
}
async function delCategory(id) {
  if (!confirm('删除分类？文章不会删除，但会失去分类关联。')) return
  await admin.categories.delete(id); categories.value = await admin.categories.list(); toast('已删除', 'info')
}

// ── 标签管理 ──
const renamingTag = ref('')
const renameTagName = ref('')

async function addTag() {
  if (!newTag.value.trim()) return
  await admin.tags.create({ name: newTag.value.trim() }); newTag.value = ''
  tagsList.value = await admin.tags.list(); toast('标签已添加', 'success')
}
function startRenameTag(t) { renamingTag.value = t.id; renameTagName.value = t.name }
async function saveRenameTag(id) {
  if (!renameTagName.value.trim()) { renamingTag.value = ''; return }
  await admin.tags.update(id, { name: renameTagName.value.trim() })
  tagsList.value = await admin.tags.list(); renamingTag.value = ''
  toast('标签已重命名', 'success')
}
async function delTag(id) {
  if (!confirm('删除标签？将移除所有文章的该标签关联。')) return
  await admin.tags.delete(id); tagsList.value = await admin.tags.list(); toast('已删除', 'info')
}

// ── 生命周期 ──
onBeforeUnmount(() => { if (editing.value) saveDraft() })

// ── 初始化 ──
const saved = localStorage.getItem('blog_auth')
if (saved) {
  try {
    const [u, p] = atob(saved).split(':')
    setAuth(u, p)
    admin.test().then(() => { authed.value = true; loadAll() }).catch(() => localStorage.removeItem('blog_auth'))
  } catch { localStorage.removeItem('blog_auth') }
}
</script>

<style scoped>
/* ── 筛选栏 ── */
.filter-bar {
  display: flex; gap: 8px; align-items: center; padding-bottom: 1rem;
  border-bottom: 1px solid var(--border, #e2e8f0); margin-bottom: 1rem;
}
.filter-input { flex: 1; min-width: 0; }
.filter-select { width: 120px; }

/* ── 文章列表项 ── */
.post-info { flex: 1; min-width: 0; }
.post-info strong { display: block; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.post-meta { display: flex; gap: 8px; font-size: 0.8rem; margin-top: 2px; }
.post-actions { display: flex; gap: 4px; flex-shrink: 0; }

.pub { color: #16a34a; font-weight: 600; }
.draft { color: #f59e0b; font-weight: 600; }

/* ── 编辑器 ── */
.editor-container { margin-bottom: 1.5rem; }
.editor-meta { margin-bottom: 0; border-radius: 6px 6px 0 0; }
.meta-row { display: flex; gap: 12px; margin-bottom: 8px; }
.meta-field { flex: 1; }
.meta-field.flex-2 { flex: 2; }
.meta-field label,
.meta-field input,
.meta-field select { display: block; width: 100%; box-sizing: border-box; }
.tag-checkboxes { display: flex; flex-wrap: wrap; gap: 8px; margin-bottom: 8px; }
.tag-label { display: flex; align-items: center; gap: 4px; font-size: 0.85rem; cursor: pointer; }
.tag-label input[type="checkbox"] { width: auto; display: inline; }

.auto-save { font-size: 0.75rem; color: #16a34a; font-weight: normal; margin-left: 12px; }

/* ── 工具栏 ── */
.toolbar {
  display: flex; gap: 4px; align-items: center; padding: 6px 10px;
  background: var(--bg-secondary, #f1f5f9); border: 1px solid var(--border, #e2e8f0);
  border-top: none;
}
.tool-btn {
  width: 28px; height: 28px; padding: 0; font-size: 0.8rem; cursor: pointer;
  border: 1px solid transparent; border-radius: 4px; background: transparent;
  display: flex; align-items: center; justify-content: center;
}
.tool-btn:hover { background: var(--bg, #fff); border-color: var(--border, #cbd5e1); }
.toolbar-spacer { flex: 1; }

/* ── 分栏 ── */
.split-pane { display: flex; border: 1px solid var(--border, #e2e8f0); border-top: none; min-height: 420px; }
.pane { flex: 1; overflow: auto; }
.pane-edit { border-right: 1px solid var(--border, #e2e8f0); }
.pane-edit textarea {
  width: 100%; height: 100%; border: none; resize: none; padding: 12px;
  font-family: 'Fira Code', 'Cascadia Code', monospace; font-size: 0.9rem; line-height: 1.6;
  box-sizing: border-box; outline: none; background: var(--bg, #fff); color: var(--text, #0f172a);
}
.pane-preview { padding: 12px; background: var(--bg, #fafaf9); }
.preview-content :deep(h1) { font-size: 1.5rem; margin: 0.5rem 0; }
.preview-content :deep(h2) { font-size: 1.25rem; margin: 0.5rem 0; }
.preview-content :deep(h3) { font-size: 1.1rem; margin: 0.5rem 0; }
.preview-content :deep(p) { margin: 0.5rem 0; line-height: 1.7; }
.preview-content :deep(pre) { background: #1e293b; color: #e2e8f0; padding: 12px; border-radius: 6px; overflow-x: auto; margin: 0.5rem 0; }
.preview-content :deep(code) { font-family: monospace; font-size: 0.85rem; }
.preview-content :deep(blockquote) { border-left: 3px solid var(--accent, #f97316); padding-left: 12px; margin: 0.5rem 0; color: var(--text-secondary, #64748b); }
.preview-content :deep(ul) { padding-left: 1.5rem; }
.preview-content :deep(img) { max-width: 100%; border-radius: 4px; }
.preview-content :deep(a) { color: var(--accent, #f97316); }
.preview-content :deep(hr) { border: none; border-top: 1px solid var(--border, #e2e8f0); margin: 1rem 0; }

/* ── 编辑器底部 ── */
.editor-footer {
  border-radius: 0 0 6px 6px; margin-top: 0; border-top: none;
  display: flex; align-items: center; gap: 16px;
}
.publish-check { display: flex; align-items: center; gap: 6px; font-size: 0.9rem; cursor: pointer; }
.publish-check input[type="checkbox"] { width: auto; display: inline; }
.form-actions { display: flex; gap: 8px; margin-left: auto; }

/* ── 元管理面板 ── */
.meta-panels { margin-top: 1.5rem; }
.rename-input { width: 120px !important; padding: 2px 6px; margin-right: 4px; }
</style>
