import { request } from '@/utils/request';

// 评教问题集信息接口
export interface EvaluationQuestionSetInfo {
  id: number;
  name: string; // 问题集名称
  remark: string; // 备注
  status: number; // 状态：0-正常，1-禁用
  createdAt: string;
  updatedAt: string;
}

// 查询参数接口
export interface QueryParmas {
  page?: number;
  limit?: number;
  keyword?: string;
  status?: string;
}

// 创建/更新评教问题集接口
export interface EvaluationQuestionSetForm {
  id?: number;
  name: string;
  remark: string;
  status: number;
}

/**
 * 获取评教问题集列表
 * @param params 查询参数
 * @returns 评教问题集列表数据
 */
export const queryEvaluationQuestionSetList = (params: QueryParmas) => {
  return request.get<{
    code: number;
    data: EvaluationQuestionSetInfo[];
    count: number;
    msg: string;
  }>(
    { url: '/manager/api/evaluation/question-set/list', params },
    {
      isTransformResponse: false,
    },
  );
};

/**
 * 获取评教问题集详情
 * @param id 问题集ID
 * @returns 问题集详情
 */
export const getEvaluationQuestionSetDetail = (id: number) => {
  return request.get<{
    code: number;
    data: EvaluationQuestionSetInfo;
    msg: string;
  }>({ url: `/manager/api/evaluation/question-set/detail/${id}` });
};

/**
 * 添加评教问题集
 * @param data 问题集信息
 * @returns 添加结果
 */
export const addEvaluationQuestionSet = (data: EvaluationQuestionSetForm) => {
  return request.post<{
    code: number;
    msg: string;
  }>(
    { url: '/manager/api/evaluation/question-set/add', data },
    { isTransformResponse: false },
  );
};

/**
 * 更新评教问题集
 * @param data 问题集信息
 * @returns 更新结果
 */
export const updateEvaluationQuestionSet = (
  data: EvaluationQuestionSetForm,
) => {
  return request.post<{
    code: number;
    msg: string;
  }>(
    { url: '/manager/api/evaluation/question-set/update', data },
    { isTransformResponse: false },
  );
};

/**
 * 删除评教问题集
 * @param id 问题集ID
 * @returns 删除结果
 */
export const deleteEvaluationQuestionSet = (id: number) => {
  return request.post<{
    code: number;
    msg: string;
  }>(
    { url: `/manager/api/evaluation/question-set/delete/${id}` },
    { isTransformResponse: false },
  );
};

/**
 * 获取所有评教问题集（用于下拉选择）
 * @returns 问题集列表
 */
export const getAllEvaluationQuestionSets = () => {
  return request.get<{
    code: number;
    data: { id: number; name: string }[];
    msg: string;
  }>({ url: '/manager/api/evaluation/question-set/all' });
};
