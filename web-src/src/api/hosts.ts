import axios from 'axios';

export interface QueryParmas {
  page: number;
  limit: number;
  [key: string]: any;
}
export interface HostsType {
	id:number | undefined;
	host:string;
	username:string;
	port:string;
	iskey:number;
	keypath:string;
	hostdesc:string;
}

export function queryHostsList(params: QueryParmas) {
  return axios.post('/api/hosts/list', params);
}
export function deleteHosts(id: any) {
  return axios.delete(`/api/hosts/delete?id=${id}`);
}

export function hostsInsert(data: HostsType) {
  return axios.post('/api/hosts/insert', data);
}
export function hostsUpdate(data: HostsType) {
  return axios.post('/api/hosts/update', data);
}
