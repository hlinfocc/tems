import localeLogin from '@/views/login/locale/en-US';
import localeTheme from '@/components/theme/locale/en-US';

import locale403 from '@/views/exception/403/locale/en-US';
import locale404 from '@/views/exception/404/locale/en-US';
import locale500 from '@/views/exception/500/locale/en-US';

import localeUserInfo from '@/views/user/info/locale/en-US';
import localeUserSetting from '@/views/user/setting/locale/en-US';

import localekanban from '@/views/board/locale/en-US';

import localeSettings from './en-US/settings';

import localeHttpError from './en-US/httpError';

export default {
  'menu.board': 'Dashboard Page',
  'menu.home': 'Monitoring page',
  'menu.work': 'workbench',
  'menu.list': 'List',
  'menu.buildcenter': 'Builder Center',
  'menu.hosts': 'Host Manager',
  'menu.result': 'Result',
  'menu.exception': 'Exception',
  'menu.form': 'Form',
  'menu.keys': 'keys for ssh',
  'menu.profile.detail': 'Basic details page',
  'menu.visualization': 'Data Visualization',
  'menu.user': 'User Management',
  'navbar.docs': 'Docs',
  'navbar.action.locale': 'Switch to English',
  'messageBox.switchRoles': 'Switch Roles',
  'messageBox.userCenter': 'User Center',
  'messageBox.userSettings': 'User Settings',
  'messageBox.logout': 'Logout',
  'menu.cloud': 'Cloud service capability',
  'menu.student': 'Student Management',
  'menu.student.add': 'Add Student',
  'menu.student.edit': 'Edit Student',
  'menu.class': 'Class Management',
  'menu.course': 'Course Management',
  'menu.evaluationTask': 'Evaluation Task Management',
  'menu.evaluationQuestionManager': 'Evaluation Question Management',
  'menu.evaluationQuestionSet': 'Evaluation Question Set Management',
  'menu.evaluationQuestionDetail': 'Evaluation Question Detail Management',
  'menu.user.list': 'User List',
  'menu.statistics': 'Statistics',
  'menu.banner': 'Banner Management',
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
