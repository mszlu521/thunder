package upload

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/mszlu521/thunder/logs"
)


type AliyunOSSUpload struct {
	client *oss.Client
	bucket *oss.Bucket
}

// InitAliyunOSSUpload 初始化阿里云OSS上传管理器
func InitAliyunOSSUpload(accessKeyID, accessKeySecret, endpoint, bucketName string) (*AliyunOSSUpload, error) {
	// 检查必要配置是否存在
	if accessKeyID == "" || accessKeySecret == "" || endpoint == "" || bucketName == "" {
		logs.Info("阿里云OSS配置不完整，跳过初始化")
		return nil,errors.New("阿里云OSS配置不完整")
	}
	
	// 创建OSSClient实例
	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		logs.Errorf("创建阿里云OSS客户端失败: %v", err)
		return nil, err
	}
	
	// 获取存储空间
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		logs.Errorf("获取阿里云OSS存储空间失败: %v", err)
		return nil, err
	}
	
	aliyunOSSUploadManager := &AliyunOSSUpload{
		client: client,
		bucket: bucket,
	}
	
	logs.Info("阿里云OSS上传服务初始化成功")
	return aliyunOSSUploadManager, nil
}

// Upload 上传文件到阿里云OSS
func (a *AliyunOSSUpload) Upload(ctx context.Context, reader io.Reader, name string) error {
	if a == nil || a.bucket == nil {
		return errors.New("no such bucket")
	}
	
	// 直接使用PutObject上传
	err := a.bucket.PutObject(name, reader)
	if err != nil {
		logs.Errorf("上传文件到阿里云OSS失败: %v", err)
		return err
	}
	
	return nil
}

// UploadWithMetadata 上传文件到阿里云OSS并附带元数据
func (a *AliyunOSSUpload) UploadWithMetadata(ctx context.Context, reader io.Reader, name string, metadata map[string]string) error {
	if a == nil || a.bucket == nil {
		return errors.New("no such bucket")
	}
	
	// 设置对象元数据
	var options []oss.Option
	for key, value := range metadata {
		options = append(options, oss.Meta(key, value))
	}
	
	// 上传对象
	err := a.bucket.PutObject(name, reader, options...)
	if err != nil {
		logs.Errorf("上传文件到阿里云OSS失败: %v", err)
		return err
	}
	
	return nil
}

// GetSignedURL 获取对象的签名URL，用于临时访问
func (a *AliyunOSSUpload) GetSignedURL(objectKey string, expiredInSec int64) (string, error) {
	if a == nil || a.bucket == nil {
		return "", errors.New("no such bucket")
	}
	
	// 生成签名URL
	signedURL, err := a.bucket.SignURL(objectKey, oss.HTTPGet, expiredInSec)
	if err != nil {
		logs.Errorf("生成阿里云OSS签名URL失败: %v", err)
		return "", err
	}
	return signedURL, nil
}

// DeleteObject 删除OSS中的对象
func (a *AliyunOSSUpload) DeleteObject(objectKey string) error {
	if a == nil || a.bucket == nil {
		return errors.New("no such bucket")
	}
	
	err := a.bucket.DeleteObject(objectKey)
	if err != nil {
		logs.Errorf("删除阿里云OSS对象失败: %v", err)
		return err
	}
	
	return nil
}

// IsAvailable 检查阿里云OSS服务是否可用
func (a *AliyunOSSUpload) IsAvailable() bool {
	return a != nil && a.bucket != nil
}

// GetObjectURL 获取对象的访问URL
func (a *AliyunOSSUpload) GetObjectURL(endpoint, bucketName, objectKey string) string {
	if a == nil || a.bucket == nil {
		return ""
	}
	
	// 处理endpoint，确保格式正确
	if !strings.HasPrefix(endpoint, "http") {
		endpoint = "https://" + endpoint
	}
	
	// 构造对象URL
	objectURL := endpoint + "/" + bucketName + "/" + objectKey
	return objectURL
}

func (a *AliyunOSSUpload) GetPublicUrl(bucket, endpoint, filename string) string {
	return fmt.Sprintf("https://%s.%s/%s", bucket, endpoint, filename)
}