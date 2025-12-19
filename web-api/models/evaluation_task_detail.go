package models

// EvaluationTaskDetail 评教任务详情模型
type EvaluationTaskDetail struct {
	BaseModel
	TaskID     uint64 `gorm:"column:task_id;type:bigint;not null;index;comment:任务ID" json:"taskId"`
	ClassID    uint64 `gorm:"column:class_id;type:bigint;not null;index;comment:班级ID" json:"classId"`
	CourseID   uint64 `gorm:"column:course_id;type:bigint;not null;index;comment:课程ID" json:"courseId"`
	TeacherID  uint64 `gorm:"column:teacher_id;type:bigint;index;comment:任课教师ID" json:"teacherId"`
	ClassName  string `gorm:"column:class_name;type:varchar(255);comment:班级名称" json:"className"`
	CourseName string `gorm:"column:course_name;type:varchar(255);comment:课程名称" json:"courseName"`
	UserName   string `gorm:"column:user_name;type:varchar(255);comment:任课教师名称" json:"userName"`
}

// TableName 指定表名
func (EvaluationTaskDetail) TableName() string {
	return "evaluation_task_details"
}

// CreateEvaluationTaskDetail 创建评教任务详情
func CreateEvaluationTaskDetail(detail *EvaluationTaskDetail) error {
	db := getDB()
	return db.Create(detail).Error
}

// CreateBatchEvaluationTaskDetail 创建批量评教任务详情
func CreateBatchEvaluationTaskDetail(details []*EvaluationTaskDetail) error {
	db := getDB()
	return db.CreateInBatches(details, len(details)).Error
}

// GetEvaluationTaskDetailByID 根据ID获取评教任务详情
func GetEvaluationTaskDetailByID(id uint) (*EvaluationTaskDetail, error) {
	db := getDB()
	var detail EvaluationTaskDetail
	err := db.Where("id = ? AND is_deleted = ?", id, false).First(&detail).Error
	if err != nil {
		return nil, err
	}
	return &detail, nil
}

// GetTaskDetailsByTaskID 根据任务ID获取详情列表
func GetTaskDetailsByTaskID(taskID uint) ([]EvaluationTaskDetail, error) {
	db := getDB()
	var details []EvaluationTaskDetail
	err := db.Where("task_id = ? AND is_deleted = ?", taskID, false).Find(&details).Error
	return details, err
}

// GetTaskDetailsByClassID 根据班级ID获取详情列表
func GetTaskDetailsByClassID(classID uint) ([]EvaluationTaskDetail, error) {
	db := getDB()
	var details []EvaluationTaskDetail
	err := db.Where("class_id = ? AND is_deleted = ?", classID, false).Find(&details).Error
	return details, err
}

// GetTaskDetailsByTaskAndClassID 根据任务ID和班级ID获取详情列表
func GetTaskDetailsByTaskAndClassID(taskID, classID uint) ([]EvaluationTaskDetail, error) {
	db := getDB()
	var details []EvaluationTaskDetail
	err := db.Where("task_id = ? AND class_id = ? AND is_deleted = ?", taskID, classID, false).Find(&details).Error
	return details, err
}

// QueryEvaluationTaskDetailsWithPagination 分页查询评教任务详情
func QueryEvaluationTaskDetailsWithPagination(page, limit int, taskID, classID uint) ([]EvaluationTaskDetail, int64, error) {
	db := getDB()
	var details []EvaluationTaskDetail
	var count int64

	// 构建查询
	query := db.Model(&EvaluationTaskDetail{}).Where("is_deleted = ?", false)

	// 添加查询条件
	if taskID > 0 {
		query = query.Where("task_id = ?", taskID)
	}
	if classID > 0 {
		query = query.Where("class_id = ?", classID)
	}

	// 获取总数
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&details).Error; err != nil {
		return nil, 0, err
	}

	return details, count, nil
}

// GetTaskDetailByTaskClassCourse 根据任务ID、班级ID和课程ID获取详情
func GetTaskDetailByTaskClassCourse(taskID, classID, courseID uint) (*EvaluationTaskDetail, error) {
	db := getDB()
	var detail EvaluationTaskDetail
	err := db.Where("task_id = ? AND class_id = ? AND course_id = ? AND is_deleted = ?", taskID, classID, courseID, false).First(&detail).Error
	if err != nil {
		return nil, err
	}
	return &detail, nil
}

// UpdateEvaluationTaskDetail 更新评教任务详情
func UpdateEvaluationTaskDetail(detail *EvaluationTaskDetail) error {
	db := getDB()
	return db.Save(detail).Error
}

// DeleteEvaluationTaskDetail 删除评教任务详情（软删除）
func DeleteEvaluationTaskDetail(id uint) error {
	db := getDB()
	return db.Model(&EvaluationTaskDetail{}).Where("id = ?", id).Update("is_deleted", true).Error
}

// DeleteTaskDetailsByTaskID 根据任务ID批量删除详情
func DeleteTaskDetailsByTaskID(taskID uint) error {
	db := getDB()
	return db.Model(&EvaluationTaskDetail{}).Where("task_id = ?", taskID).Update("is_deleted", true).Error
}
