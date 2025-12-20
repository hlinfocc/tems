package uploads

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"tems-web-api/assets"
	"tems-web-api/config"
	"tems-web-api/utils"
	"time"

	"github.com/gin-gonic/gin"
)

var maxFileSize int64 = 1024 * 1024 * 10 // 10MB
var allowedFileTypes = []string{"image/jpeg", "image/png", "image/gif"}

func initConfig() {
	if val := os.Getenv("MAX_FILE_SIZE"); val != "" {
		maxFileSize, _ = strconv.ParseInt(val, 10, 64)
	}
	if val := os.Getenv("ALLOWED_FILE_TYPES"); val != "" {
		allowedFileTypes = strings.Split(val, ",")
	}
}

// 文件类型检查函数
func checkFileType(file *multipart.FileHeader) bool {
	// 打开文件读取前512字节用于检测真实类型
	src, err := file.Open()
	if err != nil {
		return false
	}
	defer src.Close()

	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil && err != io.EOF {
		return false
	}

	// 通过魔数检测文件真实类型
	fileType := http.DetectContentType(buffer)

	// 检查是否在允许的类型列表中
	for _, allowedType := range allowedFileTypes {
		if strings.HasPrefix(fileType, allowedType) {
			return true
		}
	}
	return false
}

// 生成唯一文件名
func generateFileName(originalName string) string {
	// 使用时间戳+MD5生成唯一文件名
	timestamp := time.Now().UnixNano()
	newFileName := utils.Sha1(fmt.Sprintf("%d%s%s", timestamp, originalName, utils.RandStringBytes(10, false)))
	return newFileName
}

// 获取动态子目录路径
func getDateDir(baseDir string, dirFormat string) string {
	now := time.Now()
	switch dirFormat {
	case "year":
		year := now.Format("2006")
		return path.Join(baseDir, year)
	case "month":
		year := now.Format("2006")
		month := now.Format("01")
		return path.Join(baseDir, year, month)
	default:
		year := now.Format("2006")
		month := now.Format("01")
		day := now.Format("02")
		return path.Join(baseDir, year, month, day)
	}
}

func createDirectoryTreeWithIndex(exeBaseDir, fullPath string) error {

	// 创建目标目录
	if err := os.MkdirAll(fullPath, 0644); err != nil {
		return err
	}

	// 为目录创建空index.html
	parts := strings.Split(strings.Replace(fullPath, exeBaseDir, "", 1), string(filepath.Separator))
	currentPath := exeBaseDir

	for _, part := range parts {
		if part == "" {
			continue
		}
		currentPath = filepath.Join(currentPath, part)
		if err := os.WriteFile(path.Join(currentPath, "index.html"), []byte(""), 0644); err != nil {
			return err
		}
	}
	return nil
}

// 文件上传处理函数
func UploadHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		initConfig()
		// 1. 获取上传文件
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 500,
				"msg":  "无法获取文件: " + err.Error(),
			})
			return
		}
		// 2. 检查文件大小
		if file.Size > maxFileSize {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 500,
				"msg": fmt.Sprintf("文件大小超出限制: %d bytes > %d bytes",
					file.Size, maxFileSize),
			})
			return
		}

		// 3. 检查文件类型
		if !checkFileType(file) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 500,
				"msg":  "不支持的文件类型",
			})
			return
		}

		// 4. 确定上传目录
		uploadDir := getDateDir(assets.UploadsBasePath, "date")

		// 5. 创建目录（包括所有父目录）
		if err := createDirectoryTreeWithIndex(assets.ExecutePath, uploadDir); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "创建目录失败: " + err.Error(),
			})
			return
		}
		// 提取文件扩展名
		fileext := filepath.Ext(file.Filename)

		// 6. 生成文件名
		filename := generateFileName(file.Filename) + fileext

		// 7. 完整文件路径
		filePath := path.Join(uploadDir, filename)

		// 8. 保存文件
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "文件保存失败: " + err.Error(),
			})
			return
		}
		fmt.Println(filePath)
		imgPreviewURL := os.Getenv("SERVER_IMG_PREVIEW_URL")

		// 9. 返回成功响应
		c.JSON(http.StatusOK, gin.H{
			"code":     200,
			"msg":      "文件上传成功",
			"orgName":  file.Filename,
			"filename": filename,
			"size":     file.Size,
			"url":      fmt.Sprintf("%s%s", imgPreviewURL, strings.Replace(filePath, assets.ExecutePath+config.Separator, "/", 1)),
		})
	}
}

// 多文件上传处理函数
func MultiUploadHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		initConfig()
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":  500,
				"error": "表单解析失败: " + err.Error(),
			})
			return
		}

		files := form.File["files"]
		if len(files) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":  500,
				"error": "未找到上传文件",
			})
			return
		}

		var results []gin.H
		var failed []string

		// 确定上传目录
		uploadDir := getDateDir(assets.UploadsBasePath, "date")

		// 创建目录
		if err := createDirectoryTreeWithIndex(assets.ExecutePath, uploadDir); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "创建目录失败: " + err.Error(),
			})
			return
		}

		imgPreviewURL := os.Getenv("SERVER_IMG_PREVIEW_URL")

		for _, file := range files {
			// 检查文件大小
			if file.Size > maxFileSize {
				failed = append(failed, fmt.Sprintf("%s: 文件过大", file.Filename))
				continue
			}

			// 检查文件类型
			if !checkFileType(file) {
				failed = append(failed, fmt.Sprintf("%s: 不支持的文件类型", file.Filename))
				continue
			}
			// 提取文件扩展名
			fileext := filepath.Ext(file.Filename)

			// 生成文件名
			filename := generateFileName(file.Filename) + fileext

			filePath := path.Join(uploadDir, filename)

			// 保存文件
			if err := c.SaveUploadedFile(file, filePath); err != nil {
				failed = append(failed, fmt.Sprintf("%s: 保存失败", file.Filename))
				continue
			}

			results = append(results, gin.H{
				"orgName":  file.Filename,
				"filename": filename,
				"size":     file.Size,
				"url":      fmt.Sprintf("%s%s", imgPreviewURL, strings.Replace(filePath, assets.ExecutePath+config.Separator, "/", 1)),
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"success": len(results),
			"failed":  len(failed),
			"files":   results,
			"errors":  failed,
		})
	}
}
