package randoms

import (
	cr "crypto/rand"
	"fmt"
	"math/rand"
	v2 "math/rand/v2"
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

// Gen6Code 生成6位数字验证码
func Gen6Code() (string, error) {
	// 1. 准备种子
	var seed [32]byte
	if _, err := cr.Read(seed[:]); err != nil {
		return "", err
	}
	// 2. 创建随机源
	cha8 := v2.NewChaCha8(seed)
	// 3. 基于源创建一个随机生成器对象
	r := v2.New(cha8)
	// 4. 直接请求一个 [0, 1000000) 的随机整数
	// IntN 会自动处理边界，比自己做取模运算更均匀（没有模偏差）
	num := r.IntN(1000000)
	// 5. 格式化为 6 位字符串 (不足补0)
	return fmt.Sprintf("%06d", num), nil
}

// Gen6CodeNumber 生成6位数字验证码
func Gen6CodeNumber() (int, error) {
	// 1. 准备种子
	var seed [32]byte
	if _, err := cr.Read(seed[:]); err != nil {
		return 0, err
	}
	// 2. 创建随机源
	cha8 := v2.NewChaCha8(seed)
	// 3. 基于源创建一个随机生成器对象
	r := v2.New(cha8)
	// 4. 直接请求一个 [0, 1000000) 的随机整数
	// IntN 会自动处理边界，比自己做取模运算更均匀（没有模偏差）
	num := r.IntN(1000000)
	// 5. 格式化为 6 位字符串 (不足补0)
	return num, nil
}