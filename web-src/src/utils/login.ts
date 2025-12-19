import { request } from '@/utils/request';
import { setToken, clearToken } from '@/utils/auth';

const checkLogin = async() => {
      try {
        const res = await userCheckLogin();
        const hasToken = !!res.token; // 转换为布尔值
        if (hasToken) {
          setToken(res.token);
        } else {
          clearToken();
        }
        return hasToken;
      } catch (err) {
        clearToken();
        return false;
      }
};

function userCheckLogin(params?:any) {
  return request.get<any>({
    url:'/manager/api/users/checkLogin',
    params
  },{isTransformResponse:true});
}

export { checkLogin };
