package crypro

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io"
)

// Md5WithSalt 使用 MD5 和 salt 进行哈希
func Md5WithSalt(data, salt string) string {
	// 将 salt 和 data 组合
	combined := data + salt
	// 生成 MD5 哈希
	hash := md5.Sum([]byte(combined))
	// 转为十六进制字符串
	return hex.EncodeToString(hash[:])
}
func Md5(data []byte) string {
	// 生成 MD5 哈希
	hash := md5.Sum(data)
	// 转为十六进制字符串
	return hex.EncodeToString(hash[:])
}

// HashPassword 使用 bcrypt 生成哈希密码
func HashPassword(password string) (string, error) {
	// 使用默认成本值生成哈希
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// CheckPasswordHash 验证密码是否匹配
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Sha1(input string) string {
	hash := sha1.New()
	// 写入数据到哈希器
	hash.Write([]byte(input))
	// 计算哈希值
	hashBytes := hash.Sum(nil)
	// 将哈希值转换为十六进制字符串
	hashString := hex.EncodeToString(hashBytes)
	return hashString
}

func encodeBase64(input []byte) string {
	return base64.URLEncoding.EncodeToString(input)
}

// 使用 URL Safe Base64 解码
func decodeBase64(input string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(input)
}

const Slat = "mscommunity!@#$%"

// EncryptString 加密字符串
func EncryptString(key, plaintext string) (string, error) {
	// 将 key 转换为字节数组
	keyBytes := []byte(key)
	if len(keyBytes) != 16 {
		return "", errors.New("key 长度必须为 16 字节")
	}

	// 创建 AES 加密块
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("创建 AES 块失败: %w", err)
	}

	// 初始化随机化向量 (IV)
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("生成 IV 失败: %w", err)
	}

	// 创建加密器
	stream := cipher.NewCFBEncrypter(block, iv)
	ciphertext := make([]byte, len(plaintext))
	stream.XORKeyStream(ciphertext, []byte(plaintext))

	// 将 IV 和密文组合后返回 Base64 编码
	final := append(iv, ciphertext...)
	return encodeBase64(final), nil
}

// 解密字符串
func DecryptString(key, ciphertext string) (string, error) {
	// 将 key 转换为字节数组
	keyBytes := []byte(key)
	if len(keyBytes) != 16 {
		return "", errors.New("key 长度必须为 16 字节")
	}

	// Base64 解码
	data, err := decodeBase64(ciphertext)
	if err != nil {
		return "", fmt.Errorf("Base64 解码失败: %w", err)
	}

	// 提取 IV 和密文
	if len(data) < aes.BlockSize {
		return "", errors.New("密文长度不足")
	}
	iv := data[:aes.BlockSize]
	ciphertextBytes := data[aes.BlockSize:]

	// 创建 AES 解密块
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("创建 AES 块失败: %w", err)
	}

	// 创建解密器
	stream := cipher.NewCFBDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertextBytes))
	stream.XORKeyStream(plaintext, ciphertextBytes)

	return string(plaintext), nil
}
