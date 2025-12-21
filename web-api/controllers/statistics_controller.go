package controllers

import (
	"encoding/csv"
	"log"
	"net/http"
	"strconv"
	"tems-web-api/models"

	"github.com/gin-gonic/gin"
)

// EvaluationDetailStatistics 评教详情统计结果结构体
type EvaluationDetailStatistics struct {
	TaskName        string `json:"task_name"`
	AcademicYear    string `json:"academic_year"`
	Semester        string `json:"semester"`
	ClassName       string `json:"class_name"`
	CourseName      string `json:"course_name"`
	UserName        string `json:"user_name"`
	QuestionSetName string `json:"question_set_name"`
	Title           string `json:"title"`
	Options         string `json:"options"`
	Qty             int    `json:"qty"`
}

// StatsPageRequest 统计分页请求参数
type StatsPageRequest struct {
	Limit    int    `json:"limit" form:"limit" default:"10"`
	Page     int    `json:"page" form:"page" default:"1"`
	TaskID   uint64 `json:"taskId" form:"taskId" default:"0"`
	ClassID  uint64 `json:"classId" form:"classId" default:"0"`
	CourseID uint64 `json:"courseId" form:"courseId" default:"0"`
	UserID   uint64 `json:"userId" form:"userId" default:"0"`
}

// GetEvaluationDetailStatistics 获取评教详情统计
func GetEvaluationDetailStatistics(c *gin.Context) {
	// 获取查询参数
	var req StatsPageRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "参数绑定失败",
			"error": err.Error(),
		})
		return
	}
	var args []interface{}
	whereClause := ""
	if req.ClassID > 0 {
		whereClause += " and d.class_id = ?"
		args = append(args, req.ClassID)
	}
	if req.CourseID > 0 {
		whereClause += " and d.course_id = ?"
		args = append(args, req.CourseID)
	}
	if req.UserID > 0 {
		whereClause += " and d.teacher_id = ?"
		args = append(args, req.UserID)
	}
	if req.TaskID > 0 {
		whereClause += " and d.task_id = ?"
		args = append(args, req.TaskID)
	}
	// 准备SQL查询
	sqlQuery := `
		SELECT t.task_name,t.academic_year,case t.semester when 0 then '上学期' when 1 then '下学期' end as semester,d.class_name,d.course_name,d.user_name,t.question_set_name,
		eqd.title,concat(opt->>'key','：',opt->>'value') as options,count(er.answer) as qty 
		FROM evaluation_task_details AS d 
		INNER JOIN evaluation_tasks t ON t.id = d.task_id 
		left join evaluation_question_details eqd on t.question_set_id =eqd.question_set_id 
		CROSS JOIN LATERAL jsonb_array_elements(eqd.options::jsonb) as opt 
		left join evaluation_results er on er.question_detail_id =eqd.id and er.answer = opt->>'key' 
		and er.task_id=d.task_id and er.task_detail_id=d.id and er.question_set_id=t.question_set_id 
		WHERE 1=1 ` + whereClause + `
		group by t.task_name,t.academic_year,t.semester,t.question_set_name,d.class_name,d.course_name,d.user_name,d.course_id,eqd.id,eqd.title,eqd.question_type,opt->>'key',opt->>'value',er.answer 
		ORDER BY d.course_id,eqd.id, opt->>'key' 
		limit ? offset ?
	`
	sqlCount := `
		select count(*) from (SELECT t.task_name,t.academic_year,case t.semester when 0 then '上学期' when 1 then '下学期' end as semester,d.class_name,d.course_name,d.user_name,t.question_set_name,
		eqd.title,concat(opt->>'key','：',opt->>'value') as options,count(er.answer) as qty 
		FROM evaluation_task_details AS d 
		INNER JOIN evaluation_tasks t ON t.id = d.task_id 
		left join evaluation_question_details eqd on t.question_set_id =eqd.question_set_id 
		CROSS JOIN LATERAL jsonb_array_elements(eqd.options::jsonb) as opt 
		left join evaluation_results er on er.question_detail_id =eqd.id and er.answer = opt->>'key' 
		and er.task_id=d.task_id and er.task_detail_id=d.id and er.question_set_id=t.question_set_id 
		WHERE 1=1 ` + whereClause + `
		group by t.task_name,t.academic_year,t.semester,t.question_set_name,d.class_name,d.course_name,d.user_name,d.course_id,eqd.id,eqd.title,eqd.question_type,opt->>'key',opt->>'value',er.answer) as TB
	`
	// 执行查询
	db := models.GetDB()
	var total int64
	if err := db.Raw(sqlCount, args...).Scan(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "查询总数失败",
			"error": err.Error(),
		})
		return
	}
	var results []EvaluationDetailStatistics
	queryArgs := append(args, req.Limit, (req.Page-1)*req.Limit)
	if err := db.Raw(sqlQuery, queryArgs...).Scan(&results).Error; err != nil {
		log.Printf("查询评教详情统计失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "查询评教详情统计失败",
			"error": err.Error(),
		})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取评教详情统计成功",
		"data": gin.H{
			"total": total,
			"list":  results,
		},
	})
}

