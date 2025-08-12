1.加载配置。

2.初始化数据库连接。

3.按正确的依赖顺序，初始化所有的 Repositories 和 Services。

4.初始化所有的 Controllers，并将它们需要的服务注入进去。

5.设置 Gin 引擎和所有的 API 路由。

6.启动 Web 服务器。

即：
分离的、职责单一的服务（ImageUploadService, ImageService, WatermarkService, VisibleWatermarkService, VerificationService, LogService）。
基于 AES-CBC 的 CryptoService。
基于 go-cache 的通用缓存服务。
为不同业务功能设计的独立 API 端点。


配置驱动: 所有的“魔法字符串”和可变配置（如密钥、路径、端口）都被提取到了文件的顶部，
并且设计为可以从环境变量中读取，
这使得应用可以轻松地在不同环境（开发、测试、生产）中部署。
