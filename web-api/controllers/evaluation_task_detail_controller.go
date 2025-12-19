package controllers

import (
	"log"
	"net/http"
	"strconv"
	"tems-web-api/models"

	"github.com/gin-gonic/gin"
)

// EvaluationTaskDetailRequest 评教任务详情请求参数
type EvaluationTaskDetailRequest struct {
	Page    int    `form:"page" json:"page"`
	Limit   int    `form:"limit" json:"limit"`
	TaskID  uint64 `form:"taskId" json:"taskId"`
	ClassID uint64 `form:"classId" json:"classId"`
}

// EvaluationTaskDetailForm 评教任务详情表单数据
type EvaluationTaskDetailForm struct {
	ID         uint64 `json:"id"`
	TaskID     uint64 `json:"taskId" binding:"required"`
	ClassID    uint64 `json:"classId" binding:"required"`
	CourseID   uint64 `json:"courseId" binding:"required"`
	TeacherID  uint64 `json:"teacherId"`
	ClassName  string `json:"className"`
	CourseName string `json:"courseName"`
	UserName   string `json:"userName"`
}

// ClassInfo 班级信息
type ClassInfo struct {
	ClassID   uint64 `json:"classId"`
	ClassName string `json:"className"`
}

// EvaluationTaskDetailBatchForm 评教任务详情批量请求参数
type EvaluationTaskDetailBatchForm struct {
	TaskID  uint64      `form:"taskId" json:"taskId"`
	Classes []ClassInfo `form:"classes" json:"classes"`
}

// GetEvaluationTaskDetailList 获取评教任务详情列表（分页）
func GetEvaluationTaskDetailList(c *gin.Context) {
	var req EvaluationTaskDetailRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "参数错误",
			"error": err.Error(),
		})
		return
	}

	// 验证TaskID是否必传
	if req.TaskID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "任务ID(TaskID)是必传参数",
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
	details, count, err := models.QueryEvaluationTaskDetailsWithPagination(
		req.Page,
		req.Limit,
		uint(req.TaskID),
		uint(req.ClassID),
	)

	if err != nil {
		log.Printf("分页查询评教任务详情失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "获取评教任务详情列表失败",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":  200,
		"msg":   "成功",
		"data":  details,
		"count": count,
	})
}

// GetEvaluationTaskDetailById 获取单个评教任务详情
func GetEvaluationTaskDetailById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的详情ID",
		})
		return
	}

	detail, err := models.GetEvaluationTaskDetailByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "评教任务详情不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": detail,
	})
}

// CreateEvaluationTaskDetail 创建评教任务详情
func CreateEvaluationTaskDetail(c *gin.Context) {
	var form EvaluationTaskDetailForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "参数错误",
			"error": err.Error(),
		})
		return
	}

	// 创建任务详情
	detail := &models.EvaluationTaskDetail{
		TaskID:     form.TaskID,
		ClassID:    form.ClassID,
		CourseID:   form.CourseID,
		TeacherID:  form.TeacherID,
		ClassName:  form.ClassName,
		CourseName: form.CourseName,
		UserName:   form.UserName,
	}

	// 检查是否已存在相同的任务详情
	existingDetail, err := models.GetTaskDetailByTaskClassCourse(uint(form.TaskID), uint(form.ClassID), uint(form.CourseID))
	if err == nil && existingDetail != nil {
		c.JSON(http.StatusConflict, gin.H{
			"code": 409,
			"msg":  "该任务、班级和课程的评教详情已存在",
		})
		return
	}

	if err := models.CreateEvaluationTaskDetail(detail); err != nil {
		log.Printf("创建评教任务详情失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "创建评教任务详情失败",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "创建评教任务详情成功",
	})
}

// CreateEvaluationTaskDetailBatch 创建评教任务详情批量
func CreateEvaluationTaskDetailBatch(c *gin.Context) {
	var form EvaluationTaskDetailBatchForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "参数错误",
			"error": err.Error(),
		})
		return
	}

	var details []*models.EvaluationTaskDetail
	for _, classes := range form.Classes {
		classCourse, e := models.GetCoursesByClassID(uint(classes.ClassID))
		if e != nil || len(classCourse) == 0 {
			log.Printf("获取班级课程失败: %v", e)
			continue
		}
		for _, course := range classCourse {
			detail := &models.EvaluationTaskDetail{
				TaskID:     form.TaskID,
				ClassID:    classes.ClassID,
				CourseID:   course.ID,
				ClassName:  classes.ClassName,
				CourseName: course.CourseName,
				TeacherID:  course.TeacherID,
				UserName:   course.TeacherName,
			}
			details = append(details, detail)
		}
	}
	if len(details) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "班级课程不能为空",
		})
		return
	}
	if err := models.CreateBatchEvaluationTaskDetail(details); err != nil {
		log.Printf("创建批量评教任务详情失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "保存失败",
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "保存成功",
	})
}

// UpdateEvaluationTaskDetail 更新评教任务详情
func UpdateEvaluationTaskDetail(c *gin.Context) {
	var form EvaluationTaskDetailForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "参数错误",
			"error": err.Error(),
		})
		return
	}

	// 先获取原记录
	oldDetail, err := models.GetEvaluationTaskDetailByID(uint(form.ID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "评教任务详情不存在",
		})
		return
	}

	// 更新字段
	oldDetail.ClassName = form.ClassName
	oldDetail.CourseName = form.CourseName
	oldDetail.UserName = form.UserName

	if err := models.UpdateEvaluationTaskDetail(oldDetail); err != nil {
		log.Printf("更新评教任务详情失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "更新评教任务详情失败",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "更新评教任务详情成功",
	})
}

// DeleteEvaluationTaskDetail 删除评教任务详情
func DeleteEvaluationTaskDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的详情ID",
		})
		return
	}

	if err := models.DeleteEvaluationTaskDetail(uint(id)); err != nil {
		log.Printf("删除评教任务详情失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "删除评教任务详情失败",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除评教任务详情成功",
	})
}

// GetEvaluationTaskDetailByTaskClassCourse 根据任务ID、班级ID和课程ID获取详情
func GetEvaluationTaskDetailByTaskClassCourse(c *gin.Context) {
	taskIDStr := c.Query("taskId")
	classIDStr := c.Query("classId")
	courseIDStr := c.Query("courseId")

	if taskIDStr == "" || classIDStr == "" || courseIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "请提供任务ID、班级ID和课程ID",
		})
		return
	}

	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的任务ID",
		})
		return
	}

	classID, err := strconv.ParseUint(classIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的班级ID",
		})
		return
	}

	courseID, err := strconv.ParseUint(courseIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的课程ID",
		})
		return
	}

	detail, err := models.GetTaskDetailByTaskClassCourse(uint(taskID), uint(classID), uint(courseID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "评教任务详情不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": detail,
	})
}
