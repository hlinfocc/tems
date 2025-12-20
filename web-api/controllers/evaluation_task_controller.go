package controllers

import (
	"log"
	"net/http"
	"strconv"
	"tems-web-api/models"

	"github.com/gin-gonic/gin"
)

// EvaluationTaskRequest 评教任务请求参数
type EvaluationTaskRequest struct {
	Page         int    `form:"page" json:"page"`
	Limit        int    `form:"limit" json:"limit"`
	Keyword      string `form:"keyword" json:"keyword"`
	AcademicYear string `form:"academic_year" json:"academic_year"`
	Semester     *int   `form:"semester" json:"semester"`
}

// EvaluationTaskForm 评教任务表单数据
type EvaluationTaskForm struct {
	ID              uint64          `json:"id"`
	TaskName        string          `json:"taskName" binding:"required"`
	AcademicYear    string          `json:"academicYear" binding:"required"`
	Semester        *int            `json:"semester" binding:"required"`
	QuestionSetID   uint64          `json:"questionSetID" binding:"required"`
	QuestionSetName string          `json:"questionSetName" binding:"required"`
	StartTime       models.JsonTime `json:"startTime" binding:"required"`
	EndTime         models.JsonTime `json:"endTime" binding:"required"`
}

// GetEvaluationTaskList 获取评教任务列表
func GetEvaluationTaskList(c *gin.Context) {
	var req EvaluationTaskRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "参数错误",
			"error": err.Error(),
		})
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	// 分页查询
	tasks, count, err := models.QueryEvaluationTasksWithPagination(req.Page, req.Limit, req.Keyword, req.AcademicYear, req.Semester)
	if err != nil {
		log.Printf("查询评教任务列表失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "获取评教任务列表失败",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":  200,
		"msg":   "成功",
		"data":  tasks,
		"count": count,
	})
}

// GetEvaluationTaskDetail 获取评教任务详情
func GetEvaluationTaskDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的任务ID",
		})
		return
	}

	task, err := models.GetEvaluationTaskByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "评教任务不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": task,
	})
}

// CreateEvaluationTask 创建评教任务
func CreateEvaluationTask(c *gin.Context) {
	var form EvaluationTaskForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "参数错误",
			"error": err.Error(),
		})
		return
	}

	// 创建任务
	task := &models.EvaluationTask{
		TaskName:        form.TaskName,
		AcademicYear:    form.AcademicYear,
		Semester:        *form.Semester,
		QuestionSetID:   form.QuestionSetID,
		QuestionSetName: form.QuestionSetName,
		StartTime:       form.StartTime,
		EndTime:         form.EndTime,
	}

	if err := models.CreateEvaluationTask(task); err != nil {
		log.Printf("创建评教任务失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "创建评教任务失败",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "创建评教任务成功",
	})
}

// UpdateEvaluationTask 更新评教任务
func UpdateEvaluationTask(c *gin.Context) {
	var form EvaluationTaskForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "参数错误",
			"error": err.Error(),
		})
		return
	}

	// 更新任务
	task := &models.EvaluationTask{
		BaseModel: models.BaseModel{
			ID: form.ID,
		},
		TaskName:        form.TaskName,
		AcademicYear:    form.AcademicYear,
		Semester:        *form.Semester,
		QuestionSetID:   form.QuestionSetID,
		QuestionSetName: form.QuestionSetName,
		StartTime:       form.StartTime,
		EndTime:         form.EndTime,
	}

	if err := models.UpdateEvaluationTask(task); err != nil {
		log.Printf("更新评教任务失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "更新评教任务失败",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "更新评教任务成功",
	})
}

// DeleteEvaluationTask 删除评教任务
func DeleteEvaluationTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的任务ID",
		})
		return
	}

	// 先删除相关的任务详情
	if err := models.DeleteTaskDetailsByTaskID(uint(id)); err != nil {
		log.Printf("删除评教任务详情失败: %v", err)
	}

	// 删除任务
	if err := models.DeleteEvaluationTask(id); err != nil {
		log.Printf("删除评教任务失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "删除评教任务失败",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除评教任务成功",
	})
}

// GetAllEvaluationTasks 获取所有评教任务（用于下拉选择）
func GetAllEvaluationTasks(c *gin.Context) {
	tasks, err := models.GetAllEvaluationTasks()
	if err != nil {
		log.Printf("获取所有评教任务失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "获取评教任务列表失败",
			"error": err.Error(),
		})
		return
	}

	// 转换为选项格式
	options := make([]map[string]interface{}, 0, len(tasks))
	for _, task := range tasks {
		options = append(options, map[string]interface{}{
			"id":        task.ID,
			"task_name": task.TaskName,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": options,
	})
}

// GetActiveEvaluationTasks 获取当前有效的评教任务
func GetActiveEvaluationTasks(c *gin.Context) {
	tasks, err := models.GetActiveTasks()
	if err != nil {
		log.Printf("获取有效评教任务失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "获取有效评教任务失败",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": tasks,
	})
}
