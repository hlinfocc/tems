// app.js
App({
  globalData: {
    userInfo: null,
    studentInfo: null,
    isLoggedIn: false,
    baseUrl: 'https://pj.gychinazx.com/api'
  },

  onLaunch() {
    // 展示本地存储能力
    const logs = wx.getStorageSync('logs') || []
    logs.unshift(Date.now())
    wx.setStorageSync('logs', logs)

    // 初始化用户登录状态
    this.initializeUserStatus()
  },
  onShow (options) {
    // 初始化用户登录状态
    this.initializeUserStatus()
  },

  // 初始化用户登录状态
  initializeUserStatus() {
    // 检查是否已登录（这里通过检查是否有学生信息来模拟）
    const studentInfo = wx.getStorageSync('studentInfo')
    if (studentInfo) {
      this.globalData.studentInfo = studentInfo
      this.globalData.isLoggedIn = true
      console.log('用户已登录并绑定学生信息')
    }

    // 获取用户信息
    wx.getSetting({
      success: res => {
        if (res.authSetting['scope.userInfo']) {
          wx.getUserInfo({
            success: res => {
              this.globalData.userInfo = res.userInfo
            }
          })
        }
      }
    })
  },

  // 用户登录成功后调用
  loginSuccess(userInfo, studentInfo) {
    this.globalData.userInfo = userInfo || this.globalData.userInfo
    this.globalData.studentInfo = studentInfo
    this.globalData.isLoggedIn = true
    
    // 保存学生信息到本地存储
    wx.setStorageSync('studentInfo', studentInfo)
    console.log('用户登录成功')
  },

  // 用户退出登录
  logout() {
    this.globalData.userInfo = null
    this.globalData.studentInfo = null
    this.globalData.isLoggedIn = false
    
    // 清除本地存储的学生信息
    wx.removeStorageSync('studentInfo')
    console.log('用户退出登录')
  }
})