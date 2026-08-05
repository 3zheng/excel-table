import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: LoginView,
    },
    {
      path: '/',
      component: HomeView,
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          redirect: '/export',
        },
        {
          path: 'export',
          name: 'export',
          component: () => import('../views/ExportListView.vue'),
        },
        {
          path: 'import',
          name: 'import',
          component: () => import('../views/ImportListView.vue'),
        },
        {
          path: 'products',
          name: 'products',
          component: () => import('../views/ProductsView.vue'),
        },
        {
          path: 'users',
          name: 'users',
          meta: { requiresAdmin: true },
          component: () => import('../views/UsersView.vue'),
        },
        {
          path: 'export/:inv_no',
          name: 'export-detail',
          component: () => import('../views/ExportDetailView.vue'),
        },
        {
          path: 'import/:inv_no',
          name: 'import-detail',
          component: () => import('../views/ImportDetailView.vue'),
        },
      ],
    },
  ],
})

router.beforeEach((to, from) => {
  const token = localStorage.getItem('token')
  if (to.meta.requiresAuth && !token) {
    return '/login'
  } else if (to.path === '/login' && token) {
    return '/'
  }
  if (to.meta.requiresAdmin) {
    try {
      const user = JSON.parse(localStorage.getItem('user') || '{}')
      if (user.user_role !== 'admin') {
        return '/export'
      }
    } catch {
      return '/export'
    }
  }
})

export default router