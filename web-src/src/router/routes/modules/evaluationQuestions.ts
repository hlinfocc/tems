import { RoleType } from '@/types/roleType';

export default {
  path: 'evaluationQuestion',
  name: 'evaluationQuestion',
  id: 'evaluationQuestion',
  label: 'evaluationQuestion',
  component: () => import('@/views/routersView/index.vue'),
  customIcon: 'IconRichTextListUnordered',
  meta: {
    locale: 'menu.evaluationQuestionManager',
    requiresAuth: true,
    order: 5,
    roles: [RoleType.admin, RoleType.user],
  },
  children: [
    {
      path: 'evaluationQuestionSet',
      name: 'EvaluationQuestionSet',
      id: 'EvaluationQuestionSet',
      label: 'EvaluationQuestionSet',
      component: () => import('@/views/evaluationQuestionSet/index.vue'),
      meta: {
        locale: 'menu.evaluationQuestionSet',
        requiresAuth: true,
        order: 1,
        roles: [RoleType.admin, RoleType.user],
      },
    },
    {
      path: 'evaluation-question-detail',
      name: 'EvaluationQuestionDetail',
      id: 'EvaluationQuestionDetail',
      label: 'EvaluationQuestionDetail',
      component: () => import('@/views/evaluation-question-detail/index.vue'),
      meta: {
        locale: 'menu.evaluationQuestionDetail',
        requiresAuth: true,
        order: 2,
        roles: [RoleType.admin, RoleType.user],
      },
    },
  ],
};