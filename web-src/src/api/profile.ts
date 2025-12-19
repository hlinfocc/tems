import { request } from '@/utils/request';


export function checkLogin(params?:any) {
  return request.get<any>({
    url:'/manager/api/users/checkLogin',
    params
  },{isTransformResponse:false});
}