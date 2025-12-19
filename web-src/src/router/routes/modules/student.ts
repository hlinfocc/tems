import { RoleType } from '@/types/roleType';

export default {
  path: 'student',
  name: 'Student',
  id: 'Student',
  label: 'Student',
  component: () => import('@/views/student/index.vue'),
  customIcon: 'IconUser',
  meta: {
    locale: 'menu.student',
    requiresAuth: true,
    order: 2,
    roles: [RoleType.admin, RoleType.user],
  },
  // 移除children数组，使用抽屉模式处理添加和编辑操作
};