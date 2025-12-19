package utils

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func String2Int(s string) int {
	if len(s) == 0 {
		return -1
	}
	res, err := strconv.Atoi(s)
	if err != nil {
		return -1
	}
	return res
}

func RandStringBytes(n int, letter bool) string {
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	if letter {
		letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[r.Intn(len(letterBytes))]
	}
	return string(b)
}

func CheckPortStatus(port int) bool {
	// 监听 端口
	log.Println("check listener Port", port)
	listenerPort := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", listenerPort)
	if err != nil {
		// 如果监听失败，则说明端口已被占用
		return false
	}
	// 关闭监听器
	defer listener.Close()

	// 如果监听成功，则说明端口未被占用
	return true
}

// 判断是否注册了/static
func HasHomeStaticRoute(router *gin.Engine, prefix string) bool {
	for _, route := range router.Routes() {
		// StaticFS 会注册 GET 和 HEAD 方法，路径模式会是 "/static/*filepath"
		if strings.HasPrefix(route.Path, prefix) {
			return true
		}
	}
	return false
}

func HasPrefixIgnoreCase(s, prefix string) bool {
	if len(s) < len(prefix) {
		return false
	}
	return strings.EqualFold(s[:len(prefix)], prefix)
}

func Sha3(data string) string {
	h := sha512.New()
	h.Write([]byte(data))
	hash := h.Sum(nil)
	return hex.EncodeToString(hash)
}

func CheckErrorLog(err error) {
	if err != nil {
		log.Fatalf("无法获取可执行文件路径: %v", err.Error())
		os.Exit(1)
	}
}