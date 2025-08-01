package controllers

import (
	"bytes"
	"encoding/base64"
	"mime"
	"net/http"
	"path/filepath"
	"photo_service/db"
	"photo_service/models"
	"photo_service/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func DecryptImage(ctx *gin.Context) {
	//读取url中的加密数据
	//加密数据太大，所以这里是用id读取sqlite中的EncryptedBase64，

	id := ctx.Query("id")

	var record models.EncryptedImage //定义一个结构体接受返回数据

	//go中在函数内部对调用的参数的修改 不会影响原始变量，因为接受到的拷贝
	//&record = dest interface{} (目标对象)   id = conds ...interface{} (查询条件)

	//GORM 在后台生成并执行SQL: SELECT * FROM `file_records` WHERE `file_records`.`id` = 10 LIMIT 1;
	// 发送给 SQLite 数据库引擎去执行。
	if err := db.DB.First(&record, id).Error; err != nil {
		ctx.JSON(404, gin.H{"error": "记录未找到"})
		return
	}
	//if err := ...; err != nil 结构化地同时执行 操作 和 处理异常

	//base64-->byte
	decodedBytes, err := base64.StdEncoding.DecodeString(record.EncryptedBase64)
	if err != nil {
		ctx.JSON(404, gin.H{"error": "Base64 解码失败"})
		return
	}

	//解密
	decryptedText, err := utils.DecryptAES(key, decodedBytes)
	if err != nil {
		ctx.JSON(404, gin.H{"error": "解密失败"})
		return
	}
	//分割
	parts := bytes.SplitN(decryptedText, []byte("::"), 2) //数据库中没有存文件名最好用这个
	if len(parts) != 2 {
		// 错误处理：格式不对
		ctx.JSON(400, gin.H{"error": "解密数据格式错误"})
		return
	}

	decryptedData := parts[1]
	filename := string(parts[0])

	// return decryptedData, filename

	//图片二进制内容（字节 []byte）如何恢复为可访问的本地图片文件，用于查看、展示或下载？

	//动态读取
	contentType := mime.TypeByExtension(filepath.Ext(filename))
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	//响应头
	ctx.Header("Content-Type", contentType)
	ctx.Header("Content-Disposition", "inline; filename=\""+filename+"\"")
	//inline 显示    attachment 弹出对话框下载
	ctx.Header("Accept-Ranges", "bytes") //启用 支持断点续传、分段下载（Range 请求）
	// 判断是否HEAD 请求
	if ctx.Request.Method == http.MethodHead {
		ctx.Status(200)
		return
	}
	// 支持 Range 请求，这里才返回实际内容
	http.ServeContent(ctx.Writer, ctx.Request, filename, time.Now(), bytes.NewReader(decryptedData))

}
