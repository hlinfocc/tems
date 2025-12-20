package models

import "strings"

// Class 班级模型
type Class struct {
	BaseModel
	ClassName       string `gorm:"column:class_name;type:varchar(100);not null;comment:班级名称" json:"className"`
	Grade           string `gorm:"column:grade;type:varchar(50);comment:年级" json:"grade"`
	MajorName       string `gorm:"column:major_name;type:varchar(100);comment:专业名称" json:"majorName"`
	MajorID         uint64 `gorm:"column:major_id;type:bigint;comment:专业ID" json:"majorID"`
	HeadTeacherName string `gorm:"column:head_teacher_name;type:varchar(100);comment:班主任姓名" json:"headTeacherName"`
	HeadTeacherID   uint64 `gorm:"column:head_teacher_id;type:bigint;comment:班主任ID" json:"headTeacherId"`
}

// TableName 指定表名
func (Class) TableName() string {
	return "classes"
}

// CreateClass 创建班级
func CreateClass(class *Class) error {
	db := getDB()
	return db.Create(class).Error
}

// GetClassByID 根据ID获取班级
func GetClassByID(id uint64) (*Class, error) {
	db := getDB()
	var class Class
	err := db.First(&class, id).Error
	if err != nil {
		return nil, err
	}
	return &class, nil
}

// GetAllClasses 获取所有班级
func GetAllClasses() ([]Class, error) {
	db := getDB()
	var classes []Class
	err := db.Where("is_deleted = ?", false).Order("id ASC").Find(&classes).Error
	return classes, err
}

// UpdateClass 更新班级信息
func UpdateClass(class *Class) error {
	db := getDB()
	return db.Save(class).Error
}

// DeleteClass 删除班级（软删除）
func DeleteClass(id uint64) error {
	db := getDB()
	return db.Model(&Class{}).Where("id = ?", id).Update("is_deleted", true).Error
}

// GetClassesByGrade 根据年级获取班级
func GetClassesByGrade(grade string) ([]Class, error) {
	db := getDB()
	var classes []Class
	err := db.Model(&Class{}).Where("grade = ? AND is_deleted = ?", grade, false).Order("id ASC").Find(&classes).Error
	return classes, err
}

// GetClassesByName 根据班级名称获取班级
func GetClassesByName(className string) (Class, error) {
	db := getDB()
	var classes Class
	className = strings.TrimSpace(className)
	err := db.Model(&Class{}).Where("class_name = ? AND is_deleted = ?", className, false).First(&classes).Error
	return classes, err
}

// QueryClassesWithPagination 分页查询班级列表
func QueryClassesWithPagination(page, limit int, keyword string) ([]Class, int64, error) {
	db := getDB()
	var classes []Class
	var count int64

	// 计算偏移量
	offset := (page - 1) * limit

	// 构建查询条件
	query := db.Model(&Class{}).Where("is_deleted = ?", false)

	// 如果有关键字，添加模糊匹配条件
	if keyword != "" {
		keywordPattern := "%" + keyword + "%"
		query = query.Where("class_name LIKE ? OR major_name LIKE ? OR head_teacher_name LIKE ?", keywordPattern, keywordPattern, keywordPattern)
	}

	// 查询总数
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	err := query.Order("id DESC").Offset(offset).Limit(limit).Find(&classes).Error
	if err != nil {
		return nil, 0, err
	}

	return classes, count, nil
}

func GetClassesTotal() (int64, error) {
	db := getDB()
	var count int64
	err := db.Model(&Class{}).Where("is_deleted = ?", false).Count(&count).Error
	return count, err
}
