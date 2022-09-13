import { createWebHistory, createRouter } from 'vue-router'
import { ROUTE_NAMES } from '@/constants'
import { useAuthStore } from './store/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', redirect: '/home' },
    { path: '/home', name: ROUTE_NAMES.HOME, component: () => import('@/views/Home.vue') },
    {
      path: '/authentication',
      name: ROUTE_NAMES.AUTHENTICATION,
      component: () => import('@/views/Authentication.vue'),
      meta: { disableLayout: true }
    },
    {
      path: '/playlist/:id',
      name: ROUTE_NAMES.PLAYLIST,
      component: () => import('@/views/Playlist.vue')
    },
    {
      path: '/queue',
      name: ROUTE_NAMES.QUEUE,
      component: () => import('@/views/Queue.vue')
    }
  ]
})

router.beforeEach((to, from, next) => {
  if (!useAuthStore().isConnected && to.name !== ROUTE_NAMES.AUTHENTICATION) {
    console.log("redirect to authentication")
    next({ name: ROUTE_NAMES.AUTHENTICATION })
  } else {
    next()
  }
})

export default router
