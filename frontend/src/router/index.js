import { createRouter, createWebHistory } from 'vue-router'
import UserListView from '../views/UserListView.vue'
import LoginView from '../views/LoginView.vue'
import { useUserAccountStore } from '../stores/userAccount'
import RegisterView from '../views/RegisterView.vue'
import MainLayout from '../layouts/MainLayout.vue'
import ChatView from '../views/ChatView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      component: MainLayout,
      meta: { requiresAuth: true },  
      children: [
        {
          path: '',
          name: 'user-list',
          component: UserListView,
        },
        {
          path: '/chat',
          name: 'chat',
          component: ChatView,
        },
      ]
    },
    {
      path: '/login',
      name: 'login',
      component: LoginView,
    },
    {
      path: '/register',
      name: 'register',
      component: RegisterView,
    },
  ],
}) 

router.beforeEach(async (to, from, next) => {
  const userStore = useUserAccountStore()

  if (to.meta.requiresAuth) {
    await userStore.checkAuth()

    if (!userStore.isAuthenticated || !userStore.user === null ) {
      return next({ name: 'login' })
    }
  }
  next()
})
export default router
