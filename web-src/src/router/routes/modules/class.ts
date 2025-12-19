import { RoleType } from '@/types/roleType';

export default {
  path: 'class',
  name: 'Class',
  id: 'Class',
  label: 'Class',
  component: () => import('@/views/class/index.vue'),
  customIcon: 'IconVersiontree',
  meta: {
    locale: 'menu.class',
    requiresAuth: true,
    order: 4,
    roles: [RoleType.admin, RoleType.user],
  },
  // 使用抽屉模式处理添加和编辑操作
};
