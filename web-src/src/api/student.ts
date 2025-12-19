import { request } from '@/utils/request';

// 学生信息接口
export interface StudentInfo {
  id?: number;
  studentID: string;
  studentName: string;
  phone: string;
  classID: number;
  className: string;
  password?: string;
  // 以下是响应可能包含的字段
  openID?: string;
  isBound?: boolean;
  createdAt?: string;
  updatedAt?: string;
}

// 查询参数接口
export interface QueryParmas {
  page?: number;
  pageSize?: number;
  studentName?: string;
  studentID?: string;
  classID?: number;
}

// 班级信息接口
export interface ClassInfo {
  id: number;
  name: string;
}

// 响应数据结构
export interface ApiResponse<T = any> {
  code: number;
  msg?: string;
  data: T;
  count?: number;
}

/**
 * 获取学生列表
 * @param params 查询参数
 * @returns 学生列表数据
 */
export function queryStudentList(
  params?: QueryParmas,
): Promise<ApiResponse<any[]>> {
  return request.get({
    url: '/manager/api/students/list',
    params,
  });
}

/**
 * 获取班级列表
 * @returns 班级列表数据
 */
export function queryClassList(): Promise<ApiResponse<ClassInfo[]>> {
  return request.get({
    url: '/student/class/list',
  });
}

/**
 * 添加学生
 * @param data 学生信息
 * @returns 添加结果
 */
export function addStudent(data: StudentInfo): Promise<ApiResponse> {
  return request.post({
    url: '/manager/api/students/add',
    data,
  });
}
export function addBatchStudent(data: StudentInfo[]): Promise<ApiResponse> {
  return request.post({
    url: '/manager/api/students/batch',
    data,
  }, {isTransformResponse: false});
}

/**
 * 更新学生信息
 * @param id 学生ID
 * @param data 更新的学生信息
 * @returns 更新结果
 */
export function updateStudent(
  id: number,
  data: StudentInfo,
): Promise<ApiResponse> {
  return request.post({
    url: `/manager/api/students/update`,
    data: { ...data, id },
  });
}

/**
 * 删除学生
 * @param id 学生ID
 * @returns 删除结果
 */
export function deleteStudent(id: number): Promise<ApiResponse> {
  return request.delete({
    url: `/manager/api/students/delete/${id}`,
  });
}

/**
 * 获取学生详情
 * @param id 学生ID
 * @returns 学生详情信息
 */
export function getStudentDetail(
  id: number,
): Promise<ApiResponse<StudentInfo>> {
  return request.get({
    url: `/manager/api/students/detail/${id}`,
  });
}

/**
 * 批量删除学生
 * @param ids 学生ID数组
 * @returns 删除结果
 */
export function batchDeleteStudent(ids: number[]): Promise<ApiResponse> {
  return request.post({
    url: '/manager/api/students/batch/delete',
    data: { ids },
  });
}
