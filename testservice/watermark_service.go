package services

import (
	"image"
	// "photo_service/models"
)

// WatermarkService 接口定义了服务层对外暴露的核心能力。
// 控制器将通过这个接口与服务层交互。
type WatermarkService interface {

	// 生成水印ID
	GenerateWatermarkID(userID, imageID string, timestamp int64) (string, error)

	// 嵌入LSB水印
	EmbedLSBWatermark(img image.Image, watermarkID string) (image.Image, error)

	// 提取LSB水印
	ExtractLSBWatermark(img image.Image, length int) (string, error)

	// 处理水印图片（完整流程）
	// ProcessWatermarkedImage(userID, imageID string, imgData []byte) ([]byte, *models.WatermarkLog, error)
	ProcessAndGetWatermarkedImage(userID, imageID string) (imageData []byte, contentType string, err error)

	// 验证水印
	VerifyWatermark(imgData []byte, expectedWatermarkID string) (bool, error)
}
