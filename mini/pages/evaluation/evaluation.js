// evaluation.js
Page({
  data: {
    currentCourse: {
      courseName: "",
      teacherName: "",
    },
    questions: [],
    currentIndex: 1,
    totalQuestions: 0,
    currentQuestion: {},
    selectedOption: "",
    comment: "",
    answers: [],
    evaluation: {},
    userId: 0,
  },

  onLoad(options) {
    const { id, evaluation: encodedEvaluation } = options;
    const evaluation = JSON.parse(decodeURIComponent(encodedEvaluation));
    console.log("evaluation:", evaluation);
    this.setData({
      evaluation,
      userId: getApp().globalData.studentInfo.userId,
    });
    // this.loadCourseData(id);
    this.loadEvaluationQuestions(evaluation.questionSetID);
  },

  // 加载课程信息
  loadCourseData(id) {
    // 模拟加载课程数据
    const courses = [];

    const course = courses.find((c) => c.id === parseInt(id)) || courses[0];
    this.setData({
      currentCourse: course,
    });

    // 实际项目中应该调用API
    /*
    const app = getApp();
    wx.request({
      url: app.globalData.baseUrl + '/courses/' + id,
      success: (res) => {
        this.setData({
          currentCourse: res.data
        });
      }
    });
    */
  },

  // 加载评教问题
  loadEvaluationQuestions(qsid) {
    // 模拟评教问题数据
    const questions = [];

    this.setData({
      questions: questions,
      totalQuestions: questions.length,
      currentQuestion: questions[0],
      selectedOption: this.data.answers[1] || "",
      comment: "",
    });

    // 实际项目中应该调用API
    // 获取token
    const token = wx.getStorageSync("token") || "";

    const app = getApp();
    wx.request({
      url: app.globalData.baseUrl + "/v1/evaluations/questions/" + qsid,
      header: {
        "content-type": "application/json",
        Authorization: `${token}`,
      },
      success: (res) => {
        console.log("评教问题接口返回结果:", res);
        if (res.statusCode === 200 && res.data.code === 200 && res.data.data) {
          let questions = res.data.data;
          // 转换数据格式，确保每个问题都有options属性
          questions = questions.map((q) => ({
            ...q,
            options: JSON.parse(q.options || "[]") || [],
          }));
          console.log("questions:", questions);
          this.setData({
            questions: questions,
            totalQuestions: questions.length,
            currentQuestion: questions[0],
          });
        } else {
          console.error(
            "评教问题接口调用失败:",
            res.data.message || "未知错误"
          );
          // 使用模拟评教问题数据作为备用
          this.setData({
            questions: [],
            totalQuestions: 0,
            currentQuestion: {},
          });
        }
      },
    });
  },

  // 选择选项
  selectOption(e) {
    const value = e.currentTarget.dataset.value;
    const qindex = this.data.currentIndex - 1;
    console.log("value:", qindex, value);
    this.setData({
      selectedOption: value,
    });

    // 保存答案
    const answers = this.data.answers;
    const answerCurr = {
      taskId: this.data.evaluation.taskId,
      taskDetailId: this.data.evaluation.id,
      questionSetId: this.data.evaluation.questionSetID,
      semester: this.data.evaluation.semester,
      academicYear: this.data.evaluation.academicYear,
      questionDetailId: this.data.currentQuestion.id,
      questionType: this.data.currentQuestion.questionType,
      studentId: this.data.userId,
      answer: value,
      result: -1,
    };
    answers[qindex] = answerCurr;
    this.setData({
      answers: answers,
    });
  },

  // 输入评论
  inputComment(e) {
    this.setData({
      comment: e.detail.value,
    });
  },

  // 上一题
  prevQuestion() {
    if (this.data.currentIndex > 1) {
      const newIndex = this.data.currentIndex - 1;
      this.setData({
        currentIndex: newIndex,
        currentQuestion: this.data.questions[newIndex - 1],
        selectedOption: this.data.answers[newIndex] || "",
        comment: "",
      });
    }
  },

  // 下一题或提交
  nextQuestion() {
    // 检查是否已选择答案
    if (!this.data.selectedOption) {
      wx.showToast({
        title: "请选择一个选项",
        icon: "none",
      });
      return;
    }

    // 如果是最后一题，提交评价
    if (this.data.currentIndex === this.data.totalQuestions) {
      this.submitEvaluation();
    } else {
      // 否则切换到下一题
      const newIndex = this.data.currentIndex + 1;
      this.setData({
        currentIndex: newIndex,
        currentQuestion: this.data.questions[newIndex - 1],
        selectedOption: this.data.answers[newIndex] || "",
        comment: "",
      });
    }
  },

  // 提交评价
  submitEvaluation() {
    wx.showLoading({
      title: "提交中...",
    });

    // 构造提交数据
    const submitData = {
      courseId: this.data.currentCourse.id,
      answers: this.data.answers,
      comment: this.data.comment,
      submitTime: new Date().toISOString(),
    };
    console.log("submitData:", submitData);
    // 模拟提交成功
    // setTimeout(() => {
    //   wx.hideLoading();
    //   wx.showToast({
    //     title: "评价成功",
    //     icon: "success",
    //   });

    //   // 延迟返回首页
    //   setTimeout(() => {
    //     wx.navigateBack();
    //   }, 1500);
    // }, 1000);

    // 实际项目中应该调用API
    const token = wx.getStorageSync("token") || "";
    const app = getApp();
    wx.request({
      url: app.globalData.baseUrl + "/v1/evaluations/submitResults",
      method: "POST",
      header: {
        "content-type": "application/json",
        Authorization: `${token}`,
      },
      data: this.data.answers,
      success: (res) => {
        wx.hideLoading();
        if (res.statusCode === 200 && res.data.data.code === 200) {
          wx.hideLoading();
          wx.showToast({
            title: "提交成功",
            icon: "success",
          });
          setTimeout(() => {
            wx.navigateBack();
          }, 1500);
        } else {
          wx.showToast({
            title: res.data.data.message || "提交失败，请重试",
            icon: "none",
          });
        }
        wx.hideLoading();
      },
      fail: (err) => {
        wx.hideLoading();
        wx.showToast({
          title: "网络错误，请重试",
          icon: "none",
        });
      },
    });
  },
});
