package controllers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"photo_service/db"
	"photo_service/models"
	"photo_service/utils"

	"github.com/gin-gonic/gin"
)

// 密钥：min：16字节
// 最终输出的数据长度与密钥长度无关
var key = []byte("0123456789abcdef")

func EncryptImage(ctx *gin.Context) {

	// 获取请求中传递的图片文件名
	imageName := ctx.Query("filename")
	if imageName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "未提供图片文件名"})
		return
	}

	// 构造图片文件的完整路径
	imagePath := filepath.Join("storage/uploads", imageName)
	fmt.Println("文件路径:", imagePath)

	// 检查文件是否存在
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	// 读取图片文件的“完整原始二进制”内容
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("读取图片失败: %v", err)})
		return
	}
	// 构造要加密的数据（用 "::" 分隔，后面解密时能提取 filename）
	plaintext := append([]byte(imageName+"::"), imageData...)

	//加密
	encryptedData, err := utils.EncryptAES(key, plaintext)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "加密失败"})
		return
	}
	// 编码为 base64，便于前端传输和存储
	encoded := base64.StdEncoding.EncodeToString(encryptedData)

	record := models.EncryptedImage{
		OriginalFilename: imageName,
		EncryptedBase64:  encoded,
	}

	// 使用DB.Create 存储文件名和密文
	if err := db.DB.Create(&record).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "保存到数据库失败"})
		return
	}

	// 返回加密结果
	ctx.JSON(200, gin.H{
		"success":  "图片加密保存成功",
		"id":       record.ID,
		"filename": record.OriginalFilename,
	})

}
