import axios from 'axios';
import { UserInfo } from '@/store/modules/user/types';
import { request } from '@/utils/request';

export interface LoginData {
  account: string;
  password: string;
  verifycode?: string;
  key?: string;
}

export interface LoginDataMail {
  mailname: string;
  mailpassword: string;
}

export interface LoginRes {
  token: string;
  userInfo: UserInfo;
}
export interface UserRes {
  chartData: [];
  tableData: [];
}
export interface UserData {
  sort?: number | undefined;
  startTime?: string;
  endTime?: string;
  filterStatus?: [];
  filterType?: [];
}

export function login(data: LoginData) {
  return request.post<any>(
    {
      url: '/manager/api/login',
      data,
    },
    { isTransformResponse: false },
  );
}
export function loginMail(data: LoginDataMail) {
  return axios.post<LoginRes>('/manager/api/mail/login', data);
}

export function logout() {
  return request.post<LoginRes>({ url: '/manager/api/user/logout' });
}

export function getUserInfo() {
  return axios.get<LoginRes>(`/manager/api/users/info`);
}

export function updateUserInfo(data: UserInfo) {
  return axios.put<LoginRes>(`/manager/api/users/userInfo`, data);
}

export function getUserData(data?: UserData) {
  return axios.post<UserRes>('/manager/api/users/data', data);
}

export function registerUser(data: LoginData) {
  return axios.post<UserInfo>('/manager/api/users/register', data);
}

export function getVerifycode(params?: any) {
  return request.get<any>(
    {
      url: '/manager/verifycode',
      params,
    },
    { isTransformResponse: false },
  );
}

export function checkLogin(params?: any) {
  return request.get<any>(
    {
      url: '/manager/api/users/checkLogin',
      params,
    },
    { isTransformResponse: false },
  );
}

export interface AdminUserListItem {
  id: number;
  name: string;
  username: string;
  status: number;
  user_type: number;
  createdAt: string;
  updatedAt: string;
}

// 为了兼容其他组件，添加别名
export type UserListItem = AdminUserListItem;

export interface AdminUserListParams {
  page?: number;
  limit?: number;
  isPage?: number;
  keyword?: string;
  account?: string;
  status?: number;
  user_type?: number;
}

export interface AdminUserListResponse {
  code: number;
  msg: string;
  data: {
    list: AdminUserListItem[];
    total: number;
  };
}

export interface AdminUserRequest {
  name: string;
  username: string;
  password?: string;
  status: number;
  user_type: number;
}

export interface ResetPasswordRequest {
  password: string;
}

export interface UpdateStatusRequest {
  status: number;
}

export function getAdminUserList(params?: AdminUserListParams) {
  return request.get<AdminUserListResponse>(
    {
      url: '/manager/api/users/list',
      params,
    },
    { isTransformResponse: false },
  );
}

export function createAdminUser(data: AdminUserRequest) {
  return request.post<any>(
    {
      url: '/manager/api/users/add',
      data,
    },
    { isTransformResponse: false },
  );
}

export function updateAdminUser(id: number, data: AdminUserRequest) {
  return request.post<any>(
    {
      url: '/manager/api/users/update',
      data: { id, ...data },
    },
    { isTransformResponse: false },
  );
}

export function deleteAdminUser(id: number) {
  return request.delete<any>(
    {
      url: `/manager/api/users/delete/${id}`,
    },
    { isTransformResponse: false },
  );
}

export function updateAdminUserStatus(id: number, data: UpdateStatusRequest) {
  return request.post<any>(
    {
      url: `/manager/api/users/${id}/status`,
      data,
    },
    { isTransformResponse: false },
  );
}

export function resetAdminUserPassword(id: number, data: ResetPasswordRequest) {
  return request.post<any>(
    {
      url: `/manager/api/users/${id}/password`,
      data,
    },
    { isTransformResponse: false },
  );
}

// 获取姓名的汉语拼音
export function getPinyin(name: string) {
  return request.post<any>(
    {
      url: '/manager/api/get-pinyin',
      data: { name },
    },
    { isTransformResponse: false },
  );
}

// 为了兼容其他组件，添加getUserList作为getAdminUserList的别名
export const getUserList = getAdminUserList;
