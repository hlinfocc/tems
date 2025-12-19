import { request } from '@/utils/request';

// 轮播图信息接口
export interface BannerInfo {
  id: number;
  title: string; // 轮播图标题
  image_url: string; // 图片URL
  is_visible: boolean; // 是否可见
  sort: number; // 排序值
  createdAt: string;
  updatedAt: string;
}

// 查询参数接口
export interface QueryParams {
  page?: number;
  limit?: number;
  keyword?: string;
}

// 创建/更新轮播图接口
export interface BannerForm {
  id?: number;
  title: string;
  image_url: string;
  is_visible: boolean;
  sort: number;
}

/**
 * 获取轮播图列表
 * @param params 查询参数
 * @returns 轮播图列表数据
 */
export const queryBannerList = (params: QueryParams) => {
  return request.get<{
    code: number;
    data: BannerInfo[];
    count: number;
    msg: string;
  }>(
    { url: '/manager/api/banners/list', params },
    { isTransformResponse: true },
  );
};

/**
 * 获取轮播图详情
 * @param id 轮播图ID
 * @returns 轮播图详情
 */
export const getBannerDetail = (id: number) => {
  return request.get<{
    code: number;
    data: BannerInfo;
    msg: string;
  }>(
    { url: `/manager/api/banners/detail/${id}` },
    { isTransformResponse: false },
  );
};

/**
 * 添加轮播图
 * @param data 轮播图信息
 * @returns 添加结果
 */
export const addBanner = (data: BannerForm) => {
  return request.post<{
    code: number;
    msg: string;
  }>({ url: '/manager/api/banners/add', data }, { isTransformResponse: false });
};

/**
 * 更新轮播图
 * @param data 轮播图信息
 * @returns 更新结果
 */
export const updateBanner = (data: BannerForm) => {
  return request.post<{
    code: number;
    msg: string;
  }>(
    { url: '/manager/api/banners/update', data },
    { isTransformResponse: false },
  );
};

/**
 * 删除轮播图
 * @param id 轮播图ID
 * @returns 删除结果
 */
export const deleteBanner = (id: number) => {
  return request.delete<{
    code: number;
    msg: string;
  }>(
    { url: `/manager/api/banners/delete/${id}` },
    { isTransformResponse: false },
  );
};

/**
 * 更新轮播图可见性
 * @param data 可见性信息
 * @returns 更新结果
 */
export const updateBannerVisibility = (data: {
  id: number;
  is_visible: boolean;
}) => {
  return request.post<{
    code: number;
    msg: string;
  }>(
    { url: '/manager/api/banners/visibility', data },
    { isTransformResponse: false },
  );
};
