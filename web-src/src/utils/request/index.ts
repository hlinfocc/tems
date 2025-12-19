// axios配置  可自行根据项目进行更改，只需更改该文件即可，其他文件可以不动
import type { AxiosInstance, InternalAxiosRequestConfig } from 'axios';
import isString from 'lodash/isString';
import merge from 'lodash/merge';

import { Notify, Modal } from '@opentiny/vue';
import { getToken, clearToken } from '@/utils/auth';
import { ContentTypeEnum } from '@/utils/constants';
import { VAxios } from './Axios';
import type { AxiosTransform, CreateAxiosOptions } from './AxiosTransform';
import { formatRequestDate, joinTimestamp, setObjToUrlParams } from './utils';
// import {debounce} from '@/utils/request/debounce';

const env = import.meta.env.MODE || 'development';

const host = import.meta.env.VITE_API_URL || '';
const hostPrexy = import.meta.env.VITE_API_URL_PREFIX || '';
// 数据处理，方便区分多种处理方式
const transform: AxiosTransform = {
  // 处理请求数据。如果数据不是预期格式，可直接抛出错误
  transformRequestHook: (res, options) => {
    const { isTransformResponse, isReturnNativeResponse, errorMessageMode } =
      options;

    // 如果204无内容直接返回
    const method = res.config.method?.toLowerCase() || '';
    if (res.status === 204 && ['put', 'patch', 'delete'].includes(method)) {
      return res;
    }

    // 是否返回原生响应头 比如：需要获取响应头时使用该属性
    if (isReturnNativeResponse) {
      return res;
    }
    // 不进行任何处理，直接返回
    // 用于页面代码可能需要直接获取code，data，message这些信息时开启
    if (!isTransformResponse) {
      return res.data;
    }

    // 错误的时候返回
    const { data } = res;
    if (!data) {
      throw new Error('请求接口错误');
    }

    //  这里 code为 后台统一的字段，需要在 types.ts内修改为项目自己的接口返回格式
    // @ts-ignore
    const { code, msg } = data;

    // 这里逻辑可以根据项目进行修改
    switch (code) {
      case 200:
        return data.data;
        break;
      case 401:
      case 402:
        if (errorMessageMode && errorMessageMode === 'message') {
          Modal.message({
            message: msg,
            status: 'error',
          });
        } else if (errorMessageMode && errorMessageMode === 'notify') {
          Notify({
            type: 'error',
            title: '温馨提示',
            message: msg,
            position: 'top-right',
            duration: 5000,
            customClass: 'my-custom-cls',
          });
        } else if (errorMessageMode && errorMessageMode === 'alert') {
          dialogAlert(msg);
        } else if (errorMessageMode && errorMessageMode === 'none') {
          // none
        } else {
          Modal.message({
            message: msg,
            status: 'error',
          });
        }
        return '';
      case 405:
        return '';
      default:
        throw new Error(
          msg && msg !== '' ? msg : `请求接口错误, 错误码: ${code}`,
        );
        break;
    }
  },

  // 请求前处理配置
  beforeRequestHook: (config, options) => {
    const {
      apiUrl,
      isJoinPrefix,
      urlPrefix,
      joinParamsToUrl,
      formatDate,
      joinTime = true,
    } = options;

    // 添加接口前缀
    if (isJoinPrefix && urlPrefix && isString(urlPrefix)) {
      config.url = `${urlPrefix}${config.url}`;
    }

    // 将baseUrl拼接
    if (apiUrl && isString(apiUrl)) {
      config.url = `${apiUrl}${config.url}`;
    }
    const params = config.params || {};
    const data = config.data || false;

    if (formatDate && data && !isString(data)) {
      formatRequestDate(data);
    }
    if (config.method?.toUpperCase() === 'GET') {
      if (!isString(params)) {
        // 给 get 请求加上时间戳参数，避免从缓存中拿数据。
        config.params = Object.assign(
          params || {},
          joinTimestamp(joinTime, false),
        );
      } else {
        // 兼容restful风格
        config.url = `${config.url + params}${joinTimestamp(joinTime, true)}`;
        config.params = undefined;
      }
    } else if (!isString(params)) {
      if (formatDate) {
        formatRequestDate(params);
      }
      if (
        Reflect.has(config, 'data') &&
        config.data &&
        (Object.keys(config.data).length > 0 || data instanceof FormData)
      ) {
        config.data = data;
        config.params = params;
      } else {
        // 非GET请求如果没有提供data，则将params视为data
        config.data = params;
        config.params = undefined;
      }
      if (joinParamsToUrl) {
        config.url = setObjToUrlParams(config.url as string, {
          ...config.params,
          ...config.data,
        });
      }
    } else {
      // 兼容restful风格
      config.url += params;
      config.params = undefined;
    }
    // 处理超时
    if (
      (config as Recordable)?.requestOptions?.timeout &&
      (config as Recordable)?.requestOptions?.timeout > 0
    ) {
      config.timeout = (config as Recordable)?.requestOptions?.timeout;
    }
    return config;
  },

  // 请求拦截器处理
  requestInterceptors: (config, options) => {
    // 请求之前处理config
    const token = getToken();
    if (token && (config as Recordable)?.requestOptions?.withToken !== false) {
      // jwt token
      (config as Recordable).headers.Authorization =
        options.authenticationScheme
          ? `${options.authenticationScheme} ${token}`
          : token;
      (config as Recordable).headers.satoken = options.authenticationScheme
        ? `${options.authenticationScheme} ${token}`
        : token;
    }
    return config as InternalAxiosRequestConfig;
  },

  // 响应拦截器处理
  responseInterceptors: (res) => {
    return res;
  },

  // 响应错误处理
  responseInterceptorsCatch: (error: any, instance: AxiosInstance) => {
    const { config } = error;
    // 这里和重试配置retry.count有关
    if (
      config &&
      config.retryCount &&
      config.retryCount === 1 &&
      !error.response.data.msg
    ) {
      let errMessage = '';
      switch (error.response.status) {
        case 400:
          errMessage = `请求错误`;
          break;
        case 401:
          errMessage = `用户没有权限或登录过期!`;
          break;
        case 403:
          errMessage = `用户得到授权，但是访问是被禁止的!`;
          break;
        case 404:
          errMessage = `网络请求错误,未找到该资源!`;
          break;
        case 405:
          errMessage = `网络请求错误,请求方法未允许!`;
          break;
        case 408:
          errMessage = `网络请求超时!`;
          break;
        case 500:
          errMessage = `服务器错误,请联系管理员!`;
          break;
        case 501:
          errMessage = `网络未实现!`;
          break;
        case 502:
          errMessage = `网络错误或服务异常!`;
          break;
        case 503:
          errMessage = `服务不可用，服务器暂时过载或维护!`;
          break;
        case 504:
          errMessage = `网络超时!`;
          break;
        case 505:
          errMessage = `http版本不支持该请求!`;
          break;
        default:
      }
      if (config.requestOptions.errorMessageMode) {
        switch (config.requestOptions.errorMessageMode) {
          case 'message':
            Modal.message({
              message: errMessage,
              status: 'error',
            });
            break;
          case 'notify':
            Notify({
              type: 'error',
              title: '温馨提示',
              message: errMessage,
              position: 'top-right',
              duration: 5000,
              customClass: 'my-custom-cls',
            });
            break;
          case 'alert':
            dialogAlert(errMessage);
            break;
          default:
        }
      } else {
        Notify({
          type: 'error',
          title: '温馨提示',
          message: errMessage,
          position: 'top-right',
          duration: 5000,
          customClass: 'my-custom-cls',
        });
      }
    }
    if (!config || !config.requestOptions.retry) {
      error.message = error.response.data.msg || error.message;
      return Promise.reject(error);
    }

    config.retryCount = config.retryCount || 0;

    if (config.retryCount >= config.requestOptions.retry.count) {
      error.message = error.response.data.msg || error.message;
      return Promise.reject(error);
    }

    config.retryCount += 1;

    const backoff = new Promise((resolve) => {
      setTimeout(() => {
        resolve(config);
      }, config.requestOptions.retry.delay || 1);
    });
    config.headers = {
      ...config.headers,
      'Content-Type': ContentTypeEnum.Json,
    };

    return backoff.then((configdata) => instance.request(config));
  },
};

