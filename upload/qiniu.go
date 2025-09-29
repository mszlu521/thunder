package upload

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/http_client"
	"github.com/qiniu/go-sdk/v7/storagev2/objects"
	"github.com/qiniu/go-sdk/v7/storagev2/region"
	"github.com/qiniu/go-sdk/v7/storagev2/uploader"
)


type QiniuUpload struct {
	uploadManager *uploader.UploadManager
	bucket      *objects.Bucket
	credentials *credentials.Credentials
}

func InitQiniuUpload(regionId, bucket, accessKey, secretKey string) (*QiniuUpload, error) {
	mac := credentials.NewCredentials(accessKey, secretKey)
	clientOptions := http_client.Options{
		Credentials: mac,
		Regions:     region.GetRegionByID(regionId, true),
	}

	uploadManager := uploader.NewUploadManager(&uploader.UploadManagerOptions{
		Options: clientOptions,
	})
	objectsManager := objects.NewObjectsManager(&objects.ObjectsManagerOptions{
		Options: clientOptions,
	})
	b := objectsManager.Bucket(bucket)
	qiniuUploadManager := &QiniuUpload{
		uploadManager: uploadManager,
		bucket:      b,
		credentials: mac,
	}
	return qiniuUploadManager, nil
}

func (q *QiniuUpload) Upload(ctx context.Context, bucket string, reader io.Reader, name string) error {
	err := q.uploadManager.UploadReader(ctx, reader, &uploader.ObjectOptions{
		BucketName: bucket,
		ObjectName: &name,
		FileName:   name,
	}, nil)
	return err
}

// Delete 从七牛云存储中删除指定的文件
func (q *QiniuUpload) Delete(ctx context.Context, key string) error {
	return q.bucket.Object(key).Delete().Call(ctx)
}

// GetPublicURL 获取文件的公开访问URL
func (q *QiniuUpload) GetPublicURL(domain string, key string) string {
	// 构建公开访问的URL
	if strings.HasPrefix(domain, "http") {
		return fmt.Sprintf("%s/%s", domain, key)
	}
	u := url.URL{
		Scheme: "https",
		Host:   domain,
		Path:   key,
	}
	return u.String()
}

