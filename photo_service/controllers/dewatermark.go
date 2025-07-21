package controllers

import (
	"bytes"
	"encoding/base64"
	"image"
	"mime"
	"path/filepath"

	// "mime"
	// "net/http"
	// "path/filepath"
	"photo_service/db"
	"photo_service/models"
	"photo_service/utils"

	// "time"

	"fmt"

	//仅支持以下格式
	_ "image/gif"
	_ "image/jpeg"
	"image/png"

	"github.com/fogleman/gg"
	"github.com/gin-gonic/gin"
)

func DecryptImage2(id string, key []byte) ([]byte, string, error) {

	var record models.EncryptedImage //定义一个结构体接受返回数据

	if err := db.DB.First(&record, id).Error; err != nil {
		return nil, "", fmt.Errorf("记录未找到 (id: %s): %w", id, err)
	}

	//base64-->byte
	decodedBytes, err := base64.StdEncoding.DecodeString(record.EncryptedBase64)
	if err != nil {
		return nil, "", fmt.Errorf("Base64解码失败: %w", err)
	}

	//解密
	decryptedText, err := utils.DecryptAES(key, decodedBytes)
	if err != nil {
		return nil, "", fmt.Errorf("AES解密失败: %w", err)
	}
	//分割
	parts := bytes.SplitN(decryptedText, []byte("::"), 2) //数据库中没有存文件名最好用这个
	if len(parts) != 2 {
		return nil, "", fmt.Errorf("解密后的数据格式无效")
	}

	decryptedData := parts[1]
	filename := string(parts[0])

	return decryptedData, filename, nil

}

func AddWatermark(img image.Image, watermarkText string) (image.Image, error) {
	// image.Image 只读图片类

	w := img.Bounds().Dx()
	h := img.Bounds().Dy()

	// NewContext 创建一个具有指定宽度和高度的新图像 RGBA，将image.Image只读的东西先画上去
	dc := gg.NewContext(w, h)
	// 将文字内容渲染到图片。
	// font = 字体
	dc.DrawImage(img, 0, 0)

	//1.动态计算合适的字体大小（从大到小试探）
	// 目标宽度：图片宽度的 25%
	// 目标高度：图片高度的 5%
	maxWidth := float64(w) * 0.25
	maxHeight := float64(h) * 0.05

	// 默认起始字体大小
	fontSize := maxHeight
	fontPath := "fonts/font.ttf"

	// 尝试逐步减小字体直到合适
	for fontSize > 5 {
		if err := dc.LoadFontFace(fontPath, fontSize); err != nil {
			return nil, fmt.Errorf("无法加载字体: %v", err)
		}
		textWidth, textHeight := dc.MeasureString(watermarkText)
		if textWidth <= maxWidth && textHeight <= maxHeight {
			break
		}
		fontSize -= 1
	}

	// 2.设置水印颜色为半透明白色 + 绘制
	dc.SetRGBA(1, 1, 1, 0.6)

	// 设置水印位置（右下角）
	margin := 20.0
	textWidth, textHeight := dc.MeasureString(watermarkText)
	x := float64(w) - textWidth - margin
	y := float64(h) - margin - textHeight/8

	dc.DrawString(watermarkText, x, y+textHeight/2)
	return dc.Image(), nil

}

func DecryptAndWatermark(ctx *gin.Context) {
	id := ctx.Query("id")
	data, filename, err := DecryptImage2(id, key)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "图片解密失败"})
		return
	}
	fmt.Println("图片解密成功")

	// 图片解码
	img, formatName, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		ctx.JSON(500, gin.H{"error": "图片解码失败"})
		return
	}
	fmt.Printf("成功解码图片，格式为: %s\n", formatName)

	watermarkedImg, err := AddWatermark(img, "已加密。")
	if err != nil {
		ctx.JSON(500, gin.H{"error": "图片加水印失败"})
		return
	}

	// ctx.Header("Content-Type", "image/jpeg")
	contentType := mime.TypeByExtension(filepath.Ext(filename))
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	ctx.Header("Content-Type", contentType) //通知客户端如何处理接受的数据
	// writer 方法直接可以读img类型。
	png.Encode(ctx.Writer, watermarkedImg)
	//io.Writer
}
