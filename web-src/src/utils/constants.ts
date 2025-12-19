
// 通用请求头
export const ContentTypeEnum = {
    Json: 'application/json;charset=UTF-8',
    FormURLEncoded : 'application/x-www-form-urlencoded;charset=UTF-8',
    FormData : 'multipart/form-data;charset=UTF-8',
}

export const CONTRACT_STATUS = {
    FAIL: 0,
    AUDIT_PENDING: 1,
    EXEC_PENDING: 2,
    EXECUTING: 3,
    FINISH: 4,
  };

export const CONTRACT_STATUS_OPTIONS = [
    { value: CONTRACT_STATUS.FAIL, label: '审核失败' },
    { value: CONTRACT_STATUS.AUDIT_PENDING, label: '待审核' },
    { value: CONTRACT_STATUS.EXEC_PENDING, label: '待履行' },
    { value: CONTRACT_STATUS.EXECUTING, label: '审核成功' },
    { value: CONTRACT_STATUS.FINISH, label: '已完成' },
  ];
  