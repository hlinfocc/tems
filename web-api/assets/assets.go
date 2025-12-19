package assets

import (
	"embed"
	"net/http"
	"os"
	"path/filepath"
	"tems-web-api/utils"
)

var (
	//go:embed all:dist/*
	content embed.FS

	FileSystem http.FileSystem
	UploadFileSystem http.FileSystem
	UploadsBasePath  string
	ExecutePath      string
)

func init() {
	FileSystem = http.FS(content)
	// 获取当前可执行文件的路径
	exePath, err := os.Executable()
	utils.CheckErrorLog(err)
	// 获取可执行文件所在目录，并拼接 'upload' 子目录
	exeDir, err := filepath.EvalSymlinks(filepath.Dir(exePath))
	utils.CheckErrorLog(err)
	uploadDir := filepath.Join(exeDir, "upload")
	UploadFileSystem = http.Dir(uploadDir)
	UploadsBasePath = uploadDir
	ExecutePath = exeDir
}

func GetIndexHtml() ([]byte, error) {
	indexHTML, err := content.ReadFile("dist/index.html")
	if err != nil {
		// c.String(http.StatusInternalServerError, "Failed to read index.html: %v", err)
		return nil, err
	}
	return indexHTML, nil
}
