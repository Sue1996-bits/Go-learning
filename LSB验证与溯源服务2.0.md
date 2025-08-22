## 增强 VerificationService 以实现完整溯源流程

（技术提取 + 业务溯源）全部封装在 VerificationService 中。它将成为一个一站式的“溯源中心”。

VerificationService 现在是名副其实的“数字取证中心”，它封装了从像素到用户的完整溯源链路。
纠错算法 (correctErrors) 的加入，使得系统能够抵抗轻度的图像失真，大大增强了水印的稳健性。
清晰的API (/trace/lsb) 和返回结构 (TraceResult)，为上层应用提供了简单易用的接口。

<br>1. 解码图片<br>2. 检查尺寸<br>3. 提取3倍冗余比特流<br>4. 进行多数表决纠错<br>5. 编码为ID字符串<br>6. 查询数据库 watermark_logs 表<br>7. 组合成溯源结果

在数据库 watermark_logs 中存储的 watermark_id 实际上是一个 64个字符长的十六进制字符串。
修复方案就是：在提取出水印后，也进行同样的 Hex 编码，以保持格式一致。


## 问题2：提取水印总是提出与存储水印长度不同
实际数据流：

嵌入时：32字节 × 3份冗余 = 96字节 = 768比特
提取时：只提取256比特 = 32字节

只提取了第一份数据，忽略了后面两份冗余数据！

解决：
```
func (s *verificationServiceImpl) ExtractLSBWatermark(imgData []byte) (string, error) {
	// ❌ 只提取256比特，但实际嵌入了768比特的冗余数据！
	rawWatermark, err := s.lsbProcessor.ExtractWatermark(img, WatermarkPayloadBitLength)
}
// ✅ 正确：提取完整的768比特
longBitStream, _ := s.lsbProcessor.ExtractWatermark(img, TotalBitsToRead)

// ✅ 正确：进行纠错处理
correctedPayload, correctionSuccess, _ := s.correctErrors(longBitStream)
```
```
// 在watermark_service.go中
switch format {
case "jpeg":
    err = jpeg.Encode(&buf, watermarkedImg, &jpeg.Options{Quality: 95}) // ❌ JPEG压缩可能破坏LSB
case "png":
    err = png.Encode(&buf, watermarkedImg) // ✅ PNG无损，适合LSB
}
```
如果必须使用JPEG，质量设置为100，或者考虑使用PNG格式。

结果：长度一致，依旧查询不到。

