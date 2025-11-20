# thunder
基于Gin实现的快开框架，目前正在完善中，正在我的项目中进行使用

## 功能特性

- 基于Gin的Web框架
- 支持多种云存储服务（七牛云、阿里云OSS）
- 数据库访问（GORM + PostgreSQL）
- 配置管理（Viper）
- 日志系统
- 认证授权
- 订阅管理（支持免费版、基础版、高级版、企业版）
- 微信支付集成

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

## 订阅功能

系统支持四种订阅计划：

1. 免费版 (free)
2. 基础版 (basic)
3. 高级版 (pro)
4. 企业版 (enterprise)

支持三种付款时长：

1. 月付 (monthly)
2. 季付 (quarterly)
3. 年付 (yearly)

默认使用微信支付作为付款方式。

### API 接口

- `GET /api/v1/subscriptions/plans` - 获取所有订阅计划配置
- `GET /api/v1/subscriptions/current` - 获取当前用户有效订阅
- `POST /api/v1/subscriptions/` - 创建用户订阅
- `POST /api/v1/subscriptions/wechat-notify` - 微信支付通知回调
