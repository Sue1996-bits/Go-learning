**输入一张疑似被泄露的、带LSB水印的图片，输出泄露者的UserID以及泄露时间等关键信息。**

1.**技术阶段：水印提取** —— 从图片像素中恢复出隐藏的二进制信息。


提供一个API端点，它接收一张可能带有LSB水印的图片，然后能够“盲提取”出其中隐藏的水印信息（userID, imageID, timestamp 的哈希值），并将这个信息返回给调用者。

**一、修改 VerificationService：长度永远固定，HMAC-SHA256 的十六进制字符串**

原始：
```
// services/verification_service.go

	// ExtractLSBWatermark 从图片中提取指定长度的 LSB 水印字符串
	ExtractLSBWatermark(imgData []byte, charLength int) (string, error)

func (s *verificationServiceImpl) ExtractLSBWatermark(imgData []byte, charLength int) (string, error) {
	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	// length*8 转换为比特长度
	return s.lsbProcessor.ExtractWatermark(img, charLength*8)
}

```

现：改为固定长度
```
// services/verification_service.go (最终版)

const (
	// HMAC-SHA256 哈希值的十六进制表示是 32字节 * 2 = 64个字符
	LSBWatermarkCharLength = 64
	LSBWatermarkBitLength  = LSBWatermarkCharLength * 8
)

// VerificationService 定义了验证和提取相关的服务接口
type VerificationService interface {
	// ExtractLSBWatermark 从图片中“盲提取”出固定长度的 LSB 水印字符串
	ExtractLSBWatermark(imgData []byte) (string, error)

	// VerifyLSBWatermark 依然保留，用于需要精确比对的场景
	VerifyLSBWatermark(imgData []byte, expectedWatermarkID string) (bool, error)
}

// ExtractLSBWatermark 是我们新增的核心方法
func (s *verificationServiceImpl) ExtractLSBWatermark(imgData []byte) (string, error) {
	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	// 使用预定义的固定比特长度进行提取
	return s.lsbProcessor.ExtractWatermark(img, LSBWatermarkBitLength)
}
```

**二、修改 VerificationController**
新增：
```
// HandleExtract 是处理“盲提取”LSB水印请求的新函数
func (ctrl *VerificationController) HandleExtract(c *gin.Context) {
	// 1. 从请求体中获取图片文件
	file, _, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "image file is required in 'image' field of form-data",
		})
		return
	}
	defer file.Close()

	// 2. 读取文件内容
	imgData, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "failed to read image data from request",
		})
		return
	}

	// 3. 调用服务层执行核心的提取逻辑
	extractedID, err := ctrl.verificationService.ExtractLSBWatermark(imgData)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"success": false,
			"error":   "failed to process image for extraction: " + err.Error(),
		})
		return
	}

	// 4. 返回提取到的水印ID
	// 如果图片中没有水印，提取出的会是一串无意义的乱码，这也是预期的行为
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"extracted_watermark_id": extractedID,
		},
	})
}
```

**三、在 main.go 中注册新路由**
```
// 旧的验证路由，用于精确比对
			adminRoutes.POST("/verify/lsb", verificationController.HandleVerify)
			
			// 【新增】新的提取路由，用于盲提取
			adminRoutes.POST("/extract/lsb", verificationController.HandleExtract)
```

2.**业务阶段：信息匹配与溯源** —— 利用提取出的信息，在我们的日志系统中找到对应的记录，从而锁定泄
