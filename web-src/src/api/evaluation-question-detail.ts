import { request } from '@/utils/request';

// 评教问题详情接口
export interface QuestionDetailInfo {
  id: number;
  questionSetId: number; // 问题集ID
  title: string; // 问题标题
  questionType: number; // 问题类型: 单选、多选
  options: string; // 选项，JSON格式
  remark: string; // 备注
  createdAt: string;
  updatedAt: string;
}

// 查询参数接口
export interface QueryParams {
  page?: number;
  limit?: number;
  keyword?: string;
  questionSetId?: string;
}

// 创建/更新评教问题详情接口
export interface QuestionDetailForm {
  id?: number;
  questionSetId: number;
  title: string;
  questionType: number;
  options: string;
  remark: string;
}

/**
 * 获取评教问题详情列表
 * @param params 查询参数
 * @returns 评教问题详情列表数据
 */
export const queryQuestionDetailList = (params: QueryParams) => {
  return request.get<{
    code: number;
    data: QuestionDetailInfo[];
    count: number;
    msg: string;
  }>(
    { url: '/manager/api/evaluation/question-detail/list', params },
    { isTransformResponse: false },
  );
};

/**
 * 获取评教问题详情
 * @param id 评教问题详情ID
 * @returns 评教问题详情
 */
export const getQuestionDetail = (id: number) => {
  return request.get<{
    code: number;
    data: QuestionDetailInfo;
    msg: string;
  }>({ url: `/manager/api/evaluation/question-detail/detail/${id}` });
};

/**
 * 根据问题集ID获取问题列表
 * @param setId 问题集ID
 * @returns 问题列表
 */
export const getQuestionsBySetId = (setId: number) => {
  return request.get<{
    code: number;
    data: QuestionDetailInfo[];
    msg: string;
  }>({ url: `/manager/api/evaluation/question-detail/by-set/${setId}` });
};

/**
 * 添加评教问题详情
 * @param data 评教问题详情信息
 * @returns 添加结果
 */
export const addQuestionDetail = (data: QuestionDetailForm) => {
  return request.post<{
    code: number;
    msg: string;
  }>(
    { url: '/manager/api/evaluation/question-detail/add', data },
    { isTransformResponse: false },
  );
};

/**
 * 更新评教问题详情
 * @param data 评教问题详情信息
 * @returns 更新结果
 */
export const updateQuestionDetail = (data: QuestionDetailForm) => {
  return request.post<{
    code: number;
    msg: string;
  }>(
    { url: '/manager/api/evaluation/question-detail/update', data },
    { isTransformResponse: false },
  );
};

/**
 * 删除评教问题详情
 * @param id 评教问题详情ID
 * @returns 删除结果
 */
export const deleteQuestionDetail = (id: number) => {
  return request.post<{
    code: number;
    msg: string;
  }>(
    { url: `/manager/api/evaluation/question-detail/delete/${id}` },
    { isTransformResponse: false },
  );
};
