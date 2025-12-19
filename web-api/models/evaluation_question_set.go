package models

// EvaluationQuestionSet 评教问题集模型
type EvaluationQuestionSet struct {
	BaseModel
	Name   string `gorm:"column:name;type:varchar(255);not null;comment:问题集名称" json:"name"`
	Status int    `gorm:"column:status;type:integer;default:0;comment:状态(0-禁用,1-启用)" json:"status"`
	Year   int    `gorm:"column:year;type:integer;default:0;comment:年份" json:"year"`
	Remark string `gorm:"column:remark;type:text;comment:备注" json:"remark"`
}

// TableName 指定表名
func (EvaluationQuestionSet) TableName() string {
	return "evaluation_question_sets"
}

// CreateEvaluationQuestionSet 创建评教问题集
func CreateEvaluationQuestionSet(set *EvaluationQuestionSet) error {
	db := getDB()
	return db.Create(set).Error
}

// GetAllQuestionSets 获取所有问题集
func GetAllQuestionSets() ([]EvaluationQuestionSet, error) {
	db := getDB()
	var sets []EvaluationQuestionSet
	err := db.Where("is_deleted = ?", false).Order("createdAt DESC").Find(&sets).Error
	return sets, err
}

// GetEvaluationQuestionSetByID 根据ID获取评教问题集
func GetEvaluationQuestionSetByID(id uint64) (*EvaluationQuestionSet, error) {
	db := getDB()
	var set EvaluationQuestionSet
	err := db.Where("id = ? AND is_deleted = ?", id, false).First(&set).Error
	if err != nil {
		return nil, err
	}
	return &set, nil
}

// UpdateEvaluationQuestionSet 更新评教问题集
func UpdateEvaluationQuestionSet(set *EvaluationQuestionSet) error {
	db := getDB()
	return db.Save(set).Error
}

// DeleteEvaluationQuestionSet 删除评教问题集（软删除）
func DeleteEvaluationQuestionSet(id uint64) error {
	db := getDB()
	return db.Model(&EvaluationQuestionSet{}).Where("id = ?", id).Update("is_deleted", true).Error
}
