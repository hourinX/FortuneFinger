const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
    meta: { title: '登录', public: true },
  },
  {
    path: '/',
    name: 'Layout',
    component: () => import('@/layouts/index.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/profile/index.vue'),
        meta: { title: '个人信息' },
      },
    ],
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/placeholder/index.vue'),
    meta: { title: '页面不存在' },
  },
]

export default routes