function dialogAlert(msg?: any) {
  if (msg) {
    Modal.alert({ message: msg, title: '温馨提示', status: 'warning' });
  }
}

function createAxios(opt?: Partial<CreateAxiosOptions>) {
  return new VAxios(
    merge(
      <CreateAxiosOptions>{
        // https://developer.mozilla.org/en-US/docs/Web/HTTP/Authentication#authentication_schemes
        // 例如: authenticationScheme: 'Bearer'
        authenticationScheme: '',
        // 超时
        timeout: 10 * 1000,
        // 携带Cookie
        withCredentials: true,
        // 头信息
        headers: { 'Content-Type': ContentTypeEnum.Json },
        // 数据处理方式
        transform,
        // 配置项，下面的选项都可以在独立的接口请求中覆盖
        requestOptions: {
          // 接口地址
          apiUrl: host,
          // 是否自动添加接口前缀
          isJoinPrefix: true,
          // 接口前缀
          // 例如: https://www.api.com/api
          // urlPrefix: '/api'
          urlPrefix: hostPrexy,
          // 是否返回原生响应头 比如：需要获取响应头时使用该属性
          isReturnNativeResponse: false,
          // 需要对返回数据进行处理
          isTransformResponse: true,
          // post请求的时候添加参数到url
          joinParamsToUrl: false,
          // 格式化提交参数时间
          formatDate: true,
          // 是否加入时间戳
          joinTime: true,
          // 忽略重复请求
          ignoreRepeatRequest: true,
          // 异常消息提示类型
          errorMessageMode: 'message',
          // 是否携带token
          withToken: true,
          // 超时时间 毫秒
          timeout: 10 * 1000,
          // 重试
          retry: {
            count: 1,
            delay: 1000,
          },
        },
      },
      opt || {},
    ),
  );
}

// export const request = debounce(createAxios(), 500)
export const request = createAxios();
