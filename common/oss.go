package common

import (
	"bytes"
	"context"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"path/filepath"
	"strings"
	"time"
)

type Minio struct {
	Client *minio.Client
	Bucket string
	Region string
}

var MinioInstance *Minio

func MinioManagerInstance() *Minio {
	if MinioInstance == nil {
		MinioInstance = MinioNew()
	}
	return MinioInstance
}

func MinioNew() *Minio {
	endpoint, err := beego.AppConfig.String("oss")
	if err != nil {
		return nil
	}
	bucket, err := beego.AppConfig.String("bucket")
	if err != nil {
		return nil
	}
	region, err := beego.AppConfig.String("region")
	if err != nil {
		return nil
	}
	accessKey, err := beego.AppConfig.String("oss_access_key")
	if err != nil {
		return nil
	}
	secretKey, err := beego.AppConfig.String("oss_secret_key")
	if err != nil {
		return nil
	}
	secure, err := beego.AppConfig.Bool("oss_secure")

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: secure,
	})
	if err != nil {
		return nil
	}
	return &Minio{
		Client: minioClient,
		Bucket: bucket,
		Region: region,
	}
}

// GenerateUniqueName 生产唯一文件名
func GenerateUniqueName(fileName string) string {
	id, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}
	return strings.ReplaceAll(id.String(), "-", "") + "_" + time.Now().Format("2006_01_02_15_04_05") + "." + filepath.Ext(fileName)
}

// Upload 上传数据
func (m *Minio) Upload(fileName string, b []byte) (string, error) {
	r := bytes.NewReader(b)
	resultUpload, err := m.Client.PutObject(context.Background(), m.Bucket, GenerateUniqueName(fileName), r, r.Size(), minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}
	return resultUpload.Location, nil
}

func (m *Minio) UploadFile(fileName string, filePath string) (string, error) {
	resultUpload, err := m.Client.FPutObject(context.Background(), m.Bucket, GenerateUniqueName(fileName), filePath, minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}
	return resultUpload.Location, nil
}

func (m *Minio) PresignedGet(objectName string, duration string) (string, error) {
	expires, err := time.ParseDuration(duration)
	if err != nil {
		return "", err
	}
	result, err := m.Client.PresignedGetObject(context.Background(), m.Bucket, objectName, expires, nil)
	if err != nil {
		return "", err
	}
	return result.String(), nil
}
