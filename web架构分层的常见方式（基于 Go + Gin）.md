- controllers/   # 控制器层，处理 HTTP 请求和响应
- services/      # 服务层，封装业务逻辑
- models/        # 数据结构定义与数据库操作（ORM）
- db/            # 数据库连接与迁移初始化
- utils/         # 工具函数（如加密、Token 等）
- middlewares/   # Gin 中间件（如 JWT 验证、日志）
- config/        # 配置加载
- main.go        # 程序入口

| 层级         | 职责                  | 举例                                 |
| ---------- | ------------------- | ---------------------------------- |
| Controller | 接收请求、解析参数、调用服务、返回响应 | `BindJSON()` / `c.JSON(...)`       |
| Service    | 实现业务逻辑，组织调用模型或工具    | `CheckPassword()` / `SaveUser()`   |
| Model      | 数据库映射与基础数据操作        | `db.First(...)` / `db.Create(...)` |


| 优点   | 说明                                   |
| ---- | ------------------------------------ |
| 解耦   | controller 不依赖数据库，不操作密码加密逻辑          |
| 可测试  | 可以只测试服务层，或 mock service 测 controller |
| 可复用  | 多个 controller 可以共用同一个 service        |
| 易维护  | 变更逻辑只需修改 service，controller 不动       |
| 面向接口 | 后期 service 可实现接口，便于 mock 和扩展         |

### Controller 管请求，Service 管逻辑，Model 管数据。
