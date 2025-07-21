package models

import (
	"gorm.io/gorm"
)

// 模型（Models）与数据库中的数据表进行映射
// 结构体名 (Struct Name) -> 表名 (Table Name)：FileRecord-->file_records
// 驼峰式 (CamelCase) 转换为蛇形 (snake_case)
type EncryptedImage struct {
	gorm.Model              // 内置一个
	OriginalFilename string // 原始文件名
	EncryptedBase64  string // 加密后的文件内容
}
