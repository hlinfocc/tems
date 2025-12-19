import { request } from '@/utils/request';

// 评教任务详情信息接口
export interface EvaluationTaskDetailInfo {
  id: number;
  taskId: number; // 任务ID
  classId: number; // 班级ID
  courseId: number; // 课程ID
  className: string; // 班级名称
  courseName: string; // 课程名称
  userName: string; // 任课教师名称
  is_deleted: boolean; // 是否删除
  createdAt: string;
  updatedAt: string;
}

// 查询参数接口
export interface QueryParams {
  page?: number;
  limit?: number;
  taskId?: number;
  classId?: number;
}

// 创建/更新评教任务详情接口
export interface EvaluationTaskDetailForm {
  id?: number;
  taskId: number;
  classId: number;
  courseId: number;
  className?: string;
  courseName?: string;
  userName?: string;
}

/**
 * 获取评教任务详情列表
 * @param params 查询参数
 * @returns 评教任务详情列表数据
 */
export const queryEvaluationTaskDetailList = (params: QueryParams) => {
  return request.get<{
    code: number;
    data: EvaluationTaskDetailInfo[];
    count: number;
    msg: string;
  }>(
    { url: '/manager/api/evaluation/task-detail/detail/list', params },
    { isTransformResponse: false },
  );
};

/**
 * 根据任务ID获取评教任务详情列表
 * @param params 查询参数
 * @returns 评教任务详情列表
 */
export const getTaskDetailsByTaskId = (params: QueryParams) => {
  return request.get<{
    code: number;
    data: EvaluationTaskDetailInfo[];
    msg: string;
  }>(
    { url: '/manager/api/evaluation/task-detail/list', params },
    { isTransformResponse: false },
  );
};

/**
 * 获取评教任务详情
 * @param id 评教任务详情ID
 * @returns 评教任务详情
 */
export const getEvaluationTaskDetailById = (id: number) => {
  return request.get<{
    code: number;
    data: EvaluationTaskDetailInfo;
    msg: string;
  }>({ url: `/manager/api/evaluation/task-detail/detail/${id}` });
};

/**
 * 添加评教任务详情
 * @param data 评教任务详情信息
 * @returns 添加结果
 */
export const addEvaluationTaskDetail = (data: EvaluationTaskDetailForm) => {
  return request.post<{
    code: number;
    msg: string;
  }>(
    { url: '/manager/api/evaluation/task-detail/add', data },
    { isTransformResponse: false },
  );
};

/**
 * 批量添加评教任务详情
 * @param data 评教任务详情信息数组
 * @returns 添加结果
 */
export const batchAddEvaluationTaskDetail = (data: any) => {
  return request.post<{
    code: number;
    msg: string;
  }>(
    { url: '/manager/api/evaluation/task-detail/batch', data },
    { isTransformResponse: false },
  );
};

/**
 * 更新评教任务详情
 * @param data 评教任务详情信息
 * @returns 更新结果
 */
export const updateEvaluationTaskDetail = (data: EvaluationTaskDetailForm) => {
  return request.post<{
    code: number;
    msg: string;
  }>(
    { url: '/manager/api/evaluation/task-detail/update', data },
    { isTransformResponse: false },
  );
};

/**
 * 删除评教任务详情
 * @param id 评教任务详情ID
 * @returns 删除结果
 */
export const deleteEvaluationTaskDetail = (id: number) => {
  return request.post<{
    code: number;
    msg: string;
  }>({ url: `/manager/api/evaluation/task-detail/delete/${id}` });
};

/**
 * 根据任务ID批量删除评教任务详情
 * @param taskId 任务ID
 * @returns 删除结果
 */
export const deleteTaskDetailsByTaskId = (taskId: number) => {
  return request.post<{
    code: number;
    msg: string;
  }>({ url: `/manager/api/evaluation/task-detail/deleteByTaskId/${taskId}` });
};
