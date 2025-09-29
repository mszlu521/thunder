# thunder
基于Gin实现的快开框架，目前正在完善中，正在我的项目中进行使用

## 功能特性

- 基于Gin的Web框架
- 支持多种云存储服务（七牛云、阿里云OSS）
- 数据库访问（GORM + PostgreSQL）
- 配置管理（Viper）
- 日志系统
- 认证授权

## 云存储配置

### 七牛云配置

在 `app/etc/config.yml` 中配置七牛云参数：

```yaml
qiniu:
  bucket: "your-bucket-name"
  accessKey: "your-access-key"
  secretKey: "your-secret-key"
  region: "z0"  # 存储区域
```

### 阿里云OSS配置

在 `app/etc/config.yml` 中配置阿里云OSS参数：

```yaml
aliyun:
  accessKeyId: "your-access-key-id"
  accessKeySecret: "your-access-key-secret"
  endpoint: "oss-cn-hangzhou.aliyuncs.com"  # 根据实际区域选择
  bucket: "your-bucket-name"
```

## 使用方法

### 上传文件到阿里云OSS

```go
import "github.com/mszlu521/thunder/upload"

// 检查服务是否可用
if upload.AliyunOSSUploadManager.IsAvailable() {
    // 上传文件
    err := upload.AliyunOSSUploadManager.Upload(context.Background(), fileReader, "path/to/file")
    if err != nil {
        // 处理错误
    }
}
```

### 获取文件访问URL

```go
// 获取文件的访问URL
url := upload.AliyunOSSUploadManager.GetObjectURL("path/to/file")
```

### 生成签名URL

```go
// 生成临时访问的签名URL（有效期3600秒）
signedURL, err := upload.AliyunOSSUploadManager.GetSignedURL("path/to/file", 3600)
```