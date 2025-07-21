package utils

import (
	"bytes"
	"crypto/aes"    //aes加密算法
	"crypto/cipher" //加密模式
	"crypto/rand"
	"errors"
	"io"
)

// 1.生成AES-CBC
// AES加密要求数据长度为16字节的倍数，用PKCS7Padding 对数据进行填充
func PKCS7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// PKCS7UnPadding 去除填充数据
func PKCS7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("data is empty")
	}
	padding := int(data[length-1])
	if padding > length {
		return nil, errors.New("invalid padding size")
	}
	return data[:length-padding], nil
}

// CBC需要一个16字节的IV向量，generateIV 生成一个随机（IV）
// 确保每次加密结果也不同
func generateIV() ([]byte, error) {
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	return iv, nil
}

// 传参 明文plaintext，key，接受返回均为byte类型
// 2.EncryptAES 使用 AES-CBC 模式加密数据
func EncryptAES(key, originalData []byte) ([]byte, error) {
	//creates and returns a new [cipher.Block].并检查key长度满足16,24,32b
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 填充数据
	originalData = PKCS7Padding(originalData, block.BlockSize())

	// 生成随机 IV
	iv, err := generateIV()
	if err != nil {
		return nil, err
	}

	// 加密数据
	ciphertext := make([]byte, len(originalData))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, originalData)

	// 将 IV 和密文拼接在一起
	return append(iv, ciphertext...), nil
}

// 3.DecryptAES 使用 AES-CBC 模式解密数据
func DecryptAES(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 检查密文长度
	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	// 提取 IV 和密文
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// 解密数据
	originalData := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(originalData, ciphertext)

	// 去除填充
	return PKCS7UnPadding(originalData)
}
