<template>
  <div>
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
      <div class="admin-header">
        <h2 class="section-title">文章管理</h2>
        <button @click="newPost">+ 新建文章</button>
      </div>

      <div v-if="!editing" class="card">
        <div v-for="p in postsList" :key="p.id" class="admin-post-item">
          <span>
            <strong>{{ p.title }}</strong>
            <span :class="p.published ? 'pub' : 'draft'">{{ p.published ? '已发布' : '草稿' }}</span>
          </span>
          <span style="display:flex;gap:4px">
            <button @click="editPost(p)">编辑</button>
            <button @click="togglePublish(p)" class="secondary">{{ p.published ? '取消发布' : '发布' }}</button>
            <button @click="deletePost(p.id)" class="danger">删除</button>
          </span>
        </div>
        <p v-if="postsList.length === 0" class="muted" style="padding:1rem 0">暂无文章。</p>
      </div>

      <div v-else class="card" style="margin-bottom:1.5rem">
        <h3>{{ formId ? '编辑文章' : '新建文章' }}</h3>
        <form @submit.prevent="save" style="margin-top:1rem">
          <label>标题 <input v-model="form.title" required></label>
          <label>Slug <input v-model="form.slug" placeholder="留空自动生成"></label>
          <label>分类 <select v-model="form.category_id">
            <option value="">无分类</option>
            <option v-for="c in categories" :key="c.id" :value="c.id">{{ c.name }}</option>
          </select></label>
          <label>标签</label>
          <div class="tag-checkboxes">
            <label v-for="t in tagsList" :key="t.id" class="tag-label">
              <input type="checkbox" :value="t.id" v-model="form.tag_ids"> {{ t.name }}
            </label>
          </div>
          <label>摘要 <textarea v-model="form.summary" rows="2"></textarea></label>
          <label>
            内容（Markdown）
            <div class="upload-bar">
              <input type="file" accept="image/*" @change="uploadImage" ref="fileInput" style="display:none">
              <button type="button" class="sm secondary" @click="fileInput.click()">📎 上传图片</button>
              <span v-if="uploading" class="muted">上传中...</span>
            </div>
            <textarea v-model="form.content" rows="15" required></textarea>
          </label>
          <label><input type="checkbox" v-model="form.published" style="display:inline;width:auto"> 发布</label>
          <div class="form-actions">
            <button type="submit">保存</button>
            <button type="button" @click="cancelEdit" class="secondary">取消</button>
          </div>
        </form>
      </div>

      <div class="meta-panels">
        <div class="meta-panel">
          <h3>分类</h3>
          <form @submit.prevent="addCategory">
            <input v-model="newCat" placeholder="新分类名称" required>
            <button type="submit" class="sm">添加</button>
          </form>
          <div v-for="c in categories" :key="c.id" class="meta-item">
            <span>{{ c.name }}</span>
            <button @click="delCategory(c.id)" class="danger sm">×</button>
          </div>
        </div>
        <div class="meta-panel">
          <h3>标签</h3>
          <form @submit.prevent="addTag">
            <input v-model="newTag" placeholder="新标签名称" required>
            <button type="submit" class="sm">添加</button>
          </form>
          <div v-for="t in tagsList" :key="t.id" class="meta-item">
            <span>{{ t.name }}</span>
            <button @click="delTag(t.id)" class="danger sm">×</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, inject } from 'vue'
import { setAuth, admin, uploadFile } from '../api'

const toast = inject('toast')

const authed = ref(false)
const authUser = ref('')
const authPass = ref('')
const loginError = ref('')
const uploading = ref(false)
const fileInput = ref(null)

function login() {
  setAuth(authUser.value, authPass.value)
  admin.test().then(() => {
    authed.value = true; loginError.value = ''
    localStorage.setItem('blog_auth', btoa(authUser.value + ':' + authPass.value))
    loadAll()
  }).catch(() => loginError.value = '登录失败')
}

const postsList = ref([])
const categories = ref([])
const tagsList = ref([])
const editing = ref(false)
const formId = ref('')
const newCat = ref('')
const newTag = ref('')
const form = reactive({ title: '', slug: '', content: '', summary: '', category_id: '', tag_ids: [], published: false })

async function loadAll() {
  const [p, c, t] = await Promise.all([admin.posts.list({ size: 50 }), admin.categories.list(), admin.tags.list()])
  postsList.value = p.posts || []; categories.value = c || []; tagsList.value = t || []
}

function newPost() {
  formId.value = ''; editing.value = true
  Object.assign(form, { title: '', slug: '', content: '', summary: '', category_id: '', tag_ids: [], published: false })
}

function editPost(p) {
  formId.value = p.id; editing.value = true
  Object.assign(form, { title: p.title, slug: p.slug, content: p.content, summary: p.summary || '', category_id: p.category_id || '', tag_ids: p.tag_ids || [], published: p.published })
}

function cancelEdit() { editing.value = false }

async function save() {
  const body = { title: form.title, slug: form.slug || undefined, content: form.content, summary: form.summary, category_id: form.category_id, tag_ids: form.tag_ids, published: form.published }
  try {
    if (formId.value) { await admin.posts.update(formId.value, body) } else { await admin.posts.create(body) }
    toast('保存成功', 'success')
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

async function addCategory() {
  await admin.categories.create({ name: newCat.value }); newCat.value = ''
  categories.value = await admin.categories.list(); toast('分类已添加', 'success')
}
async function delCategory(id) {
  await admin.categories.delete(id); categories.value = await admin.categories.list(); toast('已删除', 'info')
}
async function addTag() {
  await admin.tags.create({ name: newTag.value }); newTag.value = ''
  tagsList.value = await admin.tags.list(); toast('标签已添加', 'success')
}
async function delTag(id) {
  await admin.tags.delete(id); tagsList.value = await admin.tags.list(); toast('已删除', 'info')
}

const saved = localStorage.getItem('blog_auth')
if (saved) {
  const [u, p] = atob(saved).split(':')
  setAuth(u, p)
  admin.test().then(() => { authed.value = true; loadAll() }).catch(() => localStorage.removeItem('blog_auth'))
}
</script>

<style scoped>
.upload-bar { margin-top: 0.25rem; margin-bottom: 0.25rem; display: flex; align-items: center; gap: 8px; }
</style>
