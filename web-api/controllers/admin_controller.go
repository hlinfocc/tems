package controllers

import (
	"fmt"
	"log"
	"strings"

	"net/http"
	"strconv"
	"sync"
	"tems-web-api/models"
	"tems-web-api/utils"
	"tems-web-api/utils/captcha"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mozillazg/go-pinyin"
)

// AdminLoginRequest 管理员登录请求
type AdminLoginRequest struct {
	Account   string `json:"account" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Captcha   string `json:"captcha" binding:"required"`
	CaptchaID string `json:"captchaId" binding:"required"`
}

// CreateAdminUserRequest 创建管理员用户请求
type CreateAdminUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password"`
	Status   int    `json:"status"`
	UserType int    `json:"user_type"`
}

// UpdateAdminUserRequest 更新管理员用户请求
type UpdateAdminUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Status   int    `json:"status"`
	UserType int    `json:"user_type"`
}

// UpdateAdminUserPasswordRequest 更新管理员用户密码请求
type UpdateAdminUserPasswordRequest struct {
	Password string `json:"password" binding:"required"`
}

// UpdateAdminUserStatusRequest 更新管理员用户状态请求
type UpdateAdminUserStatusRequest struct {
	Status int `json:"status" binding:"required"`
}

// captchaInfo 验证码信息
type captchaInfo struct {
	value      string
	expiration time.Time
}

// 验证码存储
var captchaStore = struct {
	mu         sync.RWMutex
	captchas   map[string]captchaInfo
	expiration time.Duration
}{captchas: make(map[string]captchaInfo), expiration: 5 * time.Minute}

// GenerateCaptcha 生成验证码
func GenerateCaptcha(ctx *gin.Context) {
	reqtype := ctx.DefaultQuery("type", "base64")
	if reqtype == "base64" {
		base64Data, code, err := captcha.GenerateBase64()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "生成验证码失败"})
			return
		}

		captchaId := utils.Sha3(fmt.Sprintf("%d-%s", time.Now().UnixNano(), utils.RandStringBytes(20, false)))
		// 缓存code
		cacheCaptchaInfo(captchaId, code)
		// log.Printf("verifycode：%s,=============key：%s\n", code, key)
		ctx.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "",
			"key":  captchaId,
			"data": base64Data,
		})
	} else {
		key, _ := ctx.GetQuery("key")
		captchaGen := captcha.NewCaptchaGenerator(140, 50, 5)
		// 生成验证码
		img, code, err := captchaGen.Generate()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "生成验证码失败"})
			return
		}
		// 转换为JPEG
		imgBytes, err := captchaGen.ImageToJPEG(img)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "编码图片失败"})
			return
		}
		if key == "" {
			key = utils.RandStringBytes(20, false)
		}
		// log.Printf("verifycode：%s ,========key：%s\n", code, key)
		// 缓存code
		cacheCaptchaInfo(key, code)
		// 设置响应头
		ctx.Header("Content-Type", "image/jpeg")
		// 返回图片
		ctx.Data(http.StatusOK, "image/jpeg", imgBytes)
	}
}

// cacheCaptchaInfo 缓存验证码信息
func cacheCaptchaInfo(captchaID string, code string) {
	// 清理过期的验证码
	cleanExpiredCaptchas()
	// 保存验证码
	captchaStore.mu.Lock()
	captchaKey := fmt.Sprintf("%s:%s", captchaID, strings.ToLower(code))
	captchaStore.captchas[captchaKey] = captchaInfo{
		value:      code,
		expiration: time.Now().Add(captchaStore.expiration),
	}
	captchaStore.mu.Unlock()
}

// cleanExpiredCaptchas 清理过期的验证码
func cleanExpiredCaptchas() {
	now := time.Now()
	captchaStore.mu.Lock()
	defer captchaStore.mu.Unlock()

	for id, info := range captchaStore.captchas {
		if now.After(info.expiration) {
			delete(captchaStore.captchas, id)
		}
	}
}

// validateCaptcha 验证验证码
func validateCaptcha(captchaId, captchaValue string) bool {
	captchaStore.mu.Lock()
	defer captchaStore.mu.Unlock()
	key := fmt.Sprintf("%s:%s", captchaId, strings.ToLower(captchaValue))
	info, exists := captchaStore.captchas[key]
	if !exists {
		log.Printf("验证码不存在：%s", key)
		return false
	}

	if time.Now().After(info.expiration) {
		log.Printf("验证码过期：%s", key)
		delete(captchaStore.captchas, key)
		return false
	}

	// 验证码不区分大小写
	isValid := strings.EqualFold(info.value, captchaValue)

	// 验证成功后删除验证码，防止重复使用
	if isValid {
		log.Printf("验证码验证成功：%s", key)
		delete(captchaStore.captchas, key)
	}
	return isValid
}

// AdminLogin 管理员登录
func AdminLogin(c *gin.Context) {
	var req AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "请求参数无效",
			"error": err.Error(),
		})
		return
	}

	// 验证验证码
	if !validateCaptcha(req.CaptchaID, req.Captcha) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "验证码错误或已过期",
		})
		return
	}

	// 查询管理员用户
	user, err := models.GetAdminUserByUsername(req.Account)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "用户名或密码错误",
		})
		return
	}

	// 验证密码
	if user.Password != utils.Sha3(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "用户名或密码错误",
		})
		return
	}

	// 检查用户状态
	if user.Status == 1 {
		c.JSON(http.StatusForbidden, gin.H{
			"code": 403,
			"msg":  "用户账号已禁用",
		})
		return
	}
	currRole := "teacher"
	if user.UserType == 0 {
		currRole = "admin"
	}
	// 生成JWT token
	token, err := utils.GenJwtToken(
		user.ID,       // 管理员ID
		user.Name,     // 管理员姓名
		0,             // 班级ID：0表示无班级
		user.Username, // 用户名
		user.Status,   // 状态：1表示正常
		1,             // 用户类型：2表示学生
		currRole,      // 角色
		0,             // App类型：1表示小程序
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "生成令牌失败",
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "登录成功",
		"data": gin.H{
			"token": token,
			"userInfo": gin.H{
				"id":           user.ID,
				"adminName":    user.Name,
				"adminAccount": user.Username,
				"adminLevel":   user.UserType,
				"role":         currRole,
			},
		},
	})
}

