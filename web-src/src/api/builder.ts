import { request } from '@/utils/request';

export interface QueryParmas {
  page: number;
  limit: number;
  [key: string]: any;
}

export interface TemplatesPre {
  tplTitle: string;
  tplName: string;
  tplContent?: string;
  tplVersion?: number;
  tplStatus?: number;
  [key: string]: any;
}

export function templatespPreList(params?:QueryParmas) {
  return request.post<any>({
    url:'/manager/api/templates/list',
    data:params
  },{isTransformResponse:false});
}

export function templatespPreInsert(params?:TemplatesPre) {
  return request.post<any>({
    url:'/manager/api/templates/insert',
    data:params
  },{isTransformResponse:false});
}

export function templatespPreUpdate(params?:TemplatesPre) {
  return request.post<any>({
    url:'/manager/api/templates/update',
    data:params
  },{isTransformResponse:false});
}
export function templatespPreUpdateCtt(params?:TemplatesPre) {
  return request.put<any>({
    url:'/manager/api/templates/update',
    data:params
  },{isTransformResponse:false});
}
export function templatespPreDelete(id?:number) {
  return request.delete<any>({
    url:'/manager/api/templates/delete',
    params:{
      id
    }
  },{isTransformResponse:true});
}

export function templatespPrePublish(params?:number[]) {
  return request.post<any>({
    url:'/manager/api/templates/publish',
    data:params
  },{isTransformResponse:true});
}