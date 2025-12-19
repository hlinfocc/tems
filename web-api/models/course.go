package models

// Course 课程模型
type Course struct {
	BaseModel
	CourseName string `gorm:"column:course_name;type:varchar(255);not null;comment:课程名称" json:"courseName"`
	Status     int    `gorm:"column:status;type:int;default:0;comment:状态: 0启用 1禁用" json:"status"`
}

// TableName 指定表名
func (Course) TableName() string {
	return "courses"
}

// CreateCourse 创建课程
func CreateCourse(course *Course) error {
	db := getDB()
	return db.Create(course).Error
}

// GetCourseByID 根据ID获取课程
func GetCourseByID(id uint64) (*Course, error) {
	db := getDB()
	var course Course
	err := db.First(&course, id).Error
	if err != nil {
		return nil, err
	}
	return &course, nil
}

// GetActiveCourses 获取启用的课程
func GetActiveCourses() ([]Course, error) {
	db := getDB()
	var courses []Course
	err := db.Where("status = ? AND is_deleted = ?", 0, false).Order("id ASC").Find(&courses).Error
	return courses, err
}

// GetAllCourses 获取所有课程
func GetAllCourses() ([]Course, error) {
	db := getDB()
	var courses []Course
	err := db.Where("is_deleted = ?", false).Order("id ASC").Find(&courses).Error
	return courses, err
}

// UpdateCourse 更新课程
func UpdateCourse(course *Course) error {
	db := getDB()
	return db.Save(course).Error
}

// DeleteCourse 删除课程（软删除）
func DeleteCourse(id uint64) error {
	db := getDB()
	return db.Model(&Course{}).Where("id = ?", id).Update("is_deleted", true).Error
}

// FindCourse 根据条件查询课程
func FindCourse(condition map[string]interface{}) (*Course, error) {
	db := getDB()
	var course Course
	err := db.Where(condition).First(&course).Error
	if err != nil {
		return nil, err
	}
	return &course, nil
}

// GetCoursesTotal 获取课程总数
func GetCoursesTotal() (int64, error) {
	db := getDB()
	var count int64
	err := db.Model(&Course{}).Where("is_deleted = ?", false).Count(&count).Error
	return count, err
}

// QueryCoursesWithPagination 分页查询课程列表（支持课程名称模糊查询）
func QueryCoursesWithPagination(page, limit int, courseName string) ([]Course, int64, error) {
	db := getDB()
	var courses []Course
	var count int64

	// 计算偏移量
	offset := (page - 1) * limit

	// 构建查询条件
	query := db.Model(&Course{}).Where("is_deleted = ?", false)

	// 如果有课程名称，添加模糊匹配条件
	if courseName != "" {
		courseNamePattern := "%" + courseName + "%"
		query = query.Where("course_name LIKE ?", courseNamePattern)
	}

	// 查询总数
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	err := query.Order("id DESC").Offset(offset).Limit(limit).Find(&courses).Error
	if err != nil {
		return nil, 0, err
	}

	return courses, count, nil
}
