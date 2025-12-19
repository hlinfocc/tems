package models

import (
	"database/sql/driver"
	"fmt"
	"tems-web-api/config"
	"tems-web-api/utils"
	"time"

	"gorm.io/gorm"
)

// 自定义时间类型
type JsonTime time.Time

// 定义时间格式常量
const (
	TimeFormat    = "2006-01-02 15:04:05"
	DateFormat    = "2006-01-02"
	ISO8601Format = "2006-01-02T15:04:05Z07:00"
)

// 实现 json.Marshaler 接口
func (j JsonTime) MarshalJSON() ([]byte, error) {
	if time.Time(j).IsZero() {
		return []byte("null"), nil
	}
	return []byte(`"` + time.Time(j).Format(TimeFormat) + `"`), nil
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (jt *JsonTime) UnmarshalJSON(data []byte) error {
	// 处理 null 值
	if string(data) == "null" {
		*jt = JsonTime(time.Time{})
		return nil
	}

	// 去除引号
	if len(data) < 2 {
		return fmt.Errorf("invalid time format")
	}
	strData := string(data[1 : len(data)-1])

	// 尝试多种时间格式
	var t time.Time
	var err error

	formats := []string{
		TimeFormat,
		DateFormat,
		ISO8601Format,
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02 15:04:05.999999999 -0700 MST", // 包含时区
	}

	for _, format := range formats {
		t, err = time.Parse(format, strData)
		if err == nil {
			*jt = JsonTime(t)
			return nil
		}
	}

	return fmt.Errorf("cannot parse time: %s, supported formats: %v", strData, formats)
}

// String 实现 Stringer 接口
func (jt JsonTime) String() string {
	return time.Time(jt).Format(TimeFormat)
}

// Time 转换为 time.Time
func (jt JsonTime) Time() time.Time {
	return time.Time(jt)
}

// Now 获取当前时间的 JsonTime
func Now() JsonTime {
	return JsonTime(time.Now())
}

// FromTime 从 time.Time 创建 JsonTime
func FromTime(t time.Time) JsonTime {
	return JsonTime(t)
}

// 实现 gorm.Scanner 接口
func (j *JsonTime) Scan(value interface{}) error {
	if value == nil {
		*j = JsonTime(time.Time{})
		return nil
	}
	if t, ok := value.(time.Time); ok {
		*j = JsonTime(t)
		return nil
	}
	return fmt.Errorf("can not convert %v to JsonTime", value)
}

// 实现 gorm.Valuer 接口
func (j JsonTime) Value() (driver.Value, error) {
	if time.Time(j).IsZero() {
		return nil, nil
	}
	return time.Time(j), nil
}

// BaseModel 基础模型结构，包含通用字段
type BaseModel struct {
	ID        uint64         `gorm:"primaryKey;comment:主键ID" json:"id"`
	CreatedAt JsonTime       `gorm:"autoCreateTime;comment:创建时间" json:"createdAt"`
	UpdatedAt JsonTime       `gorm:"autoUpdateTime;comment:更新时间" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间" json:"-"`
	IsDeleted bool           `gorm:"default:false;comment:是否删除" json:"isDeleted"`
}

// AutoMigrate 自动迁移所有数据库表
func AutoMigrate() error {
	db := getDB()
	if db == nil {
		return nil
	}

	// 自动迁移所有模型
	return db.AutoMigrate(
		&Banner{},
		&Student{},
		&AdminUser{},
		&Class{},
		&Course{},
		&ClassCourse{},
		&EvaluationTask{},
		&EvaluationTaskDetail{},
		&EvaluationQuestionSet{},
		&EvaluationQuestionDetail{},
		&EvaluationResult{},
	)
}

// InitManagerUser 初始化默认管理员用户
func InitManagerUser() error {
	db := getDB()
	if db == nil {
		return nil
	}

	// 检查用户是否存在
	var count int64
	db.Model(&AdminUser{}).Where("username = ?", "manager").Count(&count)
	if count > 0 {
		return nil // 用户已存在
	}

	// 创建默认管理员用户
	manager := &AdminUser{
		Name:     "超级管理员",
		Username: "manager",
		Password: utils.Sha3("Aa123456"), // 密码需要加密存储
		Status:   0,
		UserType: 0,
	}
	return db.Create(manager).Error
}

// 获取数据库连接的辅助函数
func getDB() *gorm.DB {
	return config.GetDB()
}

// GetDB 获取数据库连接（公开接口）
func GetDB() *gorm.DB {
	return getDB()
}
