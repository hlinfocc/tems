package controllers

import (
	"net/http"
	"strconv"
	"time"

	"tems-web-api/models"

	"github.com/gin-gonic/gin"
)

// EvaluationQuestionSetRequest 评教问题集请求参数
type EvaluationQuestionSetRequest struct {
	Page     int    `form:"page" json:"page"`
	Limit    int    `form:"limit" json:"limit"`
	IsPage   int    `form:"isPage" json:"isPage"`
	CurrYear int    `form:"currYear" json:"currYear"`
	Keyword  string `form:"keyword" json:"keyword"`
	Status   string `form:"status" json:"status"`
}

// EvaluationQuestionSetForm 评教问题集表单数据
type EvaluationQuestionSetForm struct {
	ID     uint64 `json:"id" form:"id"`
	Name   string `json:"name" form:"name" binding:"required"`
	Remark string `json:"remark" form:"remark"`
	Status int    `json:"status" form:"status"`
}

// GetEvaluationQuestionSetList 获取评教问题集列表
func GetEvaluationQuestionSetList(c *gin.Context) {
	var req EvaluationQuestionSetRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	// 初始化分页参数
	page := req.Page
	if page <= 0 {
		page = 1
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	// 构建查询条件
	db := models.GetDB()
	var questionSets []models.EvaluationQuestionSet
	var count int64

	query := db.Model(&models.EvaluationQuestionSet{}).Where("is_deleted = ?", false)

	// 如果有搜索关键字
	if req.Keyword != "" {
		query = query.Where("name LIKE ?", "%"+req.Keyword+"%")
	}

	// 如果有状态筛选
	if req.Status != "" {
		status, _ := strconv.Atoi(req.Status)
		query = query.Where("status = ?", status)
	}

	// 获取总数
	query.Count(&count)

	if req.IsPage == 0 {
		// 获取分页数据
		if err := query.Offset(offset).Limit(limit).Order("createdAt DESC").Find(&questionSets).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "获取评教问题集列表失败: " + err.Error(),
			})
			return
		}
	} else {
		// 获取所有数据
		if req.CurrYear > 0 {
			query = query.Where("year = ?", time.Now().Year())
		}
		if err := query.Order("createdAt DESC").Find(&questionSets).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 500,
				"msg":  "获取评教问题集列表失败: " + err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":  200,
		"data":  questionSets,
		"count": count,
		"msg":   "获取评教问题集列表成功",
	})
}

// GetEvaluationQuestionSetDetail 获取评教问题集详情
func GetEvaluationQuestionSetDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "无效的问题集ID",
		})
		return
	}

	questionSet, err := models.GetEvaluationQuestionSetByID(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 404,
			"msg":  "评教问题集不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": questionSet,
		"msg":  "获取评教问题集详情成功",
	})
}

// CreateEvaluationQuestionSet 创建评教问题集
func CreateEvaluationQuestionSet(c *gin.Context) {
	var form EvaluationQuestionSetForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	questionSet := models.EvaluationQuestionSet{
		Name:   form.Name,
		Remark: form.Remark,
		Status: form.Status,
		Year:   time.Now().Year(),
	}

	if err := models.CreateEvaluationQuestionSet(&questionSet); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "创建评教问题集失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "创建评教问题集成功",
	})
}

// UpdateEvaluationQuestionSet 更新评教问题集
func UpdateEvaluationQuestionSet(c *gin.Context) {
	var form EvaluationQuestionSetForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	// 检查问题集是否存在
	questionSet, err := models.GetEvaluationQuestionSetByID(form.ID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 404,
			"msg":  "评教问题集不存在",
		})
		return
	}

	// 更新字段
	questionSet.Name = form.Name
	questionSet.Remark = form.Remark
	questionSet.Status = form.Status

	if err := models.UpdateEvaluationQuestionSet(questionSet); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "更新评教问题集失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "更新评教问题集成功",
	})
}

// DeleteEvaluationQuestionSet 删除评教问题集
func DeleteEvaluationQuestionSet(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "无效的问题集ID",
		})
		return
	}

	// 检查问题集是否存在
	_, err = models.GetEvaluationQuestionSetByID(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 404,
			"msg":  "评教问题集不存在",
		})
		return
	}

	if err := models.DeleteEvaluationQuestionSet(id); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "删除评教问题集失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除评教问题集成功",
	})
}

// GetAllEvaluationQuestionSets 获取所有评教问题集（用于下拉选择）
func GetAllEvaluationQuestionSets(c *gin.Context) {
	questionSets, err := models.GetAllQuestionSets()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "获取评教问题集列表失败: " + err.Error(),
		})
		return
	}

	// 转换为下拉选项格式
	options := make([]map[string]interface{}, len(questionSets))
	for i, set := range questionSets {
		options[i] = map[string]interface{}{
			"id":   set.ID,
			"name": set.Name,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": options,
		"msg":  "获取评教问题集列表成功",
	})
}
