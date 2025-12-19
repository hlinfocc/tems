package controllers

import (
	"net/http"
	"strconv"
	"tems-web-api/models"

	"github.com/gin-gonic/gin"
)

// CreateBannerRequest 创建轮播图请求
type CreateBannerRequest struct {
	Title     string `json:"title" binding:"required"`
	ImageURL  string `json:"image_url" binding:"required"`
	IsVisible bool   `json:"is_visible"`
	Sort      int    `json:"sort"`
}

// UpdateBannerRequest 更新轮播图请求
type UpdateBannerRequest struct {
	ID        uint64   `json:"id" binding:"required"`
	Title     string `json:"title" binding:"required"`
	ImageURL  string `json:"image_url" binding:"required"`
	IsVisible bool   `json:"is_visible"`
	Sort      int    `json:"sort"`
}

// UpdateBannerVisibilityRequest 更新轮播图可见性请求
type UpdateBannerVisibilityRequest struct {
	ID        uint64   `json:"id"`
	IsVisible *bool `json:"is_visible" binding:"required"`
}

// GetBanners 获取轮播图列表（前端展示用）
func GetBanners(c *gin.Context) {
	// 获取查询参数
	limitStr := c.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 5
	}

	// 查询轮播图
	banners, err := models.GetBanners(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"msg": "获取轮播图失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg": "成功",
		"data":    banners,
	})
}

// GetAllBanners 获取所有轮播图分页列表（管理后台用）
func GetAllBanners(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	// 获取所有轮播图
	banners, err := models.GetAllBanners()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"msg": "获取轮播图失败",
			"error":   err.Error(),
		})
		return
	}

	// 计算分页
	total := len(banners)
	start := (page - 1) * pageSize
	end := start + pageSize
	if start >= total {
		start, end = 0, 0
	} else if end > total {
		end = total
	}

	var paginatedBanners []models.Banner
	if start < end {
		paginatedBanners = banners[start:end]
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg": "成功",
		"data": gin.H{
			"list":      paginatedBanners,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
			"pages":     (total + pageSize - 1) / pageSize,
		},
	})
}

// CreateBanner 创建轮播图
func CreateBanner(c *gin.Context) {
	var req CreateBannerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"msg": "请求参数无效",
			"error":   err.Error(),
		})
		return
	}

	banner := &models.Banner{
		Title:     req.Title,
		ImageURL:  req.ImageURL,
		IsVisible: req.IsVisible,
		Sort:      req.Sort,
	}

	if err := models.CreateBanner(banner); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"msg": "创建轮播图失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg": "轮播图创建成功",
		"data":    banner,
	})
}

// UpdateBanner 更新轮播图
func UpdateBanner(c *gin.Context) {
	var req UpdateBannerRequest
	if err1 := c.ShouldBindJSON(&req); err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"msg": "请求参数无效",
			"error":   err1.Error(),
		})
		return
	}

	banner, err := models.GetBannerByID(req.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"msg": "轮播图不存在",
		})
		return
	}

	banner.Title = req.Title
	banner.ImageURL = req.ImageURL
	banner.IsVisible = req.IsVisible
	banner.Sort = req.Sort

	if err := models.UpdateBanner(banner); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"msg": "更新轮播图失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg": "轮播图更新成功",
		"data":    banner,
	})
}

// DeleteBanner 删除轮播图
func DeleteBanner(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"msg": "Invalid banner ID",
		})
		return
	}

	// 检查轮播图是否存在
	_, err = models.GetBannerByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"msg": "Banner not found",
		})
		return
	}

	if err := models.DeleteBanner(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"msg": "删除轮播图失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg": "轮播图删除成功",
	})
}

// UpdateBannerVisibility 更新轮播图可见性
func UpdateBannerVisibility(c *gin.Context) {
	var req UpdateBannerVisibilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"msg": "请求参数无效",
			"error":   err.Error(),
		})
		return
	}

	if err := models.UpdateBannerVisibility(req.ID, *req.IsVisible); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"msg": "更新轮播图可见性失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"msg": "轮播图可见性更新成功",
	})
}