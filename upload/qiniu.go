package upload

import (
	"context"
	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/http_client"
	"github.com/qiniu/go-sdk/v7/storagev2/region"
	"github.com/qiniu/go-sdk/v7/storagev2/uploader"
	"io"
	"thunder/config"
)

var QiniuUploadManager *QiniuUpload

type QiniuUpload struct {
	uploadManager *uploader.UploadManager
}

func InitQiniuUpload() {
	regionId := config.Conf.Qiniu.Region
	accessKey := config.Conf.Qiniu.AccessKey
	secretKey := config.Conf.Qiniu.SecretKey
	mac := credentials.NewCredentials(accessKey, secretKey)
	options := uploader.UploadManagerOptions{
		Options: http_client.Options{
			Credentials: mac,
			Regions:     region.GetRegionByID(regionId, true),
		},
	}
	uploadManager := uploader.NewUploadManager(&options)
	QiniuUploadManager = &QiniuUpload{
		uploadManager: uploadManager,
	}
}

func (q *QiniuUpload) Upload(ctx context.Context, reader io.Reader, name string) error {
	bucket := config.Conf.Qiniu.Bucket
	err := q.uploadManager.UploadReader(ctx, reader, &uploader.ObjectOptions{
		BucketName: bucket,
		ObjectName: &name,
		FileName:   name,
	}, nil)
	return err
}
