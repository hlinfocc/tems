package models

// Banner 轮播图模型
type Banner struct {
	BaseModel
	Title     string `gorm:"column:title;type:varchar(255);not null;comment:轮播图标题" json:"title"`
	ImageURL  string `gorm:"column:image_url;type:varchar(500);not null;comment:图片URL" json:"image_url"`
	IsVisible bool   `gorm:"column:is_visible;type:boolean;default:true;comment:是否可见" json:"is_visible"`
	Sort      int    `gorm:"column:sort;type:int;default:0;comment:排序值" json:"sort"`
}

// TableName 指定表名
func (Banner) TableName() string {
	return "banners"
}

// CreateBanner 创建轮播图
func CreateBanner(banner *Banner) error {
	db := getDB()
	return db.Create(banner).Error
}

// GetBannerByID 根据ID获取轮播图
func GetBannerByID(id uint64) (*Banner, error) {
	db := getDB()
	var banner Banner
	err := db.First(&banner, id).Error
	if err != nil {
		return nil, err
	}
	return &banner, nil
}

// GetBanners 获取轮播图列表
func GetBanners(limit int) ([]Banner, error) {
	db := getDB()
	var banners []Banner
	query := db.Where("is_deleted = ? AND is_visible = ?", false, true).Order("sort ASC, id DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&banners).Error
	return banners, err
}

// GetAllBanners 获取所有轮播图
func GetAllBanners() ([]Banner, error) {
	db := getDB()
	var banners []Banner
	err := db.Where("is_deleted = ?", false).Order("sort ASC, id DESC").Find(&banners).Error
	return banners, err
}

// UpdateBanner 更新轮播图
func UpdateBanner(banner *Banner) error {
	db := getDB()
	return db.Save(banner).Error
}

// DeleteBanner 删除轮播图（软删除）
func DeleteBanner(id uint64) error {
	db := getDB()
	return db.Model(&Banner{}).Where("id = ?", id).Update("is_deleted", true).Error
}

// UpdateBannerVisibility 更新轮播图可见性
func UpdateBannerVisibility(id uint64, isVisible bool) error {
	db := getDB()
	return db.Model(&Banner{}).Where("id = ?", id).Update("is_visible", isVisible).Error
}
