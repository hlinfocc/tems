import axios from 'axios';

export interface QueryParmas {
  page: number;
  limit: number;
  [key: string]: any;
}

export interface generateParmas {
  keyname: string;
  passwd: string;
}
export function queryKeysList(params: QueryParmas) {
  return axios.post('/api/keys/list', params);
}
export function deleteKeys(id: string) {
  return axios.delete(`/api/keys/delete?id=${id}`);
}

export function keysGenerate(data: generateParmas) {
  return axios.post('/api/keys/create', data);
}

