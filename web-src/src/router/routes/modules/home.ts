import { RoleType } from '@/types/roleType';

export default {
  path: 'board',
  name: 'Board',
  id: 'Board',
  label: 'Board',
  component: () => import('@/views/board/index.vue'),
  customIcon: 'IconPublicHome',
  meta: {
    locale: 'menu.board',
    requiresAuth: true,
    order: 1,
    roles: [RoleType.admin,RoleType.classAdvisor],
  },
  children: [
    {
      path: 'home',
      name: 'Home',
      id: 'Home',
      label: 'Home',
      component: () => import('@/views/board/home/index.vue'),
      customIcon: 'IconTotal',
      meta: {
        locale: 'menu.board.home',
        requiresAuth: true,
        roles: [RoleType.admin,RoleType.classAdvisor],
      },
    },
    {
      path: 'statistics',
      name: 'Statistics',
      id: 'Statistics',
      label: 'Statistics',
      component: () => import('@/views/board/stats/index.vue'),
      customIcon: 'IconStatistics',
      meta: {
        locale: 'menu.statistics',
        requiresAuth: true,
        roles: [RoleType.admin],
      },
    },
  ],
};
