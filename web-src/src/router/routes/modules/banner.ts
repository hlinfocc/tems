import { RoleType } from '@/types/roleType';

export default {
  path: 'banner',
  name: 'Banner',
  id: 'Banner',
  label: 'Banner',
  component: () => import('@/views/banner/index.vue'),
  customIcon: 'IconVersiontree',
  meta: {
    locale: 'menu.banner',
    requiresAuth: true,
    order: 7,
    roles: [RoleType.admin],
  },
};
