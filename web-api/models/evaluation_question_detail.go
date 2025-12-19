package models

// EvaluationQuestionDetail 评教问题详情模型
type EvaluationQuestionDetail struct {
	BaseModel
	QuestionSetID uint64 `gorm:"column:question_set_id;type:bigint;not null;index;comment:问题集ID" json:"questionSetId"`
	Title         string `gorm:"column:title;type:varchar(500);not null;comment:问题标题" json:"title"`
	QuestionType  int    `gorm:"column:question_type;type:integer;not null;comment:问题类型: 单选、多选" json:"questionType"`
	Options       string `gorm:"column:options;type:text;comment:选项，JSON格式" json:"options"`
	Remark        string `gorm:"column:remark;type:text;comment:备注" json:"remark"`
}

// TableName 指定表名
func (EvaluationQuestionDetail) TableName() string {
	return "evaluation_question_details"
}

// CreateEvaluationQuestionDetail 创建评教问题详情
func CreateEvaluationQuestionDetail(detail *EvaluationQuestionDetail) error {
	db := getDB()
	return db.Create(detail).Error
}

// GetEvaluationQuestionDetailByID 根据ID获取评教问题详情
func GetEvaluationQuestionDetailByID(id uint64) (*EvaluationQuestionDetail, error) {
	db := getDB()
	var detail EvaluationQuestionDetail
	err := db.Where("id = ? AND is_deleted = ?", id, false).First(&detail).Error
	if err != nil {
		return nil, err
	}
	return &detail, nil
}

// GetQuestionsBySetID 根据问题集ID获取问题列表
func GetQuestionsBySetID(setID uint64) ([]EvaluationQuestionDetail, error) {
	db := getDB()
	var details []EvaluationQuestionDetail
	err := db.Where("question_set_id = ? AND is_deleted = ?", setID, false).Order("id ASC").Find(&details).Error
	return details, err
}

// UpdateEvaluationQuestionDetail 更新评教问题详情
func UpdateEvaluationQuestionDetail(detail *EvaluationQuestionDetail) error {
	db := getDB()
	return db.Save(detail).Error
}

// DeleteEvaluationQuestionDetail 删除评教问题详情（软删除）
func DeleteEvaluationQuestionDetail(id uint64) error {
	db := getDB()
	return db.Model(&EvaluationQuestionDetail{}).Where("id = ?", id).Update("is_deleted", true).Error
}

// DeleteQuestionsBySetID 根据问题集ID批量删除问题详情（软删除）
func DeleteQuestionsBySetID(setID uint64) error {
	db := getDB()
	return db.Model(&EvaluationQuestionDetail{}).Where("question_set_id = ?", setID).Update("is_deleted", true).Error
}
