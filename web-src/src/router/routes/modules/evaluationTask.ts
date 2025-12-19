import { RoleType } from '@/types/roleType';

export default {
  path: 'evaluation-task',
  name: 'EvaluationTask',
  id: 'EvaluationTask',
  label: 'EvaluationTask',
  component: () => import('@/views/evaluationTask/index.vue'),
  customIcon: 'IconStreamSolid',
  meta: {
    locale: 'menu.evaluationTask',
    requiresAuth: true,
    order: 6,
    roles: [RoleType.admin, RoleType.user],
  },
  // 使用抽屉模式处理添加和编辑操作
};