// ExportEvaluationDetailStatistics 导出评教详情统计为CSV
func ExportEvaluationDetailStatistics(c *gin.Context) {
	// 获取查询参数
	var req StatsPageRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "参数绑定失败",
			"error": err.Error(),
		})
		return
	}
	var args []interface{}
	whereClause := ""
	if req.ClassID > 0 {
		whereClause += " and d.class_id = ?"
		args = append(args, req.ClassID)
	}
	if req.CourseID > 0 {
		whereClause += " and d.course_id = ?"
		args = append(args, req.CourseID)
	}
	if req.UserID > 0 {
		whereClause += " and d.teacher_id = ?"
		args = append(args, req.UserID)
	}
	if req.TaskID > 0 {
		whereClause += " and d.task_id = ?"
		args = append(args, req.TaskID)
	}
	// 准备SQL查询
	sqlQuery := `
		SELECT t.task_name,t.academic_year,case t.semester when 0 then '上学期' when 1 then '下学期' end as semester,d.class_name,d.course_name,d.user_name,t.question_set_name,
		eqd.title,concat(opt->>'key','：',opt->>'value') as options,count(er.answer) as qty 
		FROM evaluation_task_details AS d 
		INNER JOIN evaluation_tasks t ON t.id = d.task_id 
		left join evaluation_question_details eqd on t.question_set_id =eqd.question_set_id 
		CROSS JOIN LATERAL jsonb_array_elements(eqd.options::jsonb) as opt 
		left join evaluation_results er on er.question_detail_id =eqd.id and er.answer = opt->>'key' 
		and er.task_id=d.task_id and er.task_detail_id=d.id and er.question_set_id=t.question_set_id 
		WHERE 1=1 ` + whereClause + `
		group by t.task_name,t.academic_year,t.semester,t.question_set_name,d.class_name,d.course_name,d.user_name,d.course_id,eqd.id,eqd.title,eqd.question_type,opt->>'key',opt->>'value',er.answer 
		ORDER BY d.course_id,eqd.id, opt->>'key' 
	`
	// 执行查询
	db := models.GetDB()
	var results []EvaluationDetailStatistics
	if err := db.Raw(sqlQuery, args...).Scan(&results).Error; err != nil {
		log.Printf("查询评教详情统计失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "查询评教详情统计失败",
			"error": err.Error(),
		})
		return
	}

	// 设置HTTP头，触发文件下载
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=evaluation_statistics.csv")
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")

	// 使用csv包生成CSV文件
	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// 写入CSV表头
	header := []string{"任务名称", "学年", "学期", "班级名称", "课程名称", "教师姓名", "问题集名称", "问题标题", "选项", "选择数量"}
	if err := writer.Write(header); err != nil {
		log.Printf("写入CSV表头失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "生成CSV文件失败",
			"error": err.Error(),
		})
		return
	}

	// 写入数据行
	for _, result := range results {
		row := []string{
			result.TaskName,
			result.AcademicYear,
			result.Semester,
			result.ClassName,
			result.CourseName,
			result.UserName,
			result.QuestionSetName,
			result.Title,
			result.Options,
			strconv.Itoa(result.Qty),
		}
		if err := writer.Write(row); err != nil {
			log.Printf("写入CSV数据失败: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":  500,
				"msg":   "生成CSV文件失败",
				"error": err.Error(),
			})
			return
		}
	}

	// 刷新缓冲区，确保所有数据都写入响应
	writer.Flush()
	if err := writer.Error(); err != nil {
		log.Printf("刷新CSV写入器失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "生成CSV文件失败",
			"error": err.Error(),
		})
		return
	}
}

func HomeTopCardStatistics(c *gin.Context) {
	// 获取班级总数
	classesTotal, err := models.GetClassesTotal()
	if err != nil {
		classesTotal = 0
	}
	// 获取课程总数
	coursesTotal, err := models.GetCoursesTotal()
	if err != nil {
		coursesTotal = 0
	}
	// 获取学生总数
	studentsTotal, err := models.GetStudentsTotal()
	if err != nil {
		studentsTotal = 0
	}
	// 获取评教任务总数
	evaluationTasksTotal, err := models.GetEvaluationTasksTotal()
	if err != nil {
		evaluationTasksTotal = 0
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取评教详情统计成功",
		"data": gin.H{
			"classesTotal":  classesTotal,
			"coursesTotal":  coursesTotal,
			"studentsTotal": studentsTotal,
			"tasksTotal":    evaluationTasksTotal,
		},
	})
}
