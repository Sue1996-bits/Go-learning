## 增强 VerificationService 以实现完整溯源流程

（技术提取 + 业务溯源）全部封装在 VerificationService 中。它将成为一个一站式的“溯源中心”。

VerificationService 现在是名副其实的“数字取证中心”，它封装了从像素到用户的完整溯源链路。
纠错算法 (correctErrors) 的加入，使得系统能够抵抗轻度的图像失真，大大增强了水印的稳健性。
清晰的API (/trace/lsb) 和返回结构 (TraceResult)，为上层应用提供了简单易用的接口。

<br>1. 解码图片<br>2. 检查尺寸<br>3. 提取3倍冗余比特流<br>4. 进行多数表决纠错<br>5. 编码为ID字符串<br>6. 查询数据库 watermark_logs 表<br>7. 组合成溯源结果

在数据库 watermark_logs 中存储的 watermark_id 实际上是一个 64个字符长的十六进制字符串。
修复方案就是：在提取出水印后，也进行同样的 Hex 编码，以保持格式一致。

