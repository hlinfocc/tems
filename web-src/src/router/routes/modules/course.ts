import { RoleType } from '@/types/roleType';

export default {
  path: 'course',
  name: 'Course',
  id: 'Course',
  label: 'Course',
  component: () => import('@/views/course/index.vue'),
  customIcon: 'IconCourse',
  meta: {
    locale: 'menu.course',
    requiresAuth: true,
    order: 3,
    roles: [RoleType.admin],
  },
};
