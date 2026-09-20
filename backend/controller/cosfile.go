package controller

import (
	"fmt"
	"home-box/pkg/storage"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 1. 直接打开 Multipart 文件流 (在内存/缓存区处理，不产生本地磁盘文件)
	srcFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "打开文件流失败"})
		return
	}
	defer srcFile.Close()
	// 2. 拼接 COS 存储路径
	cosPath := "uploads/" + fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	// 3. 直接上传文件流
	fileURL, err := storage.GlobalCOS.UploadStream(c.Request.Context(), cosPath, srcFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "上传成功", "data": gin.H{"fileURL": fileURL}})
}
