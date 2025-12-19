// profile.js
Page({
  data: {
    userInfo: null,
    studentInfo: null,
    isLoggedIn: false,
    userId: "",
    completedCount: 0,
    pendingCount: 0,
  },

  onLoad() {
    this.loadUserData();
  },

  onShow() {
    this.loadUserData();
    this.loadStats();
  },

  // 加载用户数据
  loadUserData() {
    const app = getApp();

    // 检查登录状态
    if (app.globalData.isLoggedIn && app.globalData.studentInfo) {
      this.setData({
        isLoggedIn: true,
        userInfo: app.globalData.studentInfo,
        studentInfo: app.globalData.studentInfo,
        userId: app.globalData.studentInfo.studentId,
      });
    } else {
      // 未登录状态
      this.setData({
        isLoggedIn: false,
        userInfo: null,
        studentInfo: null,
        userId: "",
        completedCount: 0,
        pendingCount: 0,
      });
    }
  },

  // 加载统计数据
  loadStats() {
    // 实际项目中应该调用API获取统计数据
    const app = getApp();
    const token = wx.getStorageSync("token") || "";
    wx.request({
      url: app.globalData.baseUrl + "/v1/evaluations/statistics",
      header: {
        "content-type": "application/json",
        Authorization: `${token}`,
      },
      success: (res) => {
        if (res.statusCode === 200 && res.data.code === 200) {
          console.log("res.data:", res.data);
          this.setData({
            completedCount: res.data.data.complete,
            pendingCount: res.data.data.incomplete,
          });
        }
      },
    });
  },

  // 登录按钮点击事件
  login() {
    wx.navigateTo({
      url: "/pages/login/login",
    });
  },

  // 退出登录
  logout() {
    wx.showModal({
      title: "退出登录",
      content: "确定要退出登录吗？",
      success: (res) => {
        if (res.confirm) {
          const app = getApp();
          app.logout();
          this.loadUserData();
          wx.showToast({
            title: "已退出登录",
            icon: "success",
          });
        }
      },
    });
  },

  // 查看评教历史
  viewHistory() {
    wx.switchTab({
      url: "/pages/evaluations/evaluations",
    });
  },

  // 查看成绩
  viewScore() {
    wx.navigateTo({
      url: "/pages/score/score",
    });
  },

  // 打开设置
  viewSettings() {
    wx.showToast({
      title: "打开设置",
      icon: "none",
    });
  },
});
