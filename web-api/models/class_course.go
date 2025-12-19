package models

// ClassCourse 班级课程模型
type ClassCourse struct {
	BaseModel
	ClassID     uint64 `gorm:"column:class_id;type:bigint;not null;index;comment:班级ID" json:"class_id"`
	CourseID    uint64 `gorm:"column:course_id;type:bigint;not null;index;comment:课程ID" json:"course_id"`
	TeacherID   uint64 `gorm:"column:teacher_id;type:bigint;not null;comment:教师ID" json:"teacher_id"`
	CourseName  string `gorm:"column:course_name;type:varchar(255);comment:课程名称" json:"courseName"`
	TeacherName string `gorm:"column:teacher_name;type:varchar(255);comment:任课教师名称" json:"teacherName"`
}

// TableName 指定表名
func (ClassCourse) TableName() string {
	return "class_courses"
}

// CreateClassCourse 创建班级课程关系
func CreateClassCourse(classCourse *ClassCourse) error {
	db := getDB()
	return db.Model(&ClassCourse{}).Create(classCourse).Error
}

// GetClassCourseByID 根据ID获取班级课程关系
func GetClassCourseByID(id uint) (*ClassCourse, error) {
	db := getDB()
	var classCourse ClassCourse
	err := db.Model(&ClassCourse{}).Where("id = ? AND is_deleted = ?", id, false).First(&classCourse).Error
	if err != nil {
		return nil, err
	}
	return &classCourse, nil
}

// GetCoursesByClassID 根据班级ID获取课程
func GetCoursesByClassID(classID uint) ([]ClassCourse, error) {
	db := getDB()
	var classCourses []ClassCourse
	err := db.Model(&ClassCourse{}).Where("class_id = ? AND is_deleted = ?", classID, false).Find(&classCourses).Error
	return classCourses, err
}

// GetClassesByCourseID 根据课程ID获取班级
func GetClassesByCourseID(courseID uint) ([]ClassCourse, error) {
	db := getDB()
	var classCourses []ClassCourse
	err := db.Model(&ClassCourse{}).Where("course_id = ? AND is_deleted = ?", courseID, false).Find(&classCourses).Error
	return classCourses, err
}

func GetCourseByClassCourseID(classID uint64,courseID uint64) (bool, error) {
	db := getDB()
	var total int64
	err := db.Model(&ClassCourse{}).Where("class_id = ? AND course_id = ? AND is_deleted = ?", classID,courseID, false).Count(&total).Error
	return total > 0, err
}

// UpdateClassCourse 更新班级课程关系
func UpdateClassCourse(classCourse *ClassCourse) error {
	db := getDB()
	return db.Model(&ClassCourse{}).Save(classCourse).Error
}

// DeleteClassCourse 删除班级课程关系（软删除）
func DeleteClassCourse(id uint) error {
	db := getDB()
	return db.Model(&ClassCourse{}).Where("id = ?", id).Update("is_deleted", true).Error
}

// DeleteClassCoursesByClassID 根据班级ID批量删除班级课程关系
func DeleteClassCoursesByClassID(classID uint) error {
	db := getDB()
	return db.Model(&ClassCourse{}).Where("class_id = ?", classID).Update("is_deleted", true).Error
}

// QueryClassCoursesWithPagination 分页查询班级课程关系
func QueryClassCoursesWithPagination(page, limit int, classID, courseID uint64, keyword string) ([]ClassCourse, int64, error) {
	db := getDB()
	var classCourses []ClassCourse
	var count int64

	// 计算偏移量
	offset := (page - 1) * limit

	// 构建查询条件
	query := db.Model(&ClassCourse{}).Where("is_deleted = ?", false)

	// 添加班级ID筛选
	if classID > 0 {
		query = query.Where("class_id = ?", classID)
	}

	// 添加课程ID筛选
	if courseID > 0 {
		query = query.Where("course_id = ?", courseID)
	}

	// 添加关键字搜索
	if keyword != "" {
		keyword = "%" + keyword + "%"
		query = query.Where("course_name LIKE ? OR teacher_name LIKE ?", keyword, keyword)
	}

	// 计算总数
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&classCourses).Error; err != nil {
		return nil, 0, err
	}

	return classCourses, count, nil
}
