Page({
  data: {
    studentInfo: {
      name: '',
      studentId: '',
      className: ''
    },
    errors: {},
    canSubmit: false,
    submitting: false
  },

  // 监听姓名输入
  onNameInput(e) {
    const { value } = e.detail;
    this.setData({
      'studentInfo.name': value
    });
    this.validateForm();
  },

  // 监听学号输入
  onStudentIdInput(e) {
    const { value } = e.detail;
    this.setData({
      'studentInfo.studentId': value
    });
    this.validateForm();
  },

  // 监听班级输入
  onClassInput(e) {
    const { value } = e.detail;
    this.setData({
      'studentInfo.className': value
    });
    this.validateForm();
  },

  // 表单验证
  validateForm() {
    const { name, studentId, className } = this.data.studentInfo;
    const errors = {};
    
    // 验证姓名
    if (!name || name.trim() === '') {
      errors.name = '请输入姓名';
    }
    
    // 验证学号
    if (!studentId || studentId.trim() === '') {
      errors.studentId = '请输入学号';
    }
    
    // 验证班级
    if (!className || className.trim() === '') {
      errors.className = '请输入班级';
    }
    
    this.setData({
      errors,
      canSubmit: Object.keys(errors).length === 0 && name && studentId && className
    });
  },

  // 提交表单
  submitForm() {
    if (!this.data.canSubmit || this.data.submitting) {
      return;
    }

    this.setData({ submitting: true });

    // 调用后端绑定学生信息接口
    const app = getApp();
    const { phone, open_id } = app.globalData.loginInfo || {};
    
    wx.request({
      url: `${app.globalData.baseUrl}/v1/students/bind`,
      method: 'POST',
      data: {
        studentName: this.data.studentInfo.name,
        studentId: this.data.studentInfo.studentId,
        className: this.data.studentInfo.className,
        phone: phone,
        openId: open_id
      },
      success: (res) => {
        console.log('绑定学生信息接口返回:', res);
        
        if (res.statusCode === 200 && res.data) {
          // 保存token
          wx.setStorageSync('token', res.data.data.token);
          // 绑定成功
          this.bindSuccess(res.data.data);
        } else {
          this.handleBindError(res.data?.message || '绑定失败');
        }
      },
      fail: (err) => {
        console.error('绑定请求失败:', err);
        this.handleBindError('网络错误，请检查网络连接');
      },
      complete: () => {
        this.setData({ submitting: false });
      }
    });
  },
  
  // 处理绑定错误
  handleBindError(message) {
    wx.showToast({
      title: message,
      icon: 'none',
      duration: 2000
    });
  },

  // 绑定成功
  bindSuccess(studentInfo) {
    // 保存学生信息到本地存储
    wx.setStorageSync('studentInfo', studentInfo);
    
    wx.showToast({
      title: '绑定成功！',
      icon: 'success',
      duration: 1500,
      success: () => {
        // 延迟后跳转到首页
        setTimeout(() => {
          // 设置全局登录状态
          const app = getApp();
          app.loginSuccess(null, studentInfo);
          
          wx.switchTab({
            url: '/pages/index/index'
          });
        }, 1500);
      }
    });
  },

  // 返回上一页
  onBackPress() {
    if (this.data.submitting) {
      return true; // 阻止返回
    }
    return false;
  },

  // 页面加载
  onLoad() {
    // 检查是否已有学生信息（用于测试）
    const savedInfo = wx.getStorageSync('studentInfo');
    if (savedInfo) {
      this.setData({
        studentInfo: savedInfo
      });
      this.validateForm();
    }
  }
});