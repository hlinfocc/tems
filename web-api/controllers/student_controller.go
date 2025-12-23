package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"tems-web-api/config"
	"tems-web-api/models"
	"tems-web-api/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// StudentLoginRequest 学生登录请求
type StudentLoginRequest struct {
	Code   string `json:"code" binding:"required"` // 获取手机号的code
	OpenID string `json:"open_id" binding:"required"`
}

// StudentBindRequest 学生绑定请求
type StudentBindRequest struct {
	StudentName string `json:"studentName" binding:"required"`
	StudentID   string `json:"studentId" binding:"required"`
	ClassName   string `json:"className"`
	ClassID     uint64 `json:"classId"`
	Phone       string `json:"phone" binding:"required"` // 获取手机号的code
	OpenID      string `json:"openId" binding:"required"`
}

// 微信access_token缓存结构
type AccessTokenCache struct {
	AccessToken string
	ExpiresIn   int64
	ExpireTime  int64 // 过期时间戳
}

var (
	accessTokenCache AccessTokenCache
	accessTokenMutex sync.Mutex
)

// 微信获取手机号的响应结构
type WechatPhoneResponse struct {
	Errcode   int    `json:"errcode"`
	Errmsg    string `json:"errmsg"`
	PhoneInfo struct {
		PhoneNumber     string `json:"phoneNumber"`
		PurePhoneNumber string `json:"purePhoneNumber"`
		CountryCode     string `json:"countryCode"`
	} `json:"phone_info"`
}

var openRegister bool = true

func initEnvParams() {
	if val := os.Getenv("OPEN_REGISTER"); val != "" {
		openRegister, _ = strconv.ParseBool(val)
	}
}

// 获取微信access_token
func getWechatAccessToken() (string, error) {
	accessTokenMutex.Lock()
	defer accessTokenMutex.Unlock()

	// 检查缓存是否有效
	now := time.Now().Unix()
	if accessTokenCache.AccessToken != "" && now < accessTokenCache.ExpireTime {
		return accessTokenCache.AccessToken, nil
	}

	// 获取配置
	cfg := config.GetGlobalConfig()

	// 调用获取access_token的接口
	accessTokenURL := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",
		cfg.WechatAppID, cfg.WechatSecret)

	resp, err := http.Get(accessTokenURL)
	if err != nil {
		return "", fmt.Errorf("获取access_token失败: %v", err)
	}
	defer resp.Body.Close()

	type AccessTokenResponse struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
		Errcode     int    `json:"errcode,omitempty"`
		Errmsg      string `json:"errmsg,omitempty"`
	}

	var tokenResp AccessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("解析access_token响应失败: %v", err)
	}

	if tokenResp.Errcode != 0 {
		return "", fmt.Errorf("获取access_token错误: %s", tokenResp.Errmsg)
	}

	// 更新缓存，设置过期时间（提前10分钟过期）
	expireTime := now + tokenResp.ExpiresIn - 600
	accessTokenCache = AccessTokenCache{
		AccessToken: tokenResp.AccessToken,
		ExpiresIn:   tokenResp.ExpiresIn,
		ExpireTime:  expireTime,
	}

	return tokenResp.AccessToken, nil
}

// 获取微信手机号
func getWechatPhone(encryptedData string) (string, error) {
	// 获取access_token
	token, err := getWechatAccessToken()
	if err != nil {
		return "", err
	}

	// 调用微信获取手机号的接口
	phoneURL := fmt.Sprintf("https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token=%s", token)

	// 构建请求体
	type PhoneRequest struct {
		Code string `json:"code"`
	}
	phoneReq := PhoneRequest{Code: encryptedData}
	reqBody, _ := json.Marshal(phoneReq)

	// 创建请求
	req, _ := http.NewRequest("POST", phoneURL, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	phoneResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("获取手机号失败: %v", err)
	}
	defer phoneResp.Body.Close()

	// 解析响应
	var wechatPhoneResp WechatPhoneResponse
	if err := json.NewDecoder(phoneResp.Body).Decode(&wechatPhoneResp); err != nil {
		return "", fmt.Errorf("解析手机号响应失败: %v", err)
	}

	if wechatPhoneResp.Errcode != 0 {
		// 如果是access_token过期或无效的错误，清除缓存并重试一次
		if wechatPhoneResp.Errcode == 40001 || wechatPhoneResp.Errcode == 42001 {
			accessTokenMutex.Lock()
			accessTokenCache = AccessTokenCache{} // 清除缓存
			accessTokenMutex.Unlock()

			// 重新获取access_token并再次调用获取手机号接口
			return getWechatPhone(encryptedData)
		}
		return "", fmt.Errorf("获取手机号错误: %s", wechatPhoneResp.Errmsg)
	}

	return wechatPhoneResp.PhoneInfo.PurePhoneNumber, nil
}

// StudentLogin 学生登录
func StudentLogin(c *gin.Context) {
	var req StudentLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "请求参数无效",
			"error": err.Error(),
		})
		return
	}

	// 获取真实手机号（前端传来的是获取手机号的code）
	phone, err := getWechatPhone(req.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "获取手机号失败",
			"error": err.Error(),
		})
		return
	}

	// 根据手机号查询学生
	student, err := models.GetStudentByPhone(phone)
	if err != nil {
		// 学生不存在，返回提示前端绑定学生信息
		c.JSON(http.StatusOK, gin.H{
			"code":      200,
			"msg":       "学生未找到，请绑定学生信息",
			"need_bind": true,
			"data": gin.H{
				"phone":   phone,
				"open_id": req.OpenID,
			},
		})
		return
	}

	// 学生存在，更新OpenID
	student.OpenID = req.OpenID
	student.IsBound = true
	student.UpdatedAt = models.Now()

	// 保存更新
	if uerr := models.UpdateStudent(student); uerr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "更新学生信息失败",
			"error": uerr.Error(),
		})
		return
	}

	// 生成JWT token
	token, err := utils.GenJwtToken(
		student.ID,
		student.StudentName,
		uint64(student.ClassID),
		student.StudentID,
		1,         // 状态：1表示正常
		2,         // 用户类型：2表示学生
		"student", // 角色
		1,         // App类型：1表示小程序
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "生成令牌失败",
			"error": err.Error(),
		})
		return
	}

	// 返回登录成功和token
	c.JSON(http.StatusOK, gin.H{
		"code":      200,
		"msg":       "登录成功",
		"need_bind": false,
		"data": gin.H{
			"userId":       student.ID,
			"student_id":   student.StudentID,
			"student_name": student.StudentName,
			"class_name":   student.ClassName,
			"phone":        student.Phone,
			"token":        token,
		},
	})
}

// BindStudent 绑定学生信息
func BindStudent(c *gin.Context) {
	var req StudentBindRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "请求参数无效",
			"error": err.Error(),
		})
		return
	}
	initEnvParams()
	var student models.Student
	if openRegister {
		student = models.Student{
			StudentName: req.StudentName,
			StudentID:   req.StudentID,
			Phone:       req.Phone,
			ClassName:   req.ClassName,
			ClassID:     req.ClassID,
			OpenID:      req.OpenID,
			IsBound:     true,
		}
		if req.ClassID <= 0 && req.ClassName == "" {
			c.JSON(http.StatusNotFound, gin.H{
				"code":  404,
				"msg":   "请输入正确的班级",
				"error": "",
			})
			return
		}
		if req.ClassID > 0 {
			classes, ce := models.GetClassByID(req.ClassID)
			if ce != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"code":  404,
					"msg":   "班级信息不存在",
					"error": "",
				})
				return
			}
			student.ClassName = classes.ClassName
		}
		if req.ClassID <= 0 && req.ClassName != "" {
			classList, cle := models.GetClassesListByName(req.ClassName, false)
			if cle != nil && len(classList) != 1 {
				c.JSON(http.StatusNotFound, gin.H{
					"code":  404,
					"msg":   "请输入正确的班级名称",
					"error": "",
				})
				return
			}
			student.ClassID = classList[0].ID
		}

		qty := models.GetCountStudentByNameIdClassId(req.StudentName, req.StudentID, req.ClassID)
		if qty > 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"code":  404,
				"msg":   "该学生姓名及学号已存在",
				"error": "请检查姓名、学号和班级信息是否正确",
			})
			return
		}
		if uerr := models.CreateStudent(&student); uerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":  500,
				"msg":   "保存信息失败，请重试",
				"error": "",
			})
			return
		}
	} else {
		// 根据姓名、学号、班级匹配学生
		student, err := models.GetStudentByNameIdClass(req.StudentName, req.StudentID, req.ClassName)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"code":  404,
				"msg":   "学生信息不匹配",
				"error": "请检查姓名、学号和班级信息是否正确",
			})
			return
		}
		// 更新绑定信息
		student.OpenID = req.OpenID
		student.IsBound = true
		student.Phone = req.Phone
		student.UpdatedAt = models.Now()

		// 保存更新
		if uerr := models.UpdateStudent(student); uerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":  500,
				"msg":   "更新学生信息失败",
				"error": uerr.Error(),
			})
			return
		}
	}

	// 生成JWT token
	token, err := utils.GenJwtToken(
		student.ID,
		student.StudentName,
		uint64(student.ClassID),
		student.StudentID,
		1,         // 状态：1表示正常
		2,         // 用户类型：2表示学生
		"student", // 角色
		1,         // App类型：1表示小程序
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
		"msg":  "学生绑定成功",
		"data": gin.H{
			"userId":       student.ID,
			"student_id":   student.StudentID,
			"student_name": student.StudentName,
			"class_name":   student.ClassName,
			"token":        token,
		},
	})
}

// GetStudentByID 根据ID获取学生信息
func GetStudentByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的学生ID",
		})
		return
	}

	student, err := models.GetStudentByID(uint64(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "学生不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "成功",
		"data": student,
	})
}

// 学生管理接口 - 管理后台

// StudentListRequest 学生列表请求
type StudentListRequest struct {
	StudentName string `form:"studentName" json:"studentName"`
	StudentID   string `form:"studentID" json:"studentID"`
	ClassID     uint64 `form:"classID" json:"classID"`
	Page        int    `form:"page" json:"page"`
	PageSize    int    `form:"pageSize" json:"pageSize"`
}

// StudentListResponse 学生列表响应
type StudentListResponse struct {
	List  []models.Student `json:"list"`
	Total int64            `json:"total"`
}

// StudentAddRequest 学生添加请求
type StudentAddRequest struct {
	StudentName string `json:"studentName" binding:"required"`
	StudentID   string `json:"studentID" binding:"required"`
	Phone       string `json:"phone"`
	ClassName   string `json:"className" binding:"required"`
	ClassID     uint64 `json:"classID"`
	Password    string `json:"password"`
}

// StudentUpdateRequest 学生更新请求
type StudentUpdateRequest struct {
	ID          uint64 `json:"id" binding:"required"`
	StudentName string `json:"studentName" binding:"required"`
	StudentID   string `json:"studentID" binding:"required"`
	Phone       string `json:"phone"`
	ClassName   string `json:"className" binding:"required"`
	ClassID     uint64 `json:"classID" binding:"required"`
	Password    string `json:"password"`
}

// GetStudentList 分页获取学生列表
func GetStudentList(c *gin.Context) {
	var req StudentListRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "请求参数无效",
			"error": err.Error(),
		})
		return
	}

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	customClaimsRole := c.GetString("customClaimsRole")
	customClaimsUserId := c.GetUint64("customClaimsUserId")

	// 构建查询参数
	params := models.StudentQueryParams{
		StudentName: req.StudentName,
		StudentID:   req.StudentID,
		ClassID:     req.ClassID,
		Page:        req.Page,
		PageSize:    req.PageSize,
	}
	if customClaimsRole == "classAdvisor" {
		params.ClassIds = models.GetClassIdByClassAdvisor(customClaimsUserId)
	}
	// 查询学生列表
	students, total, err := models.GetStudentsWithPagination(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "获取学生列表失败",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取学生列表成功",
		"data": StudentListResponse{
			List:  students,
			Total: total,
		},
	})
}

// AddStudent 添加学生
func AddStudent(c *gin.Context) {
	var req StudentAddRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "请求参数无效",
			"error": err.Error(),
		})
		return
	}

	// 检查学号是否已存在
	_, err := models.GetStudentByStudentID(req.StudentID)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "学号已存在",
		})
		return
	}

	// 创建学生
	student := &models.Student{
		StudentName: req.StudentName,
		StudentID:   req.StudentID,
		Phone:       req.Phone,
		ClassName:   req.ClassName,
		ClassID:     req.ClassID,
		Password:    req.Password,
	}

	if err := models.CreateStudent(student); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "添加学生失败",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "添加学生成功",
		"data": student,
	})
}

// AddBatchStudent 批量添加学生
func AddBatchStudent(c *gin.Context) {
	var req []StudentAddRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "请求参数无效",
			"error": err.Error(),
		})
		return
	}
	var students []*models.Student
	// 检查班级是否存在
	var errClassNameMsg []string
	for _, item := range req {
		classInfo, err := models.GetClassesByName(item.ClassName)
		if err != nil {
			log.Printf("班级不存在: %v", err)
			if !contains(errClassNameMsg, item.ClassName) {
				errClassNameMsg = append(errClassNameMsg, item.ClassName)
			}
			continue
		}
		// 创建学生
		student := &models.Student{
			StudentName: item.StudentName,
			StudentID:   item.StudentID,
			Phone:       item.Phone,
			ClassName:   classInfo.ClassName,
			ClassID:     classInfo.ID,
			Password:    "",
		}
		students = append(students, student)
	}

	// 批量创建学生
	if err := models.CreateBatchStudent(students); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "添加学生失败",
			"error": err.Error(),
		})
		return
	}
	if len(errClassNameMsg) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  fmt.Sprintf("部分添加成功，以下班级不存在：%s", strings.Join(errClassNameMsg, ",")),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  fmt.Sprintf("添加成功，共添加 %d 条记录", len(students)),
		})
	}
}

func contains(errClassNameMsg []string, s string) bool {
	panic("unimplemented")
}

// UpdateStudent 更新学生信息
func UpdateStudent(c *gin.Context) {
	var req StudentUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "请求参数无效",
			"error": err.Error(),
		})
		return
	}

	// 获取学生
	student, err := models.GetStudentByID(req.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "学生不存在",
		})
		return
	}

	// 检查学号是否已存在（排除当前学生）
	existingStudent, err := models.GetStudentByStudentID(req.StudentID)
	if err == nil && existingStudent.ID != req.ID {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "学号已存在",
		})
		return
	}

	// 更新学生信息
	student.StudentName = req.StudentName
	student.StudentID = req.StudentID
	student.Phone = req.Phone
	student.ClassName = req.ClassName
	student.ClassID = req.ClassID
	if req.Password != "" {
		student.Password = req.Password
	}
	student.UpdatedAt = models.Now()

	if err := models.UpdateStudent(student); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "更新学生信息失败",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "更新学生信息成功",
		"data": student,
	})
}

// DeleteStudent 删除学生
func DeleteStudent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的学生ID",
		})
		return
	}

	// 检查学生是否存在
	_, err = models.GetStudentByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "学生不存在",
		})
		return
	}

	// 删除学生
	if err := models.DeleteStudent(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "删除学生失败",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除学生成功",
	})
}