解决：添加调试日志
格式类似为：
```
func (s *watermarkServiceImpl) ProcessImage(userID, imageID string, imgData []byte, timestamp int64) ([]byte, *models.WatermarkLog, error) {
	// 1. 使用传入的量化时间戳生成确定性的水印ID
	watermarkID, err := s.idGenerator.GenerateID(userID, imageID, timestamp)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate watermark id: %w", err)
	}

	// 🔍 调试日志1：生成的水印ID
	fmt.Printf("🔍 [EMBED DEBUG] Generated watermarkID: %s (length: %d)\n", watermarkID, len(watermarkID))
	fmt.Printf("🔍 [EMBED DEBUG] Input params - userID: %s, imageID: %s, timestamp: %d\n", userID, imageID, timestamp)

	// 【关键】将 Hex 字符串解码回原始的二进制哈希
	rawHash, err := hex.DecodeString(watermarkID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode hex watermark id: %w", err)
	}
	// 🔍 调试日志2：原始哈希
	fmt.Printf("🔍 [EMBED DEBUG] Raw hash: %x (length: %d bytes)\n", rawHash, len(rawHash))

	// 2. 解码图像
	img, format, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode image: %w", err)
	}
	// 🔍 调试日志3：图像格式和尺寸
	bounds := img.Bounds()
	fmt.Printf("🔍 [EMBED DEBUG] Image format: %s, size: %dx%d, total pixels: %d\n",
		format, bounds.Dx(), bounds.Dy(), bounds.Dx()*bounds.Dy())

	// 3. 嵌入水印
	// 【关键修正】创建三倍冗余的水印数据
	// 我们的 watermarkID 是 Hex 字符串，可以直接拼接
	// 注意：这里的 watermarkID 是要存入数据库的【正确ID】
	var redundantPayload []byte
	redundantPayload = make([]byte, 0, len(rawHash)*RedundancyFactor) // 预分配容量
	for i := 0; i < RedundancyFactor; i++ {
		redundantPayload = append(redundantPayload, rawHash...)
	}
	// 🔍 调试日志4：冗余负载长度和内容

	fmt.Printf("🔍 [EMBED DEBUG] Redundant payload length: %d bytes (%d bits)\n",
		len(redundantPayload), len(redundantPayload)*8)
	fmt.Printf("🔍 [EMBED DEBUG] First 16 bytes of redundant payload: %x\n", redundantPayload[:16])

	// 5. 调用【新的、正确的】方法来嵌入二进制负载
	watermarkedImg, err := s.lsbProcessor.EmbedWatermark(img, redundantPayload)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to embed LSB watermark: %w", err)
	}

	// 4. 编码为字节流
	var buf bytes.Buffer
	switch format {
	case "jpeg", "jpg":
		err = png.Encode(&buf, watermarkedImg)
		fmt.Printf("🔍 [EMBED DEBUG] JPEG input converted to PNG output (lossless)\n")
	case "png":
		// PNG输入 -> PNG输出（保持无损）
		err = png.Encode(&buf, watermarkedImg)
		fmt.Printf("🔍 [EMBED DEBUG] PNG input -> PNG output (lossless)\n")
	default:
		// 其他格式 -> PNG输出
		err = png.Encode(&buf, watermarkedImg)
		fmt.Printf("🔍 [EMBED DEBUG] %s input converted to PNG output (lossless)\n", format)
	}
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encode watermarked image: %w", err)
	}
	fmt.Printf("🔍 [EMBED DEBUG] Watermarked image size: %d bytes\n", buf.Len())
	// 🔍 调试日志5：编码后的图像大小,立即验证嵌入效果

	// 5. 创建日志对象并返回
	log := &models.WatermarkLog{
		UserID:      userID,
		ImageID:     imageID,
		WatermarkID: watermarkID,
		CreatedAt:   time.Unix(timestamp, 0),
		UpdatedAt:   time.Unix(timestamp, 0),
	}
	fmt.Printf("🔍 [EMBED DEBUG] Log entry - UserID: %s, ImageID: %s, WatermarkID: %s\n",
		log.UserID, log.ImageID, log.WatermarkID)

	return buf.Bytes(), log, nil
}
```
## 问题2：postman测试：log长度一致，查询false：JPEG压缩破坏了LSB水印

因为 - 三个流的差异：
🔍 [CORRECT DEBUG] Differences - 1vs2: 116, 1vs3: 131, 2vs3: 129
🔍 [CORRECT DEBUG] Corrections made: 188, fully correct: false

768个比特中有188个需要纠错，差异率达到24.5%！ 这表明JPEG压缩严重损坏了LSB数据。

解决：强制使用PNG格式——PNG是无损压缩，完美保护LSB数据
```
	switch format {
	case "jpeg", "jpg":
		// JPEG输入 -> PNG输出（避免压缩损失）
		err = png.Encode(&buf, watermarkedImg)
		fmt.Printf("🔍 [EMBED DEBUG] JPEG input converted to PNG output (lossless)\n")
	case "png":
		// PNG输入 -> PNG输出（保持无损）
		err = png.Encode(&buf, watermarkedImg)
		fmt.Printf("🔍 [EMBED DEBUG] PNG input -> PNG output (lossless)\n")
	default:
		// 其他格式 -> PNG输出
		err = png.Encode(&buf, watermarkedImg)
		fmt.Printf("🔍 [EMBED DEBUG] %s input converted to PNG output (lossless)\n", format)
	}
	
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encode watermarked image: %w", err)
	}
```
最终响应：
json：
{
    "data": {
        "user_id": "user-12345",
        "image_id": "b3eaf7a6-86df-4dcb-a308-c0bcde8eee77",
        "watermark_id": "7ac574df1b4d8340f3b158356aef571043514b9a75474503f37ecb94c7972b64",
        "watermarked_at": "2025-08-22T14:40:00+08:00",
        "extraction_success": true,
        "correction_success": true
    },
    "success": true
}
修复生效了 - PNG格式解决了JPEG压缩损坏LSB数据的问题

水印算法工作正常 - extraction_success和correction_success都为true

数据库查询成功 - 能找到对应的记录并返回完整信息

整个流程从嵌入到提取到溯源都工作正常

完成后一步测试：

1.多格式输入测试

2.不同尺寸测试——- 小图片（100x100）

- 中等图片（500x500）

- 大图片（1000x1000）

- 边界尺寸（正好768像素）

移除多余的debug信息。
