package models

// Student 学生模型
type Student struct {
	BaseModel
	StudentName string `gorm:"column:student_name;type:varchar(100);not null;comment:学生姓名" json:"studentName"`
	StudentID   string `gorm:"column:student_id;type:varchar(50);uniqueIndex;not null;comment:学号" json:"studentID"`
	Phone       string `gorm:"column:phone;type:varchar(20);comment:手机号" json:"phone"`
	ClassName   string `gorm:"column:class_name;type:varchar(100);comment:班级名称" json:"className"`
	ClassID     uint64 `gorm:"column:class_id;type:bigint;not null;index;comment:班级ID" json:"classID"`
	OpenID      string `gorm:"column:open_id;type:varchar(100);comment:微信OpenID" json:"openID"`
	IsBound     bool   `gorm:"column:is_bound;type:boolean;default:false;comment:是否已绑定" json:"isBound"`
	Password    string `gorm:"column:password;type:varchar(100);not null;comment:密码" json:"-"`
}

// TableName 指定表名
func (Student) TableName() string {
	return "students"
}

// StudentQueryParams 学生查询参数
type StudentQueryParams struct {
	StudentName string   `json:"studentName" form:"studentName"`
	StudentID   string   `json:"studentID" form:"studentID"`
	ClassID     uint64   `json:"classID" form:"classID"`
	Page        int      `json:"page" form:"page"`
	PageSize    int      `json:"pageSize" form:"pageSize"`
	ClassIds    []uint64 `json:"classIds" form:"classIds"`
}

// GetStudentsTotal 获取学生总数
func GetStudentsTotal() (int64, error) {
	db := getDB()
	var count int64
	err := db.Model(&Student{}).Where("is_deleted = ?", false).Count(&count).Error
	return count, err
}

// CreateStudent 创建学生
func CreateStudent(student *Student) error {
	db := getDB()
	return db.Create(student).Error
}

// CreateBatchStudent 创建批量学生
func CreateBatchStudent(student []*Student) error {
	db := getDB()
	return db.CreateInBatches(student, len(student)).Error
}

// GetStudentByID 根据ID获取学生
func GetStudentByID(id uint64) (*Student, error) {
	db := getDB()
	var student Student
	err := db.First(&student, id).Error
	if err != nil {
		return nil, err
	}
	return &student, nil
}

// GetStudentByStudentID 根据学号获取学生
func GetStudentByStudentID(studentID string) (*Student, error) {
	db := getDB()
	var student Student
	err := db.Where("student_id = ? AND is_deleted = ?", studentID, false).First(&student).Error
	if err != nil {
		return nil, err
	}
	return &student, nil
}

// GetStudentsByClassID 根据班级ID获取学生列表
func GetStudentsByClassID(classID uint) ([]Student, error) {
	db := getDB()
	var students []Student
	err := db.Where("class_id = ? AND is_deleted = ?", classID, false).Order("id ASC").Find(&students).Error
	return students, err
}

// GetAllStudents 获取所有学生
func GetAllStudents() ([]Student, error) {
	db := getDB()
	var students []Student
	err := db.Where("is_deleted = ?", false).Order("id ASC").Find(&students).Error
	return students, err
}

// UpdateStudent 更新学生信息
func UpdateStudent(student *Student) error {
	db := getDB()
	return db.Save(student).Error
}

// DeleteStudent 删除学生（软删除）
func DeleteStudent(id uint64) error {
	db := getDB()
	return db.Model(&Student{}).Where("id = ?", id).Update("is_deleted", true).Error
}

// UpdateStudentPassword 更新学生密码
func UpdateStudentPassword(id uint64, password string) error {
	db := getDB()
	return db.Model(&Student{}).Where("id = ?", id).Update("password", password).Error
}

// UpdateStudentOpenID 更新学生微信OpenID
func UpdateStudentOpenID(id uint64, openID string, isBound bool) error {
	db := getDB()
	return db.Model(&Student{}).Where("id = ?", id).Updates(map[string]interface{}{
		"open_id":  openID,
		"is_bound": isBound,
	}).Error
}

// GetStudentByPhone 根据手机号获取学生
func GetStudentByPhone(phone string) (*Student, error) {
	db := getDB()
	var student Student
	err := db.Where("phone = ? AND is_deleted = ?", phone, false).First(&student).Error
	if err != nil {
		return nil, err
	}
	return &student, nil
}

// GetStudentByNameIdClass 根据姓名、学号、班级获取学生
func GetStudentByNameIdClass(studentName, studentID, className string) (*Student, error) {
	db := getDB()
	var student Student
	err := db.Where("student_name = ? AND student_id = ? AND class_name = ? AND is_deleted = ?",
		studentName, studentID, className, false).First(&student).Error
	if err != nil {
		return nil, err
	}
	return &student, nil
}

func GetCountStudentByNameIdClassId(studentName string, studentID string, classID uint64) int64 {
	db := getDB()
	var qty int64
	err := db.Where("student_name = ? AND student_id = ? and class_id=? AND is_deleted = ?",
		studentName, studentID, classID, false).Count(&qty).Error
	if err != nil {
		return -1
	}
	return qty
}

// GetStudentsWithPagination 分页获取学生列表（支持姓名、学号模糊查询，班级ID查询）
func GetStudentsWithPagination(params StudentQueryParams) ([]Student, int64, error) {
	db := getDB()
	var students []Student
	var total int64

	// 构建查询条件
	query := db.Model(&Student{}).Where("is_deleted = ?", false)

	// 添加查询条件
	if params.StudentName != "" {
		query = query.Where("student_name LIKE ?", "%"+params.StudentName+"%")
	}

	if params.StudentID != "" {
		query = query.Where("student_id LIKE ?", "%"+params.StudentID+"%")
	}

	if params.ClassID != 0 {
		query = query.Where("class_id = ?", params.ClassID)
	}
	if len(params.ClassIds) > 0 {
		query = query.Where("class_id IN ?", params.ClassIds)
	}

	// 获取总记录数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 设置默认分页参数
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}

	// 计算偏移量
	offset := (params.Page - 1) * params.PageSize

	// 执行查询
	if err := query.Offset(offset).Limit(params.PageSize).Order("id DESC").Find(&students).Error; err != nil {
		return nil, 0, err
	}

	return students, total, nil
}
