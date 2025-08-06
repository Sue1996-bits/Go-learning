// services/watermark_service_impl.go
package testservices

import (
	"bytes"
	"fmt"
	"image"
	"photo_service/db"
	"photo_service/models"
	"time"
)

// watermarkServiceImpl 是 WatermarkService 接口的具体实现。
// 它包含了实现业务逻辑所需的所有依赖。

type watermarkServiceImpl struct {
	idGenerator   *WatermarkIDGenerator
	lsbProcessor  *LSBWatermarkProcessor
	cacheService  CacheService
	cacheKeyGen   *CacheKeyGenerator
	watermarkRepo db.WatermarkRepository
	imageRepo     db.EncryptedImageRepository
	cryptoService CryptoService // <-- 新增
}

// NewWatermarkService 是 watermarkServiceImpl 的构造函数。
// 它遵循依赖注入（DI）原则，将所有必需的组件作为参数传入，
// 然后初始化并返回一个实现了 WatermarkService 接口的实例。

func NewWatermarkService(
	serverSecret string,
	cacheService CacheService,
	watermarkRepo db.WatermarkRepository,
	imageRepo db.EncryptedImageRepository,
	cryptoSvc CryptoService, // <-- 新增
	idGen *WatermarkIDGenerator,
	lsbProc *LSBWatermarkProcessor,
	cacheKeyGen *CacheKeyGenerator,
) WatermarkService {
	return &watermarkServiceImpl{

		// idGenerator:   NewWatermarkIDGenerator(serverSecret),
		// lsbProcessor:  NewLSBWatermarkProcessor("blue"),
		// cacheKeyGen:   NewCacheKeyGenerator("watermark"),
		// 3个是【内部创建】的

		idGenerator:  idGen,
		lsbProcessor: lsbProc,
		cacheKeyGen:  cacheKeyGen,

		cacheService:  cacheService,
		watermarkRepo: watermarkRepo,
		imageRepo:     imageRepo,
		cryptoService: cryptoSvc,
		// 4个是【外部注入】的
	}
}

func (s *watermarkServiceImpl) GenerateWatermarkID(userID, imageID string, timestamp int64) (string, error) {
	return s.idGenerator.GenerateID(userID, imageID, timestamp)
}

func (s *watermarkServiceImpl) EmbedLSBWatermark(img image.Image, watermarkID string) (image.Image, error) {
	return s.lsbProcessor.EmbedWatermark(img, watermarkID)
}

func (s *watermarkServiceImpl) ExtractLSBWatermark(img image.Image, length int) (string, error) {
	return s.lsbProcessor.ExtractWatermark(img, length*8) // 转换为比特长度
}

// func (s *watermarkServiceImpl) ProcessWatermarkedImage(userID, imageID string, imgData []byte) ([]byte, *models.WatermarkLog, error) {
// 	// 1. 检查缓存
// 	cacheKey := s.cacheKeyGen.WatermarkKey(userID, imageID)
// 	if cachedData, found := s.cacheService.Get(cacheKey); found {
// 		// 从数据库获取日志记录
// 		if log, err := s.watermarkRepo.GetByUserAndImage(userID, imageID); err == nil {
// 			return cachedData, log, nil
// 		}
// 	}

// 	// 2. 生成水印ID
// 	timestamp := time.Now().Unix()
// 	watermarkID, err := s.GenerateWatermarkID(userID, imageID, timestamp)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	// 3. 解码图像
// 	img, format, err := image.Decode(bytes.NewReader(imgData))
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	// 4. 嵌入水印
// 	watermarkedImg, err := s.EmbedLSBWatermark(img, watermarkID)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	// 5. 编码图像
// 	var buf bytes.Buffer
// 	switch format {
// 	case "jpeg":
// 		err = jpeg.Encode(&buf, watermarkedImg, &jpeg.Options{Quality: 95})
// 	case "png":
// 		err = png.Encode(&buf, watermarkedImg)
// 	default:
// 		err = errors.New("unsupported image format")
// 	}

// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	watermarkedData := buf.Bytes()

// 	// 6. 写入缓存
// 	s.cacheService.Set(cacheKey, watermarkedData, time.Hour)

// 	// 7. 记录日志
// 	log := &models.WatermarkLog{
// 		UserID:      userID,
// 		ImageID:     imageID,
// 		WatermarkID: watermarkID,
// 		CreatedAt:   time.Now(),
// 		UpdatedAt:   time.Now(),
// 	}

// 	if err := s.watermarkRepo.Create(log); err != nil {
// 		return nil, nil, err
// 	}

// 	return watermarkedData, log, nil
// }

// 2.0
const cacheInterval = 5 * time.Minute // 定义缓存和时间分片的窗口期

func (s *watermarkServiceImpl) ProcessAndGetWatermarkedImage(userID, imageID string) ([]byte, string, error) {
	// 1. 获取量化后的时间戳
	now := time.Now()
	truncatedTime := TruncateTime(now, cacheInterval)
	timestamp := truncatedTime.Unix()

	// 2. 生成缓存键
	cacheKey := s.cacheKeyGen.WatermarkKey(userID, imageID, timestamp)

	// 3. 检查缓存
	if cachedData, found := s.cacheService.Get(cacheKey); found {
		fmt.Printf("Watermarked image HIT cache for key: %s\n", cacheKey)
		// 缓存命中，直接返回。同样需要查询一次元数据获取 contentType
		record, err := s.imageRepo.GetByID(imageID)
		if err != nil {
			return nil, "", err
		}
		return cachedData, record.ContentType, nil
	}

	fmt.Printf("Watermarked image MISS cache for key: %s\n", cacheKey)

	// --- 缓存未命中，执行完整流程 ---

	// 4. 从数据库获取并解密
	encryptedImage, err := s.imageRepo.GetByID(imageID)
	if err != nil {
		return nil, "", err
	}
	plaintextData, err := s.cryptoService.Decrypt(encryptedImage.EncryptedData)
	if err != nil {
		return nil, "", err
	}

	// 5. 使用【量化后的时间戳】生成 LSB 水印ID
	watermarkID, err := s.idGenerator.GenerateID(userID, imageID, timestamp)
	if err != nil {
		return nil, "", err
	}

	// 6. 解码图片
	img, format, err := image.Decode(bytes.NewReader(plaintextData))
	if err != nil {
		return nil, "", err
	}

	// 7. 嵌入 LSB 水印
	watermarkedImg, err := s.lsbProcessor.EmbedWatermark(img, watermarkID)
	if err != nil {
		return nil, "", err
	}

	// 8. 将处理后的图片编码为二进制流
	var buf bytes.Buffer
	// ... (编码逻辑不变) ...
	if err != nil {
		return nil, "", err
	}
	finalImageData := buf.Bytes()

	// 9. 将【带 LSB 水印的明文图片】存入缓存
	// TTL 设置为 cacheInterval，确保它在下一个时间窗口开始时失效
	s.cacheService.Set(cacheKey, finalImageData, cacheInterval)

	// 10. 记录日志
	log := &models.WatermarkLog{
		UserID:      userID,
		ImageID:     imageID,
		WatermarkID: watermarkID, // 使用的是包含量化时间戳的ID
	}
	go s.watermarkRepo.Create(log)

	return finalImageData, encryptedImage.ContentType, nil
}

// 将当前时间“取整”。
func TruncateTime(t time.Time, interval time.Duration) time.Time {
	return t.Truncate(interval)
}

func (s *watermarkServiceImpl) VerifyWatermark(imgData []byte, expectedWatermarkID string) (bool, error) {
	// 解码图像
	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return false, err
	}

	// 提取水印
	extractedID, err := s.ExtractLSBWatermark(img, len(expectedWatermarkID))
	if err != nil {
		return false, err
	}

	return extractedID == expectedWatermarkID, nil
}
