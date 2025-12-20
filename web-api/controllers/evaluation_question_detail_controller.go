package controllers

import (
	"net/http"
	"strconv"

	"tems-web-api/models"

	"github.com/gin-gonic/gin"
)

// EvaluationQuestionDetailRequest 评教问题详情请求参数
type EvaluationQuestionDetailRequest struct {
	Page          int    `form:"page" json:"page"`
	Limit         int    `form:"limit" json:"limit"`
	Keyword       string `form:"keyword" json:"keyword"`
	QuestionSetID uint64 `form:"questionSetId" json:"questionSetId"`
	QuestionType  int    `form:"questionType" json:"questionType"`
}

// EvaluationQuestionDetailForm 评教问题详情表单数据
type EvaluationQuestionDetailForm struct {
	ID            uint64 `json:"id" form:"id"`
	QuestionSetID uint64 `json:"questionSetId" form:"questionSetId" binding:"required"`
	Title         string `json:"title" form:"title" binding:"required"`
	QuestionType  int    `json:"questionType" form:"questionType" binding:"required"`
	Options       string `json:"options" form:"options"`
	Remark        string `json:"remark" form:"remark"`
}

// GetEvaluationQuestionDetailList 获取评教问题详情列表
func GetEvaluationQuestionDetailList(c *gin.Context) {
	var req EvaluationQuestionDetailRequest
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
	var details []models.EvaluationQuestionDetail
	var count int64

	query := db.Model(&models.EvaluationQuestionDetail{}).Where("is_deleted = ?", false)

	// 如果有搜索关键字
	if req.Keyword != "" {
		query = query.Where("title LIKE ?", "%"+req.Keyword+"%")
	}

	// 如果有问题集ID筛选
	if req.QuestionSetID > 0 {
		query = query.Where("question_set_id = ?", req.QuestionSetID)
	}

	// 如果有问题类型筛选
	if req.QuestionType > 0 {
		query = query.Where("question_type = ?", req.QuestionType)
	}

	// 获取总数
	query.Count(&count)

	// 获取分页数据
	if err := query.Offset(offset).Limit(limit).Order("id ASC").Find(&details).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "获取评教问题详情列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":  200,
		"data":  details,
		"count": count,
		"msg":   "获取评教问题详情列表成功",
	})
}

// GetEvaluationQuestionDetail 获取评教问题详情
func GetEvaluationQuestionDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "无效的问题详情ID",
		})
		return
	}

	detail, err := models.GetEvaluationQuestionDetailByID(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 404,
			"msg":  "评教问题详情不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": detail,
		"msg":  "获取评教问题详情成功",
	})
}

// GetQuestionsDetailBySetID 根据问题集ID获取问题详情列表
func GetQuestionsDetailBySetID(c *gin.Context) {
	setIDStr := c.Param("setId")
	setID, err := strconv.ParseUint(setIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "无效的问题集ID",
		})
		return
	}

	details, err := models.GetQuestionsBySetID(setID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "获取问题列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": details,
		"msg":  "获取问题列表成功",
	})
}

// CreateEvaluationQuestionDetail 创建评教问题详情
func CreateEvaluationQuestionDetail(c *gin.Context) {
	var form EvaluationQuestionDetailForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	detail := models.EvaluationQuestionDetail{
		QuestionSetID: form.QuestionSetID,
		Title:         form.Title,
		QuestionType:  form.QuestionType,
		Options:       form.Options,
		Remark:        form.Remark,
	}

	if err := models.CreateEvaluationQuestionDetail(&detail); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "创建评教问题详情失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "创建评教问题详情成功",
	})
}

// UpdateEvaluationQuestionDetail 更新评教问题详情
func UpdateEvaluationQuestionDetail(c *gin.Context) {
	var form EvaluationQuestionDetailForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	// 检查问题详情是否存在
	detail, err := models.GetEvaluationQuestionDetailByID(form.ID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 404,
			"msg":  "评教问题详情不存在",
		})
		return
	}

	// 更新字段
	detail.QuestionSetID = form.QuestionSetID
	detail.Title = form.Title
	detail.QuestionType = form.QuestionType
	detail.Options = form.Options
	detail.Remark = form.Remark

	if err := models.UpdateEvaluationQuestionDetail(detail); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "更新评教问题详情失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "更新评教问题详情成功",
	})
}

// DeleteEvaluationQuestionDetail 删除评教问题详情
func DeleteEvaluationQuestionDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "无效的问题详情ID",
		})
		return
	}

	// 检查问题详情是否存在
	_, err = models.GetEvaluationQuestionDetailByID(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 404,
			"msg":  "评教问题详情不存在",
		})
		return
	}

	if err := models.DeleteEvaluationQuestionDetail(id); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "删除评教问题详情失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除评教问题详情成功",
	})
}
