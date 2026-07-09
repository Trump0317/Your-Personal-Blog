import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  { path: '/',              name: 'Home',       component: () => import('../views/Home.vue') },
  { path: '/post/:slug',    name: 'Post',       component: () => import('../views/Post.vue') },
  { path: '/archive',       name: 'Archive',    component: () => import('../views/Archive.vue') },
  { path: '/categories',    name: 'Categories', component: () => import('../views/Categories.vue') },
  { path: '/categories/:slug', name: 'Category', component: () => import('../views/Category.vue') },
  { path: '/tags/:slug',    name: 'Tag',        component: () => import('../views/Tag.vue') },
  { path: '/about',         name: 'About',      component: () => import('../views/About.vue') },
  { path: '/admin',         name: 'Admin',      component: () => import('../views/Admin.vue') },
]

export default createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() { return { top: 0 } },
})
