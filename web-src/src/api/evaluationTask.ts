import { request } from '@/utils/request';

// 评教任务信息接口
export interface EvaluationTaskInfo {
  id: number;
  task_name: string; // 任务名称
  academic_year: string; // 学年
  semester: number; // 学期: 0上学期 1下学期
  question_set_name: string; // 问题集名称
  question_set_id: number; // 问题集ID
  start_time: string; // 开始时间
  end_time: string; // 结束时间
  is_deleted: boolean; // 是否删除
  createdAt: string;
  updatedAt: string;
}

// 查询参数接口
export interface QueryParams {
  page?: number;
  limit?: number;
  keyword?: string;
  academic_year?: string;
  semester?: number;
}

// 创建/更新评教任务接口
export interface EvaluationTaskForm {
  id?: number;
  taskName: string;
  academicYear: string;
  semester: number;
  questionSetID: number;
  startTime: string;
  endTime: string;
}

/**
 * 获取评教任务列表
 * @param params 查询参数
 * @returns 评教任务列表数据
 */
export const queryEvaluationTaskList = (params: QueryParams) => {
  return request.get<{
    code: number;
    data: EvaluationTaskInfo[];
    count: number;
    msg: string;
  }>(
    { url: '/manager/api/evaluation/task/list', params },
    { isTransformResponse: false },
  );
};

/**
 * 获取评教任务详情
 * @param id 评教任务ID
 * @returns 评教任务详情
 */
export const getEvaluationTaskDetail = (id: number) => {
  return request.get<{
    code: number;
    data: EvaluationTaskInfo;
    msg: string;
  }>({ url: `/manager/api/evaluation/task/detail/${id}` });
};

/**
 * 添加评教任务
 * @param data 评教任务信息
 * @returns 添加结果
 */
export const addEvaluationTask = (data: EvaluationTaskForm) => {
  return request.post<{
    code: number;
    msg: string;
  }>(
    { url: '/manager/api/evaluation/task/add', data },
    { isTransformResponse: false },
  );
};

/**
 * 更新评教任务
 * @param data 评教任务信息
 * @returns 更新结果
 */
export const updateEvaluationTask = (data: EvaluationTaskForm) => {
  return request.post<{
    code: number;
    msg: string;
  }>(
    { url: '/manager/api/evaluation/task/update', data },
    { isTransformResponse: false },
  );
};

/**
 * 删除评教任务
 * @param id 评教任务ID
 * @returns 删除结果
 */
export const deleteEvaluationTask = (id: number) => {
  return request.post<{
    code: number;
    msg: string;
  }>(
    { url: `/manager/api/evaluation/task/delete/${id}` },
    { isTransformResponse: false },
  );
};

/**
 * 获取所有评教任务（用于下拉选择）
 * @returns 评教任务列表
 */
export const getAllEvaluationTasks = () => {
  return request.get<{
    code: number;
    data: { id: number; task_name: string }[];
    msg: string;
  }>(
    { url: '/manager/api/evaluation/task/all' },
    { isTransformResponse: false },
  );
};

/**
 * 获取当前有效的评教任务
 * @returns 有效的评教任务列表
 */
export const getActiveEvaluationTasks = () => {
  return request.get<{
    code: number;
    data: EvaluationTaskInfo[];
    msg: string;
  }>({ url: '/manager/api/evaluation/task/active' });
};
