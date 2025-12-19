// index.js
Page({
  data: {
    banners: [],
    evaluations: [],
    baseUrl: '',
    isAuthenticated: false
  },

  onLoad() {
    const app = getApp();
    this.baseUrl = app.globalData.baseUrl;
    app.initializeUserStatus();
    this.loadBanners();
    this.loadTaskData();
  },

  onShow() {
    // 页面显示时重新加载数据
    this.loadTaskData();
  },

  // 加载数据
  loadTaskData() {
    const app = getApp();

    // 调用增强版评课列表接口
    const studentInfo = app.globalData.studentInfo;
    console.log('studentInfo:', studentInfo);
    // 检查是否已登录并绑定学生信息
    console.log('studentInfo tttt:', !studentInfo);
    if (!studentInfo || !studentInfo.token) {
      console.log('用户未登录或未绑定学生信息，使用模拟数据');
      // 使用模拟数据作为备用
      this.setData({
        evaluations: []
      });
      return;
    }
    
    // 设置认证状态为true
    this.setData({
      isAuthenticated: app.globalData.isLoggedIn
    });
    
    // 获取token
    const token = wx.getStorageSync('token') || '';
    
    wx.request({
      url: this.baseUrl + '/v1/evaluations/enhanced-tasks',
      method: 'POST',
      header: {
        'content-type': 'application/json',
        'Authorization': `${token}`
      },
      data: {
        //student_id: studentInfo.student_id,
        // class_id: studentInfo.class_id,
        Result: 0,  // 只查询未评的（0表示未评）
        Page: 1,
        Limit: 10
      },
      success: (res) => {
        console.log('获取增强版评课列表成功:', res.data);
        if (res.data.code === 200 && res.data.data && res.data.data.list) {
          // 转换数据格式适配页面展示
          const formattedEvaluations = res.data.data.list.map(item => {
            item.status = '待评教';
            return item;
          });
          
          this.setData({
            evaluations: formattedEvaluations
          });
        }else{
          if(res.data.code === 401){
            // 未登录或token过期，清除本地token
            wx.removeStorageSync('token');
            app.globalData.isLoggedIn = false;
            app.globalData.studentInfo = null;
            // 重新登录
            wx.navigateTo({
              url: '/pages/profile/profile'
            });
            return;
          }
          // 使用模拟数据作为备用
          this.setData({
            evaluations: []
          });
        }
      },
      fail: (err) => {
        console.error('获取增强版评课列表失败', err);
        wx.showToast({
          title: '网络错误，请稍后重试',
          icon: 'none'
        });
        // 使用模拟数据作为备用
        this.setData({
          evaluations: []
        });
      }
    });
  },
  
  // 加载轮播图数据
  loadBanners: function() {
    wx.request({
      url: `${this.baseUrl}/v1/banners`,
      method: 'GET',
      header: {
        'content-type': 'application/json',
        'Authorization': `${wx.getStorageSync('token') || ''}`
      },
      success: (res) => {
        console.log('轮播图接口返回结果:', res);
        if (res.statusCode === 200 && res.data.code === 200 && res.data.data) {
          // 处理返回的轮播图数据
          const banners = res.data.data.map(item => ({
            id: item.id,
            imageUrl: item.image_url || item.imageUrl,
            title: item.title
          }));
          this.setData({ banners });
        } else {
          console.error('轮播图接口调用失败:', res.data.message || '未知错误');
          // 使用模拟轮播图数据作为备用
          this.setData({
            banners: this.getMockBanners()
          });
        }
      },
      fail: (err) => {
        console.error('轮播图网络请求失败:', err);
        // 使用模拟轮播图数据作为备用
        this.setData({
          banners: this.getMockBanners()
        });
      }
    });
  },
  
  // 获取模拟轮播图数据
  getMockBanners: function() {
    return [
      { id: 1, imageUrl: 'http://www.gychinazx.com/uploads/allimg/251110/1-2511101936293H.png', title: '评教系统上线了' },
      { id: 2, imageUrl: 'http://www.gychinazx.com/uploads/allimg/251110/1-251110193615921.png', title: '及时评教，共建校园' }
    ];
  },

  // 跳转到评教页面
  goToEvaluation(e) {
    console.log('跳转到评教页面:',e.currentTarget.dataset);
    const id = e.currentTarget.dataset.id;
    const index = e.currentTarget.dataset.index;
    const evaluations = this.data.evaluations;
    const evaluation = evaluations[index];
    const encodedEvaluation = encodeURIComponent(JSON.stringify(evaluation));
    wx.navigateTo({
      url: `/pages/evaluation/evaluation?id=${id}&evaluation=${encodedEvaluation}`
    });
  }
})