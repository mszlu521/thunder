package httputils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPClient 定义一个 HTTP 客户端工具类
type HTTPClient struct {
	client *http.Client
}

// NewHTTPClient 创建一个新的 HTTP 客户端
func NewHTTPClient(timeout time.Duration) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{Timeout: timeout},
	}
}

// GET 发送 GET 请求
func (h *HTTPClient) GET(url string, headers map[string]string) (string, error) {
	// 创建请求
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("创建 GET 请求失败: %w", err)
	}

	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 执行请求
	resp, err := h.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送 GET 请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应体失败: %w", err)
	}

	return string(body), nil
}

// POST 发送 POST 请求（支持 JSON 数据）
func (h *HTTPClient) POST(url string, headers map[string]string, data interface{}) (string, error) {
	// 将数据编码为 JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("JSON 编码失败: %w", err)
	}

	// 创建请求
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建 POST 请求失败: %w", err)
	}

	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	req.Header.Set("Content-Type", "application/json")

	// 执行请求
	resp, err := h.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送 POST 请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应体失败: %w", err)
	}

	return string(body), nil
}
