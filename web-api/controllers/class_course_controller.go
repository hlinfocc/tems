package controllers

import (
	"log"
	"net/http"
	"strconv"
	"tems-web-api/models"

	"github.com/gin-gonic/gin"
)

// ClassCourseRequest 班级课程请求参数
type ClassCourseRequest struct {
	Page     int    `form:"page" json:"page"`
	Limit    int    `form:"limit" json:"limit"`
	ClassID  uint64 `form:"classId" json:"classId"`
	CourseID uint64 `form:"courseId" json:"courseId"`
	Keyword  string `form:"keyword" json:"keyword"`
}

// ClassCourseForm 班级课程表单数据
type ClassCourseForm struct {
	ID          uint64 `json:"id"`
	ClassID     uint64 `json:"classID" binding:"required"`
	CourseID    uint64 `json:"courseID" binding:"required"`
	TeacherID   uint64 `json:"teacherID" binding:"required"`
	CourseName  string `json:"courseName" binding:"required"`
	TeacherName string `json:"teacherName" binding:"required"`
}

// GetClassCourseList 获取班级课程列表（分页，支持筛选和搜索）
func GetClassCourseList(c *gin.Context) {
	var req ClassCourseRequest
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

	// 分页查询（支持筛选和搜索）
	classCourses, count, err := models.QueryClassCoursesWithPagination(
		req.Page,
		req.Limit,
		req.ClassID,
		req.CourseID,
		req.Keyword,
	)

	if err != nil {
		log.Printf("分页查询班级课程列表失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "获取班级课程列表失败",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":  200,
		"msg":   "成功",
		"data":  classCourses,
		"count": count,
	})
}

// GetClassCourseByID 获取单个班级课程详情
func GetClassCourseByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的班级课程ID",
		})
		return
	}

	classCourse, err := models.GetClassCourseByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "班级课程不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": classCourse,
	})
}

// CreateClassCourse 创建班级课程关系
func CreateClassCourse(c *gin.Context) {
	var form ClassCourseForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "参数错误",
			"error": err.Error(),
		})
		return
	}

	classCourse := &models.ClassCourse{
		ClassID:     form.ClassID,
		CourseID:    form.CourseID,
		TeacherID:   form.TeacherID,
		CourseName:  form.CourseName,
		TeacherName: form.TeacherName,
	}
	// GetCourseByClassCourseID
	exists, ee := models.GetCourseByClassCourseID(form.ClassID, form.CourseID)
	if ee != nil {
		log.Printf("查询班级课程关系失败: %v", ee)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "查询班级课程关系失败",
			"error": ee.Error(),
		})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{
			"code": 409,
			"msg":  "班级课程关系已存在",
		})
		return
	}

	if err := models.CreateClassCourse(classCourse); err != nil {
		log.Printf("创建班级课程关系失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "创建班级课程关系失败",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "创建成功",
		"data": classCourse,
	})
}

// UpdateClassCourse 更新班级课程关系
func UpdateClassCourse(c *gin.Context) {
	var form ClassCourseForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "参数错误",
			"error": err.Error(),
		})
		return
	}

	// 检查班级课程是否存在
	classCourse, err := models.GetClassCourseByID(uint(form.ID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "班级课程不存在",
		})
		return
	}

	// 更新班级课程信息
	classCourse.ClassID = form.ClassID
	classCourse.CourseID = form.CourseID
	classCourse.TeacherID = form.TeacherID
	classCourse.CourseName = form.CourseName
	classCourse.TeacherName = form.TeacherName

	if err := models.UpdateClassCourse(classCourse); err != nil {
		log.Printf("更新班级课程关系失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "更新班级课程关系失败",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "更新成功",
		"data": classCourse,
	})
}

// DeleteClassCourse 删除班级课程关系
func DeleteClassCourse(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的班级课程ID",
		})
		return
	}

	// 检查班级课程是否存在
	_, err = models.GetClassCourseByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "班级课程不存在",
		})
		return
	}

	// 删除班级课程关系（软删除）
	if err := models.DeleteClassCourse(uint(id)); err != nil {
		log.Printf("删除班级课程关系失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "删除班级课程关系失败",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}