func CheckAdminLogin(c *gin.Context) {
	// 从请求头中获取token
	token := c.GetHeader("Authorization")
	if token == "" {
		token = c.GetHeader("token")
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "登录成功",
		"data": gin.H{
			"token":    token,
			"userInfo": gin.H{},
		},
	})
}

func GetAdminUserInfo(c *gin.Context) {
	// 从请求头中获取token
	token := c.GetHeader("Authorization")
	if token == "" {
		token = c.GetHeader("token")
	}
	customClaims, e := utils.JWTParse(token, 0)
	if e != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "鉴权失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "登录成功",
		"data": gin.H{
			"token": token,
			"userInfo": gin.H{
				"id":           customClaims.UserID,
				"adminName":    customClaims.RealName,
				"adminAccount": customClaims.Account,
				"adminLevel":   customClaims.UserType,
				"role":         customClaims.Role,
			},
		},
	})
}

// GetAdminUsers 获取管理员用户分页列表
func GetAdminUsers(c *gin.Context) {
	var params models.AdminUserRequestParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "无效的请求参数",
			"error": err.Error(),
		})
		return
	}
	// 设置默认值
	if *params.Page <= 0 {
		*params.Page = 1
	}
	if *params.Limit <= 0 {
		*params.Limit = 10
	}

	// 获取所有管理员用户
	adminUsers, total, err := models.GetAllAdminUsers(&params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "获取管理员用户失败",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": gin.H{
			"list":  adminUsers,
			"total": total,
		},
	})
}

// CreateAdminUser 创建管理员用户
func CreateAdminUser(c *gin.Context) {
	var req CreateAdminUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "无效的请求参数",
			"error": err.Error(),
		})
		return
	}

	adminUser := &models.AdminUser{
		Name:     req.Name,
		Username: req.Username,
		Password: utils.Sha3(req.Password),
		Status:   req.Status,
		UserType: req.UserType,
	}

	if err := models.CreateAdminUser(adminUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "创建管理员用户失败",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "管理员用户创建成功",
		"data": adminUser,
	})
}

// UpdateAdminUser 更新管理员用户
func UpdateAdminUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的管理员用户ID",
		})
		return
	}

	var req UpdateAdminUserRequest
	if err1 := c.ShouldBindJSON(&req); err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "无效的请求参数",
			"error": err1.Error(),
		})
		return
	}

	adminUser, err := models.GetAdminUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "管理员用户不存在",
		})
		return
	}

	adminUser.Name = req.Name
	adminUser.Username = req.Username
	adminUser.Status = req.Status
	adminUser.UserType = req.UserType

	if err := models.UpdateAdminUser(adminUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "更新管理员用户失败",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "管理员用户更新成功",
		"data": adminUser,
	})
}

// DeleteAdminUser 删除管理员用户
func DeleteAdminUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的管理员用户ID",
		})
		return
	}

	// 检查用户是否存在
	_, err = models.GetAdminUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "管理员用户不存在",
		})
		return
	}

	if err := models.DeleteAdminUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "删除管理员用户失败",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "管理员用户删除成功",
	})
}

// UpdateAdminUserStatus 更新管理员用户状态
func UpdateAdminUserStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的管理员用户ID",
		})
		return
	}

	var req UpdateAdminUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "无效的请求参数",
			"error": err.Error(),
		})
		return
	}

	if err := models.UpdateAdminUserStatus(id, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "更新管理员用户状态失败",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "管理员用户状态更新成功",
	})
}

// UpdateAdminUserPassword 更新管理员用户密码
func UpdateAdminUserPassword(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的管理员用户ID",
		})
		return
	}

	var req UpdateAdminUserPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "无效的请求参数",
			"error": err.Error(),
		})
		return
	}

	if err := models.UpdateAdminUserPassword(id, utils.Sha3(req.Password)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "更新管理员用户密码失败",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "管理员用户密码更新成功",
	})
}

// GetPinyinRequest 获取汉语拼音请求
type GetPinyinRequest struct {
	Name string `json:"name" binding:"required"`
}

// GetPinyin 获取传入姓名的汉语拼音
func GetPinyin(c *gin.Context) {
	var req GetPinyinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "无效的请求参数",
			"error": err.Error(),
		})
		return
	}

	// 设置拼音转换选项
	args := pinyin.NewArgs()
	args.Style = pinyin.Normal
	args.Heteronym = true

	// 转换为拼音
	pinyinSlice := pinyin.Pinyin(req.Name, args)

	// 拼接拼音结果
	var pinyinResult []string
	for _, p := range pinyinSlice {
		pinyinResult = append(pinyinResult, p[0])
	}

	pinyinStr := strings.Join(pinyinResult, "")
	log.Printf("姓名 %s 的汉语拼音为 %s", req.Name, pinyinStr)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取汉语拼音成功",
		"data": gin.H{
			"name":   req.Name,
			"pinyin": pinyinStr,
		},
	})
}
