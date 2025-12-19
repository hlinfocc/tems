import { defineStore } from 'pinia';
import {
  login as userLogin,
  loginMail as userLoginMail,
  getUserInfo,
  updateUserInfo,
  LoginData,
  LoginDataMail,
  checkLogin as userCheckLogin
} from '@/api/user';
import { setToken, clearToken } from '@/utils/auth';
import { removeRouteListener } from '@/utils/route-listener';
import { Modal } from '@opentiny/vue';
import { UserState, UserInfo } from './types';

// 从 localStorage 加载初始状态
const loadInitialState = (): UserState => {
  const savedState = localStorage.getItem('USER_INFO_STORE_STATE');
  return savedState ? JSON.parse(savedState) : {} as any;
};

const useUserStore = defineStore('user', {
  state: (): UserState => loadInitialState(),

  getters: {
    userInfo(state: UserState): UserState {
      return state;
    },
  },

  actions: {
    // 保存状态到 localStorage
    saveState() {
      localStorage.setItem('USER_INFO_STORE_STATE', JSON.stringify(this.$state));
    },

    switchRoles() {
      return new Promise((resolve) => {
        this.role = this.role === 'user' ? 'admin' : 'user';
        this.saveState();
        resolve(this.role);
      });
    },

    // Set user's information
    setInfo(partial: Partial<UserState>) {
      this.$patch(partial);
      this.saveState();
    },

    // Reset user's information
    resetInfo() {
      this.$reset();
      localStorage.removeItem('USER_INFO_STORE_STATE');
    },

    // Reset filter information
    resetFilterInfo() {
      this.startTime = '';
      this.endTime = '';
      this.filterStatus = [];
      this.filterType = [];
      this.saveState();
    },

    // Get user's information
    async info() {
      const res = await getUserInfo();
      console.log("info:", res);
      if (res) {
        this.setInfo(res.data.userInfo);
      }
    },

    async updateInfo(data: UserInfo) {
      return new Promise((resolve) => {
        updateUserInfo(data).then((res)=>{
          if(res.data){
            this.setInfo(res.data.userInfo);
          }
          resolve(res)
        });
      })
    },


      // Login
      async login(loginForm: LoginData) {
        try {
          console.log("LoginData:", loginForm);
          await userLogin(loginForm).then((res) => {
            console.log("res:", res);
            if (res && res.code===200){
              const { token, userInfo } = res.data;
              setToken(token);
              this.setInfo(userInfo);
              Modal.message({
                message: res.msg,
                status: res.code===200?'success':'error',
              });
            }else{
              throw new Error(res && res.msg&&res.msg!==''?res.msg:`登录失败`);
            }
          }).catch((e) => {
            console.log("e>>>>>>>", e.message,e);
            throw e;
          });
        } catch (err) {
          clearToken();
          throw err;
        }
      },
  
    async loginMail(loginForm: LoginDataMail) {
      try {
        const res = await userLoginMail(loginForm);
        setToken(res.data.token);
      } catch (err) {
        clearToken();
        throw err;
      }
    },

    async checkLogin() {
      try {
        const res = await userCheckLogin();
        const hasToken = !!res.data.token; // 转换为布尔值
        if (hasToken) {
          setToken(res.data.token);
        } else {
          clearToken();
        }
        return hasToken;
      } catch (err) {
        clearToken();
        return false;
      }
    },

    // Logout
    async logout() {
      this.resetInfo();
      clearToken();
      removeRouteListener();
    },
  },
});

export default useUserStore;