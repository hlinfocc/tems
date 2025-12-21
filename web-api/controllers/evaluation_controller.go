package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"tems-web-api/models"
	"tems-web-api/utils"

	"github.com/gin-gonic/gin"
)

// EvaluationResultRequest 评教结果提交请求
type EvaluationResultRequest struct {
	TaskID           uint64 `json:"task_id" binding:"required"`
	QuestionSetID    uint64 `json:"question_set_id" binding:"required"`
	QuestionDetailID uint64 `json:"question_detail_id" binding:"required"`
	StudentID        uint64 `json:"student_id" binding:"required"`
	Result           int    `json:"result" binding:"required"`
}

// EnhancedEvaluationItem 增强版评课列表项
type EnhancedEvaluationItem struct {
	ID              uint64          `json:"id"`
	TaskID          uint64          `json:"taskId"`
	ClassName       string          `json:"className"`
	CourseName      string          `json:"courseName"`
	UserName        string          `json:"userName"`
	CreatedAt       models.JsonTime `json:"createdAt"`
	TaskName        string          `json:"taskName"`
	AcademicYear    string          `json:"academicYear"`
	Semester        int             `json:"semester"`
	QuestionSetName string          `json:"questionSetName"`
	QuestionSetID   uint64          `json:"questionSetID"`
	StartTime       models.JsonTime `json:"startTime"`
	EndTime         models.JsonTime `json:"endTime"`
	Result          *int            `json:"result"` // 使用指针表示可能为null
}

// GetEnhancedEvaluationListPage 增强版评课列表请求参数
type GetEnhancedEvaluationListPage struct {
	Result int `json:"result" binding:"omitempty"` // -1:不筛选 0:未评 1:正确 2:错误
	Page   int `json:"page" binding:"required"`
	Limit  int `json:"limit" binding:"required"`
}
type StatisticsResults struct {
	Total      int `json:"total"`
	Complete   int `json:"complete"`
	Incomplete int `json:"incomplete"`
}

// GetEvaluationTasks 获取评教任务列表
func GetEvaluationTasks(c *gin.Context) {
	tasks, err := models.GetActiveTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "获取评教任务失败",
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

// GetQuestionsBySetID 根据问题集ID获取问题列表
func GetQuestionsBySetID(c *gin.Context) {
	setIDStr := c.Param("setId")
	setID, err := strconv.ParseUint(setIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的问题集ID",
		})
		return
	}
	log.Printf("GetQuestionsBySetID: setID=%d", setID)
	questions, err := models.GetQuestionsBySetID(uint64(setID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "获取问题失败",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": questions,
	})
}

// SubmitEvaluationResult 提交评教结果（批量）
func SubmitEvaluationResult(c *gin.Context) {
	var results []models.EvaluationResult
	if err := c.ShouldBindJSON(&results); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "请求参数无效",
			"error": err.Error(),
		})
		return
	}

	// 验证结果数组不为空
	if len(results) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "评教结果数组不能为空",
		})
		return
	}

	// 批量验证结果值
	for i := range results {
		if results[i].Answer == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":  400,
				"msg":   "无效的结果值，答案不能为空",
				"error": "第" + strconv.Itoa(i+1) + "个评教结果的答案为空",
			})
			return
		}
	}

	// 使用批量保存方法
	if err := models.BatchCreatesEvaluationResults(results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "提交评教结果失败",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "评教结果提交成功",
		"data": gin.H{
			"count": len(results),
		},
	})
}

