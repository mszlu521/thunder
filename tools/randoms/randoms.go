package randoms

import (
	cr "crypto/rand"
	"fmt"
	"math/rand"
	"time"
)

func Generate4Number() int {
	// 使用当前时间作为种子
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	// 生成一个范围在 1000 到 9999 之间的随机数
	randomNumber := r.Intn(9000) + 1000
	return randomNumber
}
func Generate6Number() int {
	// 使用当前时间作为种子
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	randomNumber := r.Intn(900000) + 100000
	return randomNumber
}

func GenerateTicket() string {
	timestamp := time.Now().UnixNano() // 获取当前时间的纳秒级时间戳
	randomBytes := make([]byte, 16)    // 生成16字节随机数
	_, _ = cr.Read(randomBytes)
	return fmt.Sprintf("%x%x", timestamp, randomBytes) // 拼接时间戳和随机数
}
