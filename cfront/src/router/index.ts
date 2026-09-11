import {
  createRouter,
  createWebHashHistory,
} from 'vue-router'
import { useCachedViewStore } from '@/store/modules/cached-view'
import { getAuthToken } from '@/utils/auth-token'
import setPageTitle from '@/utils/set-page-title'
import routes from './routes'

const router = createRouter({
  history: createWebHashHistory(),
  routes,
})

router.beforeEach((to, _from, next) => {
  if (to.meta.requiresAuth && !getAuthToken()) {
    next({ name: 'Auth' })
    return
  }
  if (to.name === 'Auth' && getAuthToken()) {
    next({ name: 'Home' })
    return
  }
  // 路由缓存
  useCachedViewStore().addCachedView(to)
  // 页面 title
  setPageTitle(to.meta.title)
  next()
})

export default router