// GetEnhancedEvaluationList 获取增强版评课列表
func GetEnhancedEvaluationList(c *gin.Context) {
	// 获取查询参数
	var req GetEnhancedEvaluationListPage
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "请求参数无效",
			"error": err.Error(),
		})
		return
	}

	customClaimsStr, exists := c.Get("customClaims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "未授权",
		})
		return
	}

	// 使用安全的类型断言，注意middleware中设置的是值类型而非指针类型
	customClaims, ok := customClaimsStr.(utils.CustomClaims)
	if !ok {
		log.Printf("customClaimsStr类型错误，实际类型: %T, 值: %v", customClaimsStr, customClaimsStr)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "无效的身份验证信息格式",
		})
		return
	}

	studentID := customClaims.UserID
	classID := customClaims.ClassId

	// 参数验证
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 20 // 默认每页20条
	}
	// Result值验证：-1表示不筛选，0表示未评(空或0)，1表示正确，2表示错误
	if req.Result < -1 || req.Result > 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的结果筛选值，必须是-1(不筛选)、0(未评)、1(正确)或2(错误)",
		})
		return
	}

	whereClause := ""
	// 添加result筛选条件
	// Result = -1: 不筛选评教状态
	// Result = 0: 筛选未评（result为空或0）
	// Result = 1: 筛选已评（result为1）
	if req.Result == 0 {
		// 未评：result为空或0
		whereClause = " AND temp.results = 0"
	} else if req.Result > 0 {
		// 已评：result为1或2
		whereClause = " AND temp.results = 1"
	}

	// 准备主查询SQL
	sqlQuery := `
		select * from (
		SELECT 
			d.id,
			d.task_id,
			d.class_name,
			d.course_name,
			d.user_name,
			d.created_at,
			t.task_name,
			t.academic_year,
			t.semester,
			t.question_set_name,
			t.question_set_id,
			t.start_time,
			t.end_time,
			CASE 
				WHEN EXISTS (
					SELECT 1 FROM evaluation_results r 
					WHERE r.task_id = d.task_id
						AND r.task_detail_id = d.id
						AND r.question_set_id = t.question_set_id
						AND r.semester = t.semester
						AND r.academic_year = t.academic_year
						AND r.student_id = ?
				) THEN 1 
				ELSE 0 
			END as results 
		FROM evaluation_task_details AS d
		INNER JOIN evaluation_tasks t ON t.id = d.task_id AND CURRENT_TIMESTAMP BETWEEN t.start_time AND t.end_time
		where d.is_deleted=false and d.class_id IN (?)
		) as temp where 1=1 ` + whereClause + `
		ORDER BY temp.created_at desc 
		LIMIT ? OFFSET ?
	`

	// 准备总数查询SQL
	totalQuery := `
		select count(*) from (
		SELECT 
			d.id,
			d.task_id,
			d.class_name,
			d.course_name,
			d.user_name,
			d.created_at,
			t.task_name,
			t.academic_year,
			t.semester,
			t.question_set_name,
			t.question_set_id,
			t.start_time,
			t.end_time,
			CASE 
				WHEN EXISTS (
					SELECT 1 FROM evaluation_results r 
					WHERE r.task_id = d.task_id
						AND r.task_detail_id = d.id
						AND r.question_set_id = t.question_set_id
						AND r.semester = t.semester
						AND r.academic_year = t.academic_year
						AND r.student_id = ?
				) THEN 1 
				ELSE 0 
			END as results 
		FROM evaluation_task_details AS d
		INNER JOIN evaluation_tasks t ON t.id = d.task_id AND CURRENT_TIMESTAMP BETWEEN t.start_time AND t.end_time
		where d.is_deleted=false and d.class_id IN (?)
		) as temp where 1=1 ` + whereClause + `
	`

	// 添加查询参数
	// 准备SQL查询条件
	var args []interface{}
	args = append(args, classID)
	totalArgs := append([]interface{}{studentID}, args...)
	queryArgs := append([]interface{}{studentID}, args...)
	// 计算偏移量
	offset := (req.Page - 1) * req.Limit
	queryArgs = append(queryArgs, req.Limit, offset)

	// 获取总数
	db := models.GetDB()
	var total int64
	if err := db.Raw(totalQuery, totalArgs...).Scan(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "查询总数失败",
			"error": err.Error(),
		})
		return
	}

	// 执行查询
	rows, err := db.Raw(sqlQuery, queryArgs...).Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "查询评课列表失败",
			"error": err.Error(),
		})
		return
	}
	defer rows.Close()

	// 处理查询结果
	var evaluationList []EnhancedEvaluationItem
	for rows.Next() {
		var item EnhancedEvaluationItem
		var result sql.NullInt64

		err := rows.Scan(
			&item.ID,
			&item.TaskID,
			&item.ClassName,
			&item.CourseName,
			&item.UserName,
			&item.CreatedAt,
			&item.TaskName,
			&item.AcademicYear,
			&item.Semester,
			&item.QuestionSetName,
			&item.QuestionSetID,
			&item.StartTime,
			&item.EndTime,
			&result,
		)
		if err != nil {
			log.Printf("扫描查询结果失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":  500,
				"msg":   "处理查询结果失败",
				"error": err.Error(),
			})
			return
		}

		// 处理可能为null的result字段
		if result.Valid {
			resultValue := int(result.Int64)
			item.Result = &resultValue
		} else {
			item.Result = nil
		}

		evaluationList = append(evaluationList, item)
	}

	// 检查迭代过程中的错误
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "迭代查询结果失败",
			"error": err.Error(),
		})
		return
	}

	// 计算总页数
	totalPages := int(total) / req.Limit
	if int(total)%req.Limit > 0 {
		totalPages++
	}

	// 返回结果（包含分页信息）
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取评课列表成功",
		"data": gin.H{
			"list":       evaluationList,
			"total":      total,
			"page":       req.Page,
			"limit":      req.Limit,
			"totalPages": totalPages,
		},
	})
}

// GetStatisticsResults 获取评教任务统计结果
func GetStatisticsResults(c *gin.Context) {

	customClaimsStr, exists := c.Get("customClaims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "未授权",
		})
		return
	}

	// 使用安全的类型断言，注意middleware中设置的是值类型而非指针类型
	customClaims, ok := customClaimsStr.(utils.CustomClaims)
	if !ok {
		log.Printf("customClaimsStr类型错误，实际类型: %T, 值: %v", customClaimsStr, customClaimsStr)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "无效的身份验证信息格式",
		})
		return
	}

	studentID := customClaims.UserID
	classID := customClaims.ClassId

	// 准备主查询SQL
	sqlQuery := `
		select
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE results >= eqds) as complete,
			COUNT(*) FILTER (WHERE results < eqds) as incomplete
		from (
		SELECT d.id,d.course_name,t.question_set_id,
		(
		select count(*) as qty from evaluation_results r where r.task_id=d.task_id
		and r.task_detail_id =d.id
		and r.question_set_id =t.question_set_id
		and r.semester =t.semester
		and r.academic_year=t.academic_year
		AND r.student_id = ?
		) as results,
		(select count(*) from evaluation_question_details q where q.question_set_id=t.question_set_id) as eqds
		FROM evaluation_task_details AS d
		INNER JOIN evaluation_tasks t ON t.id = d.task_id
			AND CURRENT_TIMESTAMP BETWEEN t.start_time AND t.end_time
		WHERE d.is_deleted=false and d.class_id IN (?)
		) as TB
	`

	// 添加查询参数
	queryArgs := []interface{}{studentID, classID}

	// 获取总数
	db := models.GetDB()
	var statRes StatisticsResults
	if err := db.Raw(sqlQuery, queryArgs...).Scan(&statRes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "查询统计结果失败",
			"error": err.Error(),
		})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取评教任务统计结果成功",
		"data": statRes,
	})
}
