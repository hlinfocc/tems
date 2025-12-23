package models

// AdminUser 管理员用户模型
type AdminUser struct {
	BaseModel
	Name     string `gorm:"column:name;type:varchar(100);not null;comment:管理员姓名" json:"name"`
	Username string `gorm:"column:username;type:varchar(50);uniqueIndex;not null;comment:用户名" json:"username"`
	Password string `gorm:"column:password;type:text;not null;comment:密码" json:"-"`
	Status   int    `gorm:"column:status;type:integer;default:0;comment:状态:0启用1禁用" json:"status"`           // 0 启用 1禁用
	UserType int    `gorm:"column:user_type;type:integer;default:0;comment:用户类型:0管理员1班主任" json:"user_type"` // 0 管理员 1 老师
}

// TableName 指定表名
func (AdminUser) TableName() string {
	return "admin_user"
}

// AdminUserRequestParams 管理员用户请求参数
type AdminUserRequestParams struct {
	Page     *int   `form:"page" json:"page"`
	Limit    *int   `form:"limit" json:"limit"`
	IsPage   *int   `form:"isPage" json:"isPage"`
	Status   *int   `form:"status" json:"status"`
	UserType *int   `form:"user_type" json:"user_type"`
	Keyword  string `form:"keyword" json:"keyword"`
	Account  string `form:"account" json:"account"`
}

// CreateAdminUser 创建管理员用户
func CreateAdminUser(adminUser *AdminUser) error {
	db := getDB()
	return db.Create(adminUser).Error
}

// GetAdminUserByID 根据ID获取管理员用户
func GetAdminUserByID(id uint64) (*AdminUser, error) {
	db := getDB()
	var adminUser AdminUser
	err := db.First(&adminUser, id).Error
	if err != nil {
		return nil, err
	}
	return &adminUser, nil
}

// GetAdminUserByUsername 根据用户名获取管理员用户
func GetAdminUserByUsername(username string) (*AdminUser, error) {
	db := getDB()
	var adminUser AdminUser
	err := db.Where("username = ? AND is_deleted = ?", username, false).First(&adminUser).Error
	if err != nil {
		return nil, err
	}
	return &adminUser, nil
}

// GetAllAdminUsers 获取所有管理员用户
func GetAllAdminUsers(params *AdminUserRequestParams) ([]AdminUser, int64, error) {
	db := getDB()
	var adminUsers []AdminUser
	query := db.Model(&AdminUser{}).Where("is_deleted = ?", false).Order("id ASC")
	if params.Keyword != "" {
		query = query.Where("name LIKE ? OR username LIKE ?", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}
	if params.Account != "" {
		query = query.Where("username like ?", "%"+params.Account+"%")
	}
	// 状态查询
	if params.Status != nil {
		query = query.Where("status = ?", *params.Status)
	}
	// 用户类型查询
	if params.UserType != nil {
		query = query.Where("user_type = ?", *params.UserType)
	}
	// 获取总数
	var count int64
	if err := query.Count(&count).Error; err != nil {

		return nil, 0, err
	}
	// 分页查询
	// 是否分页
	if *params.IsPage == 0 {
		offset := (*params.Page - 1) * *params.Limit
		if err := query.Offset(offset).Limit(*params.Limit).Find(&adminUsers).Error; err != nil {
			return nil, 0, err
		}
		return adminUsers, count, nil
	} else {
		if err := query.Find(&adminUsers).Error; err != nil {
			return nil, 0, err
		}
	}
	return adminUsers, count, nil
}

// UpdateAdminUser 更新管理员用户信息
func UpdateAdminUser(adminUser *AdminUser) error {
	db := getDB()
	return db.Save(adminUser).Error
}

// DeleteAdminUser 删除管理员用户（软删除）
func DeleteAdminUser(id uint64) error {
	db := getDB()
	return db.Model(&AdminUser{}).Where("id = ?", id).Update("is_deleted", true).Error
}

// UpdateAdminUserPassword 更新管理员用户密码
func UpdateAdminUserPassword(id uint64, password string) error {
	db := getDB()
	return db.Model(&AdminUser{}).Where("id = ?", id).Update("password", password).Error
}

// UpdateAdminUserStatus 更新管理员用户状态
func UpdateAdminUserStatus(id uint64, status int) error {
	db := getDB()
	return db.Model(&AdminUser{}).Where("id = ?", id).Update("status", status).Error
}
