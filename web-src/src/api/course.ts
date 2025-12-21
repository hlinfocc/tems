import { request } from '@/utils/request';

// 课程信息接口
export interface CourseInfo {
  id: number;
  courseName: string; // 课程名称
  status: number; // 状态：0-正常，1-禁用
  createdAt: string;
  updatedAt: string;
}

// 查询参数接口
export interface QueryParmas {
  page?: number;
  limit?: number;
  courseName?: string;
  status?: string;
}

// 创建/更新课程接口
export interface CourseForm {
  id?: number;
  courseName: string;
  status: number;
}

/**
 * 获取课程列表
 * @param params 查询参数
 * @returns 课程列表数据
 */
export const queryCourseList = (params: QueryParmas) => {
  return request.get<{
    code: number;
    data: CourseInfo[];
    count: number;
    msg: string;
  }>(
    { url: '/manager/api/courses/list', params },
    { isTransformResponse: false },
  );
};

/**
 * 获取课程详情
 * @param id 课程ID
 * @returns 课程详情
 */
export const getCourseDetail = (id: number) => {
  return request.get<{
    code: number;
    data: CourseInfo;
    msg: string;
  }>(
    { url: `/manager/api/courses/detail/${id}` },
    { isTransformResponse: false },
  );
};

/**
 * 添加课程
 * @param data 课程信息
 * @returns 添加结果
 */
export const addCourse = (data: CourseForm) => {
  return request.post<{
    code: number;
    msg: string;
  }>({ url: '/manager/api/courses/add', data }, { isTransformResponse: false });
};

export const batchAddCourse = (data: CourseForm[]) => {
  return request.post<{
    code: number;
    msg: string;
  }>({ url: '/manager/api/courses/batchAdd', data }, { isTransformResponse: false });
};

/**
 * 更新课程
 * @param data 课程信息
 * @returns 更新结果
 */
export const updateCourse = (data: CourseForm) => {
  return request.post<{
    code: number;
    msg: string;
  }>(
    { url: '/manager/api/courses/update', data },
    { isTransformResponse: false },
  );
};

/**
 * 删除课程
 * @param id 课程ID
 * @returns 删除结果
 */
export const deleteCourse = (id: number) => {
  return request.post<{
    code: number;
    msg: string;
  }>(
    { url: `/manager/api/courses/delete/${id}` },
    { isTransformResponse: false },
  );
};

/**
 * 获取所有课程（用于下拉选择）
 * @returns 课程列表
 */
export const getAllCourses = () => {
  return request.get<{
    code: number;
    data: { id: number; course_name: string; course_code: string }[];
    msg: string;
  }>({ url: '/manager/api/courses/all' }, { isTransformResponse: false });
};
