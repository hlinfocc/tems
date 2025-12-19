import localeLogin from '@/views/login/locale/zh-CN';
import localeTheme from '@/components/theme/locale/zh-CN';

import locale403 from '@/views/exception/403/locale/zh-CN';
import locale404 from '@/views/exception/404/locale/zh-CN';
import locale500 from '@/views/exception/500/locale/zh-CN';

import localeUserInfo from '@/views/user/info/locale/zh-CN';
import localeUserSetting from '@/views/user/setting/locale/zh-CN';

import localekanban from '@/views/board/locale/zh-CN';

import localeSettings from './zh-CN/settings';
import localeHttpError from './zh-CN/httpError';

export default {
  'menu.board': '首页',
  'menu.home': '欢迎',
  'menu.work': '工作台',
  'menu.list': '列表页',
  'menu.buildcenter': '建设中心',
  'menu.contentCenter': '内容中心',
  'menu.hosts': '主机管理',
  'menu.result': '结果页',
  'menu.exception': '异常页',
  'menu.notFound': '未知页面',
  'menu.keys': '密钥管理',
  'menu.profile.detail': '基础详情页',
  'menu.visualization': '数据可视化',
  'menu.user': '用户管理',
  'navbar.docs': '文档中心',
  'navbar.action.locale': '切换为中文',
  'messageBox.switchRoles': '切换角色',
  'messageBox.userCenter': '用户中心',
  'messageBox.userSettings': '用户设置',
  'messageBox.logout': '退出登录',
  'menu.cloud': '云服务能力展示',
  'menu.student': '学生管理',
  'menu.student.add': '添加学生',
  'menu.student.edit': '编辑学生',
  'menu.class': '班级管理',
  'menu.course': '课程管理',
  'menu.evaluationTask': '评教任务管理',
  'menu.evaluationQuestionManager': '评教问题管理',
  'menu.evaluationQuestionSet': '评教问题集管理',
  'menu.evaluationQuestionDetail': '评教问题详情管理',
  'menu.user.list': '用户列表',
  'menu.statistics': '统计',
  'menu.banner': '轮播图管理',
  ...localeTheme,
  ...localeSettings,
  ...localeLogin,
  ...locale403,
  ...locale404,
  ...locale500,
  ...localeUserInfo,
  ...localeUserSetting,
  ...localekanban,
  ...localeHttpError,
};
