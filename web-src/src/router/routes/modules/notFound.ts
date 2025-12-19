import { RoleType } from '@/types/roleType';

export default {
  path: 'notFound',
  name: 'notFound',
  id: 'notFound',
  label: 'notFound',
  component: () => import('@/views/not-found/index.vue'),
  meta: {
    locale: 'menu.notFound',
    requiresAuth: false,
    order: 1,
    hideInMenu:true
  }
};
