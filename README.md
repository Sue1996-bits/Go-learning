---
---

# 智能水印服务 (Smart Watermark Service)

![Go Version](https://img.shields.io/badge/Go-1.18%2B-blue.svg)
![Framework](https://img.shields.io/badge/Framework-Gin-green.svg)
![Database](https://img.shields.io/badge/Database-SQLite-orange.svg)

一个高性能的 Go 后端服务，提供数字图片的安全加密存储、动态可见水印和不可见 LSB (最低有效位) 隐写水印功能，并包含强大的后台溯源取证能力。

---

## ✨ 项目特性 (Features)

- **安全存储**: 所有上传的图片均通过 AES-CBC 加密后存储，确保静态数据的安全性。
- **动态 LSB 水印**: 为图片嵌入不可见的、基于用户和时间的唯一 LSB 水印，用于泄露追踪和溯源。
- **冗余纠错**: LSB 水印采用三倍冗余嵌入，配合“多数表决”算法，能抵抗轻微的图像压缩和失真。
- **可见文字水印**: 支持在图片上动态添加自定义的可见文字水印。
- **后台溯源**: 提供强大的后台 API，可从图片中“盲提取” LSB 水印信息，并与数据库日志比对，完成溯源。
- **JWT 认证**: 所有核心 API 均通过 JWT 进行保护，确保只有授权用户才能访问。
- **高性能**: 采用 Gin 框架，并为 LSB 水印流程设计了基于时间分片的内存缓存策略，有效降低重复请求的服务器负载。
- **分层架构**: 采用清晰的“控制器-服务-仓库”分层架构，代码高度解耦，易于测试和维护。
- **结构化配置**: 提供完善的配置管理，支持通过环境变量进行灵活部署。

---

## 🛠️ 技术栈 (Tech Stack)

- **语言**: Go (1.18+)
- **Web 框架**: Gin
- **ORM**: GORM
- **数据库**: SQLite (可轻松扩展至 PostgreSQL)
- **加密**: AES-CBC, HMAC-SHA256
- **认证**: JWT (`golang-jwt/jwt/v5`)
- **缓存**: `go-cache` (内存缓存)
- **图像处理**: Go 标准库 (`image`), `fogleman/gg` (用于可见水印)

---

## 🚀 快速开始 (Getting Started)

### 1. 前置条件 (Prerequisites)

- **Go**: 版本 1.18 或更高。
- **C 编译器**: SQLite 驱动需要 CGO，请确保已安装 (macOS/Linux 通常自带，Windows 用户可安装 [TDM-GCC](https://jmeubank.github.io/tdm-gcc/)).
- **Git**

### 2. 克隆项目 (Clone)

```bash
git clone <your-repository-url>
cd <project-folder>
```

### 3. 安装依赖 (Installation)

```bash
go mod tidy
```

### 4. 配置 (Configuration)

本项目通过环境变量或 `.env` 文件进行配置。

1.  复制示例配置文件：
    ```bash
    cp .env.example .env
    ```

2.  打开并编辑 `.env` 文件。对于本地开发，默认配置通常可以直接使用。**在生产环境中，请务必生成并替换所有密钥！**

    ```env
    # 服务器配置
    PORT=8080

    # 数据库配置
    DB_PATH=./data/watermark.db

    # 安全密钥 (生产环境必须替换！)
    AES_KEY=1234567890123456
    JWT_SECRET=a_very_strong_and_long_jwt_secret
    SERVER_SECRET=a_very_strong_hmac_server_secret
    ```

### 5. 运行项目 (Run)

```bash
go run main.go
```

服务器将在 `http://localhost:8080` 启动。现在你可以使用 API 工具 (如 Postman) 开始测试了。

---

## 📖 API 使用指南 (API Usage)

详细的 API 接口文档请参阅 `docs/api-documentation.pdf`。

### 典型工作流

1.  **获取 Token**: `POST /api/v1/auth/token`
2.  **上传图片**: `POST /api/v1/images/upload` (携带 Token)
3.  **获取带 LSB 水印的图片**: `GET /api/v1/images/{image_id}/lsb-watermarked` (携带 Token)
4.  **(溯源)**: `POST /api/v1/admin/trace/lsb` (携带 Token 和泄露的图片)

---

## 🧪 测试 (Testing)

(可选) 如果项目包含测试文件，可以运行：
```bash
go test ./...
```

---

## 🤝 贡献 (Contributing)

欢迎提交问题 (Issues) 和合并请求 (Pull Requests)。


---
