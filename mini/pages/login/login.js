Page({
  data: {
    loading: false,
    loadingText: '登录中...'
  },

  // 获取手机号
  getPhoneNumber(e) {
    console.log('getPhoneNumber 事件:', e);
    const { detail } = e;
    
    if (detail.errMsg === 'getPhoneNumber:ok') {
      // 用户同意授权
      this.setData({
        loading: true,
        loadingText: '登录中...'
      });

      // 获取微信登录code
      wx.login({
        success: (res) => {
          if (res.code) {
            // 调用后端登录接口
            this.loginToBackend(res.code, detail);
          } else {
            this.setData({ loading: false });
            wx.showToast({
              title: '获取登录凭证失败',
              icon: 'none',
              duration: 2000
            });
          }
        },
        fail: (err) => {
          this.setData({ loading: false });
          wx.showToast({
            title: '登录失败，请重试',
            icon: 'none',
            duration: 2000
          });
        }
      });
    } else {
      // 用户拒绝授权
      wx.showToast({
        title: '请授权手机号以继续',
        icon: 'none',
        duration: 2000
      });
    }
  },

  // 调用后端登录接口
  loginToBackend(code, phoneDetail) {
    const app = getApp();
    
    wx.request({
      url: `${app.globalData.baseUrl}/v1/students/login`,
      method: 'POST',
      data: {
        code: phoneDetail.code,
        open_id: code
      },
      success: (res) => {
        console.log('登录接口返回:', res);
        
        if (res.statusCode === 200 && res.data) {
          if (res.data.need_bind) {
            // 需要绑定学生信息，跳转到绑定页面
            this.setData({ loading: false });
            // 保存手机号和open_id到全局，供绑定页面使用
            app.globalData.loginInfo = {
              phone: res.data.data.phone,
              open_id: res.data.data.open_id
            };
            wx.navigateTo({
              url: '/pages/bind-student/bind-student'
            });
          } else {
            // 登录成功，已绑定学生信息
            const studentInfo = res.data.data;
            // 保存token
            wx.setStorageSync('token', studentInfo.token);
            // 调用登录成功处理
            this.loginSuccess(studentInfo);
          }
        } else {
          this.handleLoginError(res.data?.message || '登录失败');
        }
      },
      fail: (err) => {
        console.error('登录请求失败:', err);
        this.handleLoginError('网络错误，请检查网络连接');
      }
    });
  },
  
  // 处理登录错误
  handleLoginError(message) {
    this.setData({ loading: false });
    wx.showToast({
      title: message,
      icon: 'none',
      duration: 2000
    });
  },

  // 登录成功
  loginSuccess(studentInfo) {
    this.setData({
      loadingText: '登录成功！'
    });
    
    // 调用全局登录成功方法
    const app = getApp();
    app.loginSuccess(null, studentInfo);
    
    // 延迟后跳转到首页
    setTimeout(() => {
      wx.switchTab({
        url: '/pages/index/index'
      });
    }, 500);
  },

  // 查看隐私政策
  viewPrivacy() {
    wx.showModal({
      title: '隐私政策',
      content: '我们重视您的隐私保护，将严格按照隐私政策处理您的个人信息。',
      showCancel: true,
      cancelText: '关闭',
      confirmText: '了解'
    });
  },

  // 查看用户协议
  viewTerms() {
    wx.showModal({
      title: '用户协议',
      content: '使用本系统前，请仔细阅读并同意用户协议。',
      showCancel: true,
      cancelText: '关闭',
      confirmText: '了解'
    });
  }
});