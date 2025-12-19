import { request } from '@/utils/request';

// 班级信息接口
export interface ClassInfo {
  id: number;
  className: string; // 班级名称
  majorName: string; // 专业名称
  majorID: number; // 专业ID
  grade: string; // 年级
  headTeacherId: number; // 班主任ID
  headTeacherName: string; // 班主任姓名
  isDeleted: boolean;
  createdAt: string;
  updatedAt: string;
}

// 查询参数接口
export interface QueryParmas {
  page?: number;
  limit?: number;
  keyword?: string;
  grade?: string;
  majorId?: number;
}

// 创建/更新班级接口
export interface ClassForm {
  id?: number;
  className: string;
  majorName: string;
  majorID: number;
  grade: string;
  headTeacherId: number;
  headTeacherName: string;
}

/**
 * 获取班级列表
 * @param params 查询参数
 * @returns 班级列表数据
 */
export const queryClassList = (params: QueryParmas) => {
  return request.get<{
    code: number;
    data: ClassInfo[];
    count: number;
    msg: string;
  }>(
    { url: '/manager/api/classes/list', params },
    { isTransformResponse: false },
  );
};

/**
 * 获取班级详情
 * @param id 班级ID
 * @returns 班级详情
 */
export const getClassDetail = (id: number) => {
  return request.get<{
    code: number;
    data: ClassInfo;
    msg: string;
  }>(
    { url: `/manager/api/classes/detail/${id}` },
    { isTransformResponse: false },
  );
};

/**
 * 添加班级
 * @param data 班级信息
 * @returns 添加结果
 */
export const addClass = (data: ClassForm) => {
  return request.post<{
    code: number;
    msg: string;
  }>({ url: '/manager/api/classes/add', data }, { isTransformResponse: false });
};

/**
 * 更新班级
 * @param data 班级信息
 * @returns 更新结果
 */
export const updateClass = (data: ClassForm) => {
  return request.post<{
    code: number;
    msg: string;
  }>(
    { url: '/manager/api/classes/update', data },
    { isTransformResponse: false },
  );
};

/**
 * 删除班级
 * @param id 班级ID
 * @returns 删除结果
 */
export const deleteClass = (id: number) => {
  return request.post<{
    code: number;
    msg: string;
  }>(
    { url: `/manager/api/classes/delete/${id}` },
    { isTransformResponse: false },
  );
};
