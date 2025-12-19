package routes

import (
	"tems-web-api/assets"
	"tems-web-api/controllers"
	"tems-web-api/middleware"
	"tems-web-api/uploads"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置所有路由
func SetupRoutes(r *gin.Engine) {
	// 添加CORS中间件
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 健康检查路由
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "TEMS Web API is running",
		})
	})

	// 使用嵌入的静态资源
	r.StaticFS("/manager/static", assets.FileSystem)
	r.StaticFS("/upload", assets.UploadFileSystem)

	// API版本分组
	v1 := r.Group("/api/v1")
	v1.Use(middleware.JWTAuthMiddleware(1))
	{
		// 轮播图相关路由 - 仅前端展示用接口
		bannerRoutes := v1.Group("/banners")
		{
			bannerRoutes.GET("", controllers.GetBanners)
		}

		// 学生相关路由
		studentRoutes := v1.Group("/students")
		{
			studentRoutes.POST("/bind", controllers.BindStudent)
			studentRoutes.POST("/login", controllers.StudentLogin)
			studentRoutes.GET("/:id", controllers.GetStudentByID)
		}

		// 评教相关路由
		evaluationRoutes := v1.Group("/evaluations")
		{
			evaluationRoutes.GET("/tasks", controllers.GetEvaluationTasks)
			evaluationRoutes.POST("/enhanced-tasks", controllers.GetEnhancedEvaluationList) // 增强版评课列表
			evaluationRoutes.GET("/questions/:setId", controllers.GetQuestionsBySetID)
			evaluationRoutes.POST("/submitResults", controllers.SubmitEvaluationResult)
			evaluationRoutes.GET("/statistics", controllers.GetStatisticsResults)
		}
	}

	// 管理后台接口分组
	// 获取管理后台首页HTML
	r.GET("/manager/", controllers.GetIndexHtml)
	// 验证码相关接口（不需要JWT验证）
	r.GET("/manager/verifycode", controllers.GenerateCaptcha) // 获取验证码图片

	manager := r.Group("/manager/api")
	manager.Use(middleware.JWTAuthMiddleware(0))
	{
		// 管理员登录接口（不需要JWT验证）
		manager.POST("/login", controllers.AdminLogin)
		manager.POST("/upload", uploads.UploadHandler())

		// 轮播图管理接口（实际项目中应该添加权限控制）
		bannerRoutes := manager.Group("/banners")
		{
			bannerRoutes.GET("/list", controllers.GetAllBanners)
			bannerRoutes.POST("/add", controllers.CreateBanner)
			bannerRoutes.POST("/update", controllers.UpdateBanner)
			bannerRoutes.DELETE("/delete/:id", controllers.DeleteBanner)
			bannerRoutes.POST("/visibility", controllers.UpdateBannerVisibility)
		}

		// 管理员用户管理接口（实际项目中应该添加权限控制）
		adminUserRoutes := manager.Group("/users")
		{
			adminUserRoutes.GET("/list", controllers.GetAdminUsers)
			adminUserRoutes.POST("/add", controllers.CreateAdminUser)
			adminUserRoutes.POST("/update", controllers.UpdateAdminUser)
			adminUserRoutes.DELETE("/delete/:id", controllers.DeleteAdminUser)
			adminUserRoutes.POST("/:id/status", controllers.UpdateAdminUserStatus)
			adminUserRoutes.POST("/:id/password", controllers.UpdateAdminUserPassword)
			adminUserRoutes.GET("/info", controllers.GetAdminUserInfo)
			adminUserRoutes.GET("/checkLogin", controllers.CheckAdminLogin)
		}

		// 评教任务管理接口
		evalTaskRoutes := manager.Group("/evaluation/task")
		{
			evalTaskRoutes.GET("/list", controllers.GetEvaluationTaskList)
			evalTaskRoutes.GET("/detail/:id", controllers.GetEvaluationTaskDetail)
			evalTaskRoutes.POST("/add", controllers.CreateEvaluationTask)
			evalTaskRoutes.POST("/update", controllers.UpdateEvaluationTask)
			evalTaskRoutes.POST("/delete/:id", controllers.DeleteEvaluationTask)
			evalTaskRoutes.GET("/all", controllers.GetAllEvaluationTasks)
			evalTaskRoutes.GET("/active", controllers.GetActiveEvaluationTasks)
		}

		// 评教问题集管理
		evalQuestionSetGroup := manager.Group("/evaluation/question-set")
		{
			evalQuestionSetGroup.GET("/list", controllers.GetEvaluationQuestionSetList)
			evalQuestionSetGroup.GET("/detail/:id", controllers.GetEvaluationQuestionSetDetail)
			evalQuestionSetGroup.POST("/add", controllers.CreateEvaluationQuestionSet)
			evalQuestionSetGroup.POST("/update", controllers.UpdateEvaluationQuestionSet)
			evalQuestionSetGroup.POST("/delete/:id", controllers.DeleteEvaluationQuestionSet)
			evalQuestionSetGroup.GET("/all", controllers.GetAllEvaluationQuestionSets)
		}

		// 评教问题详情管理
		evalQuestionDetailGroup := manager.Group("/evaluation/question-detail")
		{
			evalQuestionDetailGroup.GET("/list", controllers.GetEvaluationQuestionDetailList)
			evalQuestionDetailGroup.GET("/detail/:id", controllers.GetEvaluationQuestionDetail)
			evalQuestionDetailGroup.GET("/by-set/:setId", controllers.GetQuestionsDetailBySetID)
			evalQuestionDetailGroup.POST("/add", controllers.CreateEvaluationQuestionDetail)
			evalQuestionDetailGroup.POST("/update", controllers.UpdateEvaluationQuestionDetail)
			evalQuestionDetailGroup.POST("/delete/:id", controllers.DeleteEvaluationQuestionDetail)
		}

		// 评教任务详情管理
		evalTaskDetailGroup := manager.Group("/evaluation/task-detail")
		{
			evalTaskDetailGroup.GET("/list", controllers.GetEvaluationTaskDetailList)
			evalTaskDetailGroup.GET("/detail/:id", controllers.GetEvaluationTaskDetailById)
			evalTaskDetailGroup.GET("/by-task-class-course", controllers.GetEvaluationTaskDetailByTaskClassCourse)
			evalTaskDetailGroup.POST("/add", controllers.CreateEvaluationTaskDetail)
			evalTaskDetailGroup.POST("/update", controllers.UpdateEvaluationTaskDetail)
			evalTaskDetailGroup.POST("/delete/:id", controllers.DeleteEvaluationTaskDetail)
			evalTaskDetailGroup.POST("/batch", controllers.CreateEvaluationTaskDetailBatch)
		}

		// 班级管理接口
		classGroup := manager.Group("/classes")
		{
			classGroup.GET("/list", controllers.GetClassList)
			classGroup.GET("/detail/:id", controllers.GetClassByID)
			classGroup.POST("/add", controllers.CreateClass)
			classGroup.POST("/update", controllers.UpdateClass)
			classGroup.POST("/delete/:id", controllers.DeleteClass)
		}

		// 课程管理接口
		courseGroup := manager.Group("/courses")
		{
			courseGroup.GET("/list", controllers.GetCourseList)
			courseGroup.GET("/detail/:id", controllers.GetCourseByID)
			courseGroup.POST("/add", controllers.CreateCourse)
			courseGroup.POST("/update", controllers.UpdateCourse)
			courseGroup.POST("/delete/:id", controllers.DeleteCourse)
		}

		// 班级课程管理接口
		classCourseGroup := manager.Group("/class/course")
		{
			classCourseGroup.GET("/list", controllers.GetClassCourseList)
			classCourseGroup.GET("/detail/:id", controllers.GetClassCourseByID)
			classCourseGroup.POST("/add", controllers.CreateClassCourse)
			classCourseGroup.POST("/update", controllers.UpdateClassCourse)
			classCourseGroup.POST("/delete/:id", controllers.DeleteClassCourse)
		}

		// 学生管理接口
		studentGroup := manager.Group("/students")
		{
			studentGroup.GET("/list", controllers.GetStudentList)
			studentGroup.POST("/add", controllers.AddStudent)
			studentGroup.POST("/batch", controllers.AddBatchStudent)
			studentGroup.POST("/update", controllers.UpdateStudent)
			studentGroup.POST("/delete/:id", controllers.DeleteStudent)
		}

		// 拼音转换接口
		manager.POST("/get-pinyin", controllers.GetPinyin)

		// 统计相关接口
		statRoutes := manager.Group("/statistics")
		{
			statRoutes.GET("/list", controllers.GetEvaluationDetailStatistics)
			statRoutes.GET("/export", controllers.ExportEvaluationDetailStatistics)
			statRoutes.GET("/top-card", controllers.HomeTopCardStatistics)
		}
	}
}
