import { request } from '@/utils/request';

export interface QueryParmas {
  page: number;
  limit: number;
  [key: string]: any;
}

export interface ITreeNodeData {
  // node-key='id' 设置节点的唯一标识
  id: number | string;
  // 节点显示文本
  label: string;
  // 子节点
  children?: ITreeNodeData[];
  // 链接
  url: string;
  // 组件
  component: string;
  // 图标
  customIcon: string;
  // 类型
  menuType: string;
  // 父节点
  parentId: number;
  // 排序
  order: number;
  // 国际化
  locale: string;
}

export type CreateMenuDto = {
  order: number;
  menuType: string;
  name: string;
  path: string;
  component: string;
  icon: string;
  locale: string;
  parentId: number | null;
};

export function getAllMenu() {
  return request.get<ITreeNodeData[]>({url:'/api/menu'});
}

export function getRoleMenu(role?: string) {
  return request.get({url:`/api/menu/role`});
}

export function updateMenu(data: any) {
  return request.post({url:`/api/menu`, data});
}

export function deleteMenu(id: number, parentId: number) {
  return request.delete({url:`/api/menu?id=${id}&parentId=${parentId}`});
}

export function createMenu(data: any) {
  return request.post({url:`/api/menu`, data});
}


