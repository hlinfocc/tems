package controllers

import (
	"log"
	"net/http"
	"strconv"
	"tems-web-api/models"

	"github.com/gin-gonic/gin"
)

// ClassRequest 班级请求参数
type ClassRequest struct {
	Page    int    `form:"page" json:"page"`
	Limit   int    `form:"limit" json:"limit"`
	Grade   string `form:"grade" json:"grade"`
	MajorID uint64 `form:"majorId" json:"majorId"`
	Keyword string `form:"keyword" json:"keyword"`
}

// ClassForm 班级表单数据
type ClassForm struct {
	ID              uint64 `json:"id"`
	ClassName       string `json:"className" binding:"required"`
	Grade           string `json:"grade"`
	MajorName       string `json:"majorName"`
	MajorID         uint64 `json:"majorID"`
	HeadTeacherName string `json:"headTeacherName"`
	HeadTeacherID   uint64 `json:"headTeacherId"`
}

// GetClassList 获取班级列表（分页）
func GetClassList(c *gin.Context) {
	var req ClassRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"msg": "参数错误",
			"error":   err.Error(),
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
	classes, count, err := models.QueryClassesWithPagination(
		req.Page,
		req.Limit,
		req.Keyword,
	)

	if err != nil {
		log.Printf("分页查询班级列表失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"msg": "获取班级列表失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg": "成功",
		"data":    classes,
		"count":   count,
	})
}

// GetClassByID 获取单个班级信息
func GetClassByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"msg": "无效的班级ID",
		})
		return
	}

	class, err := models.GetClassByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"msg": "班级不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg": "成功",
		"data":    class,
	})
}

// CreateClass 创建班级
func CreateClass(c *gin.Context) {
	var form ClassForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"msg": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	class := &models.Class{
		ClassName:       form.ClassName,
		Grade:           form.Grade,
		MajorName:       form.MajorName,
		MajorID:         form.MajorID,
		HeadTeacherName: form.HeadTeacherName,
		HeadTeacherID:   form.HeadTeacherID,
	}

	if err := models.CreateClass(class); err != nil {
		log.Printf("创建班级失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"msg": "创建班级失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg": "创建成功",
		"data":    class,
	})
}

// UpdateClass 更新班级信息
func UpdateClass(c *gin.Context) {
	var form ClassForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"msg": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 检查班级是否存在
	class, err := models.GetClassByID(form.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"msg": "班级不存在",
		})
		return
	}

	// 更新班级信息
	class.ClassName = form.ClassName
	class.Grade = form.Grade
	class.MajorName = form.MajorName
	class.MajorID = form.MajorID
	class.HeadTeacherName = form.HeadTeacherName
	class.HeadTeacherID = form.HeadTeacherID

	if err := models.UpdateClass(class); err != nil {
		log.Printf("更新班级失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"msg": "更新班级失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg": "更新成功",
		"data":    class,
	})
}

// DeleteClass 删除班级
func DeleteClass(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"msg": "无效的班级ID",
		})
		return
	}

	// 检查班级是否存在
	_, err = models.GetClassByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"msg": "班级不存在",
		})
		return
	}

	// 删除班级（软删除）
	if err := models.DeleteClass(id); err != nil {
		log.Printf("删除班级失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"msg": "删除班级失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg": "删除成功",
	})
}
