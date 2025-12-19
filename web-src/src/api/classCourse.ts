import { request } from '@/utils/request';

/**
 * 班级课程信息接口
 */
export interface ClassCourseInfo {
  id: number;
  classId: number;
  courseId: number;
  teacherId: number;
  className: string;
  courseName: string;
  teacherName: string;
  createdAt: string;
  updatedAt: string;
}

/**
 * 班级课程查询参数
 */
export interface QueryParams {
  page: number;
  limit: number;
  classId?: number;
  courseId?: number;
  teacherId?: number;
}

/**
 * 班级课程表单数据
 */
export interface ClassCourseForm {
  id?: number;
  classId: number;
  courseId: number;
  teacherId: number;
  className?: string;
  courseName?: string;
  teacherName?: string;
}

/**
 * 获取班级课程列表
 * @param params 查询参数
 * @returns 班级课程列表
 */
export const getClassCourseList = (params: QueryParams) => {
  return request.get<{
    code: number;
    data: ClassCourseInfo[];
    count: number;
    msg: string;
  }>(
    { url: '/manager/api/class/course/list', params },
    { isTransformResponse: false },
  );
};

/**
 * 根据ID获取班级课程信息
 * @param id 班级课程ID
 * @returns 班级课程信息
 */
export const getClassCourseByID = (id: number) => {
  return request.get<{
    code: number;
    data: ClassCourseInfo;
    msg: string;
  }>(
    { url: `/manager/api/class/course/${id}` },
    { isTransformResponse: false },
  );
};

/**
 * 添加班级课程
 * @param data 班级课程信息
 * @returns 添加结果
 */
export const addClassCourse = (data: ClassCourseForm) => {
  return request.post<{
    code: number;
    msg: string;
  }>(
    { url: '/manager/api/class/course/add', data },
    { isTransformResponse: false },
  );
};

/**
 * 更新班级课程
 * @param id 班级课程ID
 * @param data 班级课程信息
 * @returns 更新结果
 */
export const updateClassCourse = (id: number, data: ClassCourseForm) => {
  return request.put<{
    code: number;
    msg: string;
  }>(
    { url: `/manager/api/class/course/update/${id}`, data },
    { isTransformResponse: false },
  );
};

/**
 * 删除班级课程
 * @param id 班级课程ID
 * @returns 删除结果
 */
export const deleteClassCourse = (id: number) => {
  return request.delete<{
    code: number;
    msg: string;
  }>(
    { url: `/manager/api/class/course/delete/${id}` },
    { isTransformResponse: false },
  );
};
