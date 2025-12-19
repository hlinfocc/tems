import { RoleType } from '@/types/roleType';

export default {
  path: 'user',
  name: 'User',
  id: 'User',
  label: 'User',
  component: () => import('@/views/user/index.vue'),
  customIcon:'IconSetting',
  meta: {
    locale: 'menu.user',
    requiresAuth: true,
    order: 10,
    roles: [RoleType.admin, RoleType.user],
  },
  children: [
    {
      path: 'info',
      name: 'Info',
      id: 'Info',
      label: 'Info',
      component: () => import('@/views/user/list/index.vue'),
      meta: {
        locale: 'menu.user.list',
        requiresAuth: true,
        roles: [RoleType.admin, RoleType.user],
      },
    },
    // {
    //   path: 'setting',
    //   name: 'Setting',
    //   id: 'Setting',
    //   label: 'Setting',
    //   component: () => import('@/views/user/setting/index.vue'),
    //   meta: {
    //     locale: 'menu.user.setting',
    //     requiresAuth: true,
    //     roles: [RoleType.admin, RoleType.user],
    //   },
    // },
  ],
};
