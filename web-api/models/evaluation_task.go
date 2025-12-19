package models

import (
	"time"
)

// EvaluationTask 评教任务模型
type EvaluationTask struct {
	BaseModel
	TaskName        string   `gorm:"column:task_name;type:varchar(255);not null;comment:任务名称" json:"taskName"`
	AcademicYear    string   `gorm:"column:academic_year;type:varchar(50);not null;comment:学年" json:"academicYear"`
	Semester        int      `gorm:"column:semester;type:int;not null;comment:学期: 0上学期 1下学期" json:"semester"`
	QuestionSetName string   `gorm:"column:question_set_name;type:varchar(255);not null;comment:问题集名称" json:"questionSetName"`
	QuestionSetID   uint64   `gorm:"column:question_set_id;type:bigint;not null;index;comment:问题集ID" json:"questionSetID"`
	StartTime       JsonTime `gorm:"column:start_time;type:timestamp;comment:开始时间" json:"startTime"`
	EndTime         JsonTime `gorm:"column:end_time;type:timestamp;comment:结束时间" json:"endTime"`
}

// TableName 指定表名
func (EvaluationTask) TableName() string {
	return "evaluation_tasks"
}

// CreateEvaluationTask 创建评教任务
func CreateEvaluationTask(task *EvaluationTask) error {
	db := getDB()
	return db.Create(task).Error
}

// GetEvaluationTaskByID 根据ID获取评教任务
func GetEvaluationTaskByID(id uint64) (*EvaluationTask, error) {
	db := getDB()
	var task EvaluationTask
	err := db.Where("id = ? AND is_deleted = ?", id, false).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// GetActiveTasks 获取当前有效的评教任务
func GetActiveTasks() ([]EvaluationTask, error) {
	db := getDB()
	now := time.Now()
	var tasks []EvaluationTask
	err := db.Where("start_time <= ? AND end_time >= ? AND is_deleted = ?", now, now, false).Order("createdAt DESC").Find(&tasks).Error
	return tasks, err
}

// GetAllEvaluationTasks 获取所有评教任务
func GetAllEvaluationTasks() ([]EvaluationTask, error) {
	db := getDB()
	var tasks []EvaluationTask
	err := db.Where("is_deleted = ?", false).Order("createdAt DESC").Find(&tasks).Error
	return tasks, err
}

// GetTasksByAcademicYear 根据学年获取评教任务
func GetTasksByAcademicYear(academicYear string) ([]EvaluationTask, error) {
	db := getDB()
	var tasks []EvaluationTask
	err := db.Where("academic_year = ? AND is_deleted = ?", academicYear, false).Order("createdAt DESC").Find(&tasks).Error
	return tasks, err
}

// UpdateEvaluationTask 更新评教任务
func UpdateEvaluationTask(task *EvaluationTask) error {
	db := getDB()
	// 只更新非零值字段
	return db.Model(&EvaluationTask{}).Where("id = ? AND is_deleted = ?", task.ID, false).Updates(task).Error
}

// DeleteEvaluationTask 删除评教任务（软删除）
func DeleteEvaluationTask(id uint64) error {
	db := getDB()
	return db.Model(&EvaluationTask{}).Where("id = ?", id).Update("is_deleted", true).Error
}

// QueryEvaluationTasksWithPagination 分页查询评教任务
func QueryEvaluationTasksWithPagination(page, limit int, keyword, academicYear string, semester *int) ([]EvaluationTask, int64, error) {
	db := getDB()
	var tasks []EvaluationTask
	var count int64

	query := db.Model(&EvaluationTask{}).Where("is_deleted = ?", false)

	// 添加搜索条件
	if keyword != "" {
		query = query.Where("task_name LIKE ?", "%"+keyword+"%")
	}

	if academicYear != "" {
		query = query.Where("academic_year = ?", academicYear)
	}

	if semester != nil {
		query = query.Where("semester = ?", *semester)
	}

	// 获取总数
	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * limit
	err = query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&tasks).Error
	if err != nil {
		return nil, 0, err
	}

	return tasks, count, nil
}

// GetTaskOptions 获取评教任务选项（用于下拉选择）
func GetTaskOptions() ([]map[string]interface{}, error) {
	db := getDB()
	var options []map[string]interface{}
	err := db.Model(&EvaluationTask{}).
		Select("id, task_name as name").
		Where("is_deleted = ?", false).
		Order("created_at DESC").
		Find(&options).Error
	return options, err
}

// IsTaskActive 检查任务是否处于活动状态
func IsTaskActive(id uint64) (bool, error) {
	db := getDB()
	now := time.Now()
	var count int64
	err := db.Model(&EvaluationTask{}).
		Where("id = ? AND is_deleted = ? AND start_time <= ? AND end_time >= ?", id, false, now, now).
		Count(&count).Error
	return count > 0, err
}

// GetTasksByTimeRange 根据时间范围获取评教任务
func GetTasksByTimeRange(start, end time.Time) ([]EvaluationTask, error) {
	db := getDB()
	var tasks []EvaluationTask
	err := db.Where("is_deleted = ? AND ((start_time >= ? AND start_time <= ?) OR (end_time >= ? AND end_time <= ?) OR (start_time <= ? AND end_time >= ?))", false, start, end, start, end, start, end).
		Order("start_time ASC").
		Find(&tasks).Error
	return tasks, err
}

func GetEvaluationTasksTotal() (int64, error) {
	db := getDB()
	var count int64
	err := db.Model(&EvaluationTask{}).Where("is_deleted = ?", false).Count(&count).Error
	return count, err
}
