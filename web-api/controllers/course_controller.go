package controllers

import (
	"log"
	"net/http"
	"strconv"
	"tems-web-api/models"

	"github.com/gin-gonic/gin"
)

// CourseRequest 课程请求参数
type CourseRequest struct {
	Page       int    `form:"page" json:"page"`
	Limit      int    `form:"limit" json:"limit"`
	CourseName string `form:"courseName" json:"courseName"`
}

// CourseForm 课程表单数据
type CourseForm struct {
	ID         uint64 `json:"id"`
	CourseName string `json:"courseName" binding:"required"`
	Status     int    `json:"status"`
}

// GetCourseList 获取课程列表（分页，支持课程名称模糊查询）
func GetCourseList(c *gin.Context) {
	var req CourseRequest
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

	// 分页查询（支持课程名称模糊查询）
	courses, count, err := models.QueryCoursesWithPagination(
		req.Page,
		req.Limit,
		req.CourseName,
	)

	if err != nil {
		log.Printf("分页查询课程列表失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "获取课程列表失败",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":  200,
		"msg":   "成功",
		"data":  courses,
		"count": count,
	})
}

// GetCourseByID 获取单个课程信息
func GetCourseByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的课程ID",
		})
		return
	}

	course, err := models.GetCourseByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "课程不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": course,
	})
}

// CreateCourse 创建课程
func CreateCourse(c *gin.Context) {
	var form CourseForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "参数错误",
			"error": err.Error(),
		})
		return
	}

	course := &models.Course{
		CourseName: form.CourseName,
		Status:     form.Status,
	}

	if err := models.CreateCourse(course); err != nil {
		log.Printf("创建课程失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "创建课程失败",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "创建成功",
		"data": course,
	})
}

// UpdateCourse 更新课程信息
func UpdateCourse(c *gin.Context) {
	var form CourseForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "参数错误",
			"error": err.Error(),
		})
		return
	}

	// 检查课程是否存在
	course, err := models.GetCourseByID(form.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "课程不存在",
		})
		return
	}

	// 更新课程信息
	course.CourseName = form.CourseName
	course.Status = form.Status

	if err := models.UpdateCourse(course); err != nil {
		log.Printf("更新课程失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "更新课程失败",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "更新成功",
		"data": course,
	})
}

// DeleteCourse 删除课程
func DeleteCourse(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的课程ID",
		})
		return
	}

	// 检查课程是否存在
	_, err = models.GetCourseByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "课程不存在",
		})
		return
	}

	// 删除课程（软删除）
	if err := models.DeleteCourse(id); err != nil {
		log.Printf("删除课程失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "删除课程失败",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}
