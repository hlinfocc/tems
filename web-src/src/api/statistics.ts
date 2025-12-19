import { request } from '@/utils/request';

export interface QueryParmas {
  page: number;
  limit: number;
  [key: string]: any;
}

export function homeCount() {
  return request.get({
    url: '/manager/api/statistics/top-card',
  });
}
