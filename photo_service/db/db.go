import (
	"log"
	"photo_service/models"

	"gorm.io/driver/sqlite"

	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	var err error
	// gorm.Open 函数来打开一个 SQLite 数据库文件,并传入一个 gorm.Config 结构体来配置一些选项。
	DB, err = gorm.Open(sqlite.Open("photo.db"), &gorm.Config{}) //!!!注意这里不是:= 是= 即赋值而非新建！！！
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	// 自动建表:db.AutoMigrate 方法来根据结构体定义自动创建数据表。
	if err := DB.AutoMigrate(&models.EncryptedImage{}); err != nil {
		log.Fatal("自动迁移失败:", err)
	}
	log.Println("✅ 数据库初始化完成")

}
