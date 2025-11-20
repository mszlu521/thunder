# Thunder

[![Go Report Card](https://goreportcard.com/badge/github.com/mszlu521/thunder)](https://goreportcard.com/report/github.com/mszlu521/thunder)
[![GoDoc](https://pkg.go.dev/badge/github.com/mszlu521/thunder)](https://pkg.go.dev/github.com/mszlu521/thunder)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Thunder is a fast and lightweight web framework built on top of Gin, designed to accelerate Go web application development. It provides out-of-the-box solutions for common requirements like authentication, cloud storage, database access, payment integration, and more.

## Features

- 🔥 **Gin-based**: Built on top of the popular Gin web framework for high performance
- ☁️ **Cloud Storage**: Support for multiple cloud storage services (Qiniu Cloud, Alibaba Cloud OSS)
- 🗄️ **Database Access**: GORM integration with PostgreSQL and MySQL support
- ⚙️ **Configuration Management**: Viper-based configuration system with hot reloading
- 📝 **Logging**: Structured logging with customizable output formats
- 🔐 **Authentication & Authorization**: JWT-based authentication middleware
- 💰 **Payment Integration**: WeChat Pay integration for native, JSAPI, and H5 payments
- 🔄 **Event System**: Built-in event management system
- 🛡️ **Security**: CORS middleware, request validation
- 📤 **File Upload**: Easy integration with cloud storage services
- 📋 **Subscription Management**: Built-in subscription system with plan support (free, basic, pro, enterprise)

## Installation

```bash
go get github.com/mszlu521/thunder
```

## Quick Start

```go
package main

import (
    "github.com/mszlu521/thunder/config"
    "github.com/mszlu521/thunder/server"
)

func main() {
    // Initialize configuration
    conf := config.Init()
    config := config.GetConfig()
    
    // Create server instance
    s := server.NewServer(config)
    
    // Start server
    s.Start()
}
```

## Configuration

Thunder uses YAML configuration files. Create a `config.yml` file in your `etc` directory:

```yaml
server:
  mode: "release"
  host: "127.0.0.1"
  port: 8080
  readTimeout: "5s"
  writeTimeout: "5s"

log:
  level: "info"
  format: "json"
  addSource: false

jwt:
  secret: "your-jwt-secret"
  expire: "24h"

db:
  postgres:
    host: "localhost"
    port: 5432
    user: "postgres"
    password: "password"
    database: "mydb"
    sslmode: "disable"
    maxIdleConns: 10
    maxOpenConns: 100
    pingTimeout: "5s"
    
  redis:
    addr: "localhost:6379"
    password: ""
    db: 0
    poolSize: 100
    maxIdleConns: 10
    maxOpenConns: 1000

auth:
  isAuth: true
  ignores:
    - "/api/v1/public/**"
    - "/health"
  needLogins:
    - "/api/v1/user/**"

upload:
  prefix: "/uploads"

qiniu:
  bucket: "your-qiniu-bucket"
  accessKey: "your-qiniu-access-key"
  secretKey: "your-qiniu-secret-key"
  region: "z0"

aliyun:
  accessKeyId: "your-aliyun-access-key-id"
  accessKeySecret: "your-aliyun-access-key-secret"
  endpoint: "oss-cn-hangzhou.aliyuncs.com"
  bucket: "your-oss-bucket"

pay:
  wxPay:
    appId: "your-wechat-app-id"
    mchId: "your-merchant-id"
    mchSerialNo: "your-merchant-serial-no"
    apiV3Key: "your-api-v3-key"
    privateKey: "your-private-key"
    appSecret: "your-app-secret"
    notifyUrl: "https://yourdomain.com/pay/notify"
```

## Cloud Storage

### Qiniu Cloud

```go
import "github.com/mszlu521/thunder/upload"

// Initialize Qiniu upload manager
qiniuManager, err := upload.InitQiniuUpload("region-id", "bucket", "access-key", "secret-key")
if err != nil {
    // Handle error
}

// Upload file
err = qiniuManager.Upload(context.Background(), "bucket-name", fileReader, "path/to/file")
if err != nil {
    // Handle error
}

// Get public URL
url := qiniuManager.GetPublicURL("your-domain.com", "path/to/file")
```

### Alibaba Cloud OSS

```go
import "github.com/mszlu521/thunder/upload"

// Initialize Alibaba Cloud OSS upload manager
ossManager, err := upload.InitAliyunOSSUpload("access-key-id", "access-key-secret", "endpoint", "bucket-name")
if err != nil {
    // Handle error
}

// Check if service is available
if ossManager.IsAvailable() {
    // Upload file
    err := ossManager.Upload(context.Background(), fileReader, "path/to/file")
    if err != nil {
        // Handle error
    }
}

// Get object URL
url := ossManager.GetObjectURL("endpoint", "bucket-name", "path/to/file")

// Generate signed URL (valid for 3600 seconds)
signedURL, err := ossManager.GetSignedURL("path/to/file", 3600)
```

## Authentication

Thunder provides JWT-based authentication middleware:

```go
import "github.com/mszlu521/thunder/tools/jwt"

// Generate token
token, err := jwt.GenToken("user-id", "username", 24*time.Hour)
if err != nil {
    // Handle error
}

// Parse token
claims, err := jwt.ParseToken(tokenString)
if err != nil {
    // Handle error
}
```

Configure authentication in your config.yml:

```yaml
auth:
  isAuth: true
  ignores:
    - "/api/v1/auth/**"     # Public authentication endpoints
    - "/health"             # Health check endpoint
  needLogins:
    - "/api/v1/user/**"     # User-specific endpoints
```

## Payment Integration

Thunder integrates with WeChat Pay for various payment scenarios:

```go
import "github.com/mszlu521/thunder/pay/wxPay"

// Native payment
payBody := &wxPay.PayBody{
    Description: "Product Description",
    OutTradeNo:  "order-number",
    TimeExpire:  "2025-12-31T10:00:00+08:00",
    Amount:      100, // Amount in cents
    ClientIp:    "127.0.0.1",
}

codeUrl, err := wxPay.Native(context.Background(), payBody)
if err != nil {
    // Handle error
}

// JSAPI payment
jsapiParams, err := wxPay.JsApi(context.Background(), payBody)
if err != nil {
    // Handle error
}

// H5 payment
h5Url, err := wxPay.H5Pay(context.Background(), payBody)
if err != nil {
    // Handle error
}
```

## Event System

Thunder includes a simple event system for decoupling components:

```go
import "github.com/mszlu521/thunder/event"

// Register event handler
event.Register("user.registered", func(e event.Event) (any, error) {
    // Handle user registration event
    userData := e.Data.(map[string]interface{})
    // Process user data
    return "success", nil
})

// Trigger event
result, err := event.Trigger("user.registered", map[string]interface{}{
    "userId": 123,
    "email": "user@example.com",
})
```

## Database Access

Thunder uses GORM for database operations with PostgreSQL and MySQL support:

```go
import (
    "github.com/mszlu521/thunder/database"
    "github.com/mszlu521/thunder/config"
)

// Initialize PostgreSQL
pgConfig := config.GetConfig().DB.Postgres
database.InitPostgres(pgConfig)

// Get database instance
db := database.GetPostgresDB().GormDB

// Perform database operations
var users []User
db.Find(&users)
```

## Contributing

We welcome contributions to Thunder! Please follow these steps:

1. Fork the repository
2. Create a new branch (`git checkout -b feature/your-feature`)
3. Make your changes
4. Commit your changes (`git commit -am 'Add some feature'`)
5. Push to the branch (`git push origin feature/your-feature`)
6. Create a new Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgements

- [Gin](https://github.com/gin-gonic/gin) - HTTP web framework
- [GORM](https://gorm.io/) - ORM library for Go
- [Viper](https://github.com/spf13/viper) - Configuration solution
- [Go-Redis](https://github.com/go-redis/redis) - Redis client for Go
- [Go-Pay](https://github.com/go-pay/gopay) - Payment library