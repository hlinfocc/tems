package controllers

import (
	"net/http"
	"tems-web-api/assets"

	"github.com/gin-gonic/gin"
)

func GetIndexHtml(c *gin.Context) {
	indexHTML, err := assets.GetIndexHtml()
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to read index.html: %v", err)
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTML)
}
