// evaluations.js
Page({
  data: {
    evaluations: [],
    filteredEvaluations: [],
    filterTabs: [
      { id: 'all', name: '全部' },
      { id: 'pending', name: '待评教' },
      { id: 'completed', name: '已完成' }
    ],
    currentFilter: 'all',
    searchText: '',
    isAuthenticated: false
  },

  onLoad() {
    this.loadData();
  },

  onShow() {
    // 页面显示时重新加载数据，确保状态更新
    this.loadData();
  },

  // 加载数据
  loadData() {
    const app = getApp();
    const studentInfo = app.globalData.studentInfo;
    
    // 检查是否已登录并绑定学生信息
    if (!studentInfo || !studentInfo.student_id) {
      console.log('用户未登录或未绑定学生信息，使用模拟数据');
      // 使用模拟数据作为备用
      const allEvaluations = [];
      
      this.setData({
        evaluations: allEvaluations,
        filteredEvaluations: allEvaluations
      });
      this.applyFilters();
      // wx.showModal({
      //   title: '温馨提示',
      //   content: '您需要先登录才能查看',
      //   showCancel: false,
      //   success (res) {
      //     if (res.confirm) {
      //       wx.switchTab({
      //         url: '/pages/profile/profile'
      //       })
      //     }
      //   }
      // })
      return;
    }
    // 设置认证状态为true
    this.setData({
      isAuthenticated: app.globalData.isLoggedIn
    });
    // 获取token
    const token = wx.getStorageSync('token') || '';
    
    // 根据当前筛选条件决定是否只查询未评
    let resultParam = -1; // 默认不筛选
    if (this.data.currentFilter === 'pending') {
      resultParam = 0; // 只查询未评的
    }
    
    wx.request({
      url: app.globalData.baseUrl + '/v1/evaluations/enhanced-tasks',
      method: 'POST',
      header: {
        'content-type': 'application/json',
        'Authorization': `${token}`
      },
      data: {
        student_id: studentInfo.student_id,
        class_id: studentInfo.class_id,
        Result: resultParam, // 根据筛选条件设置
        Page: 1,
        Limit: 50 // 获取更多数据以便在前端进行筛选
      },
      success: (res) => {
        console.log('获取增强版评课列表成功:', res.data);
        if (res.data.code === 200 && res.data.data && res.data.data.list) {
          // 转换数据格式适配页面展示
          const formattedEvaluations = res.data.data.list.map(item => {
            // 根据result值判断状态
            let status = '待评教';
            if (item.result === 1 || item.result === 2) {
              status = '已完成';
            }
            
            return {
              id: item.taskId || item.id,
              courseName: item.courseName,
              teacherName: item.userName,
              deadline: item.endTime,
              status: status,
              result: item.result // 保存原始result值用于筛选
            };
          });
          
          this.setData({
            evaluations: formattedEvaluations,
            filteredEvaluations: formattedEvaluations
          });
          this.applyFilters();
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
        // 使用空数据作为备用
        this.setData({
          evaluations: [],
          filteredEvaluations: []
        });
      }
    });
  },

  // 切换筛选条件
  changeFilter(e) {
    const filterId = e.currentTarget.dataset.id;
    this.setData({
      currentFilter: filterId
    });
    this.applyFilters();
  },

  // 搜索输入
  onSearchInput(e) {
    this.setData({
      searchText: e.detail.value
    });
  },

  // 执行搜索
  onSearch() {
    this.applyFilters();
  },

  // 应用筛选条件
  applyFilters() {
    let filtered = [...this.data.evaluations];
    
    // 应用状态筛选
    if (this.data.currentFilter === 'pending') {
      filtered = filtered.filter(item => item.status === '待评教');
    } else if (this.data.currentFilter === 'completed') {
      filtered = filtered.filter(item => item.status === '已完成');
    }
    
    // 应用搜索筛选
    if (this.data.searchText) {
      const searchLower = this.data.searchText.toLowerCase();
      filtered = filtered.filter(item => 
        item.courseName.toLowerCase().includes(searchLower) ||
        item.teacherName.toLowerCase().includes(searchLower)
      );
    }
    
    // 更新筛选后的列表
    this.setData({
      filteredEvaluations: filtered
    });
  },

  // 跳转到评教页面
  goToEvaluation(e) {
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