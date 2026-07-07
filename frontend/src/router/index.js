import { createRouter, createWebHistory } from 'vue-router'
import Home from '../views/Home.vue'
import Post from '../views/Post.vue'
import Admin from '../views/Admin.vue'
import Archive from '../views/Archive.vue'
import Category from '../views/Category.vue'

const routes = [
  { path: '/', name: 'Home', component: Home },
  { path: '/post/:slug', name: 'Post', component: Post },
  { path: '/archive', name: 'Archive', component: Archive },
  { path: '/category/:slug', name: 'Category', component: Category },
  { path: '/admin', name: 'Admin', component: Admin },
]

export default createRouter({
  history: createWebHistory(),
  routes,
})
