package models

// EvaluationResult 评教结果模型
type EvaluationResult struct {
	BaseModel
	TaskID           uint64 `gorm:"column:task_id;type:bigint;not null;index;comment:任务ID" json:"taskId"`
	TaskDetailID     uint64 `gorm:"column:task_detail_id;type:bigint;not null;index;comment:任务详情ID" json:"taskDetailId"`
	QuestionSetID    uint64 `gorm:"column:question_set_id;type:bigint;not null;index;comment:问题集ID" json:"questionSetId"`
	QuestionDetailID uint64 `gorm:"column:question_detail_id;type:bigint;not null;index;comment:问题详情ID" json:"questionDetailId"`
	StudentID        uint64 `gorm:"column:student_id;type:bigint;not null;index;comment:学生ID" json:"studentId"`
	QuestionType     int    `gorm:"column:question_type;type:integer;not null;comment:问题类型: 单选、多选" json:"questionType"`
	Semester         int `gorm:"column:semester;type:integer;not null;comment:学期" json:"semester"`
	AcademicYear     string `gorm:"column:academic_year;type:varchar(20);not null;comment:学年" json:"academicYear"`
	Answer           string `gorm:"column:answer;type:text;comment:答案" json:"answer"`
	Result           int    `gorm:"column:result;type:integer;default:0;comment:结果: 0未评 1正确 2错误" json:"result"`
}

// TableName 指定表名
func (EvaluationResult) TableName() string {
	return "evaluation_results"
}

// CreateEvaluationResult 创建评估结果
func CreateEvaluationResult(result *EvaluationResult) error {
	db := getDB()
	return db.Create(result).Error
}

// UpdateEvaluationResult 更新评估结果
func UpdateEvaluationResult(result *EvaluationResult) error {
	db := getDB()
	return db.Save(result).Error
}

// CreateOrUpdateEvaluationResult 创建或更新评教结果
func CreateOrUpdateEvaluationResult(result *EvaluationResult) error {
	db := getDB()
	// 使用Upsert操作，通过task_id, student_id, question_detail_id进行唯一性约束
	return db.Where("task_id = ? AND student_id = ? AND question_detail_id = ?",
		result.TaskID, result.StudentID, result.QuestionDetailID).
		Assign(EvaluationResult{Result: result.Result}).
		FirstOrCreate(result).Error
}

// GetEvaluationResultByID 根据ID获取评估结果
func GetEvaluationResultByID(id uint) (*EvaluationResult, error) {
	db := getDB()
	var result EvaluationResult
	err := db.Where("id = ? AND is_deleted = ?", id, false).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetEvaluationResultsByTaskID 根据任务ID获取评估结果列表
func GetEvaluationResultsByTaskID(taskID uint) ([]EvaluationResult, error) {
	db := getDB()
	var results []EvaluationResult
	err := db.Where("task_id = ? AND is_deleted = ?", taskID, false).Find(&results).Error
	return results, err
}

// GetEvaluationResultsByStudentID 根据学生ID获取评估结果列表
func GetEvaluationResultsByStudentID(studentID uint) ([]EvaluationResult, error) {
	db := getDB()
	var results []EvaluationResult
	err := db.Where("student_id = ? AND is_deleted = ?", studentID, false).Find(&results).Error
	return results, err
}

// GetStudentResultsByTaskID 根据任务ID和学生ID获取评教结果
func GetStudentResultsByTaskID(taskID, studentID uint) ([]EvaluationResult, error) {
	db := getDB()
	var results []EvaluationResult
	err := db.Where("task_id = ? AND student_id = ? AND is_deleted = ?", taskID, studentID, false).Find(&results).Error
	return results, err
}

// DeleteEvaluationResult 删除评估结果（软删除）
func DeleteEvaluationResult(id uint) error {
	db := getDB()
	return db.Model(&EvaluationResult{}).Where("id = ?", id).Update("is_deleted", true).Error
}

// BatchCreatesEvaluationResults 批量插入评教结果
func BatchCreatesEvaluationResults(results []EvaluationResult) error {
	db := getDB()
	// 使用CreateInBatches批量插入记录（每批100条）
	return db.CreateInBatches(results, len(results)).Error
}
