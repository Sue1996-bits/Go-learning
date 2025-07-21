Phase 0: 环境准备与安全重构

安全修复： 将所有硬编码的密钥（AES Key）改为从环境变量读取。
日志系统引入： 引入slog，替换掉项目中所有的fmt.Println和log，建立结构化日志标准。
项目结构优化： 创建services/目录和middlewares/目录，为后续开发做好准备。
配置管理： 创建一个config包，用于统一加载和管理所有配置（环境变量、默认值等）。

1.数据库迁移：
/db/db.go —— 数据库初始化：连接、配置、迁移、索引

/db/migration.go —— 执行数据库结构迁移与索引创建(sql)

/db/watermark_repository.go —— WatermarkLog 表的数据访问层（Repository 模式）,封装具体 CRUD 操作（抽象数据库访问）

### 业务逻辑依赖接口!!!!!!

` 接口（Interface）+ 实现（Struct） 的设计模式`

` 接口封装了所有对表格[WatermarkLog]的操作 + 解耦业务逻辑和数据库实现，便于测试和扩展`
```

// 接口定义包含方法
type WatermarkRepository interface {
	Create(log *models.WatermarkLog) error
	GetByUserAndImage(userID, imageID string) (*models.WatermarkLog, error)
	GetByWatermarkID(watermarkID string) (*models.WatermarkLog, error)
	GetUserLogs(userID string, limit, offset int) ([]*models.WatermarkLog, error)
	GetImageLogs(imageID string, limit, offset int) ([]*models.WatermarkLog, error)
}
// 定义一个结构体，通过这个结构体实现接口
type watermarkRepository struct {
	db *gorm.DB
}

// 出生了）构造函数
// ！返回的是接口 WatermarkRepository，而不是实现结构体 *watermarkRepository
// 这样让上层代码 只依赖抽象，不依赖具体实现。
func NewWatermarkRepository(db *gorm.DB) WatermarkRepository {
	return &watermarkRepository{db: db}
}
```
` Go 语言中的接口设计思想：面向接口编程（program to an interface）`

！这里可以不用实现数据库即可测试。

mock（区别于watermarkRepository使用 *gorm.DB 实现接口；定义一个新的 struct，满足接口中的所有方法签名。）如：
```
type MockWatermarkRepo struct{}

func (m *MockWatermarkRepo) SaveLog(...) error {
    return nil
}

··

repo := &MockWatermarkRepo{}
//方法的接受者类型为*Type，就必须使用指针去调用它。
//接口绑定时不关心=值/指针，只要实现了所有方法：
HandleSave(repo)
```
如果返回具体结构体，则可直接访问 repo.db，这样会打破封装、耦合实现
`代码难以更换、复用、测试和维护。`

`使用接口作为返回值，是为了“隔离实现”，实现高内聚低耦合，遵循 Go 的最佳实践 —— 面向接口编程。`
类比java：
```
interface Repository { void save(); }
class MyRepository implements Repository { ... }
Repository repo = new MyRepository();
```
next:--对接 slog 的 repository 调用日志

