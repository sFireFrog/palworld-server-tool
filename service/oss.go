package service

import (
	"fmt"
	"path/filepath"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func NewOssClient(endpoint string, accessKeyID string, accessKeySecret string) (*oss.Client, error) {
	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		fmt.Println("New Client Error:", err)
		return nil, err
	}
	return client, err
}

func UploadFileToOss(src string, bucketName string, prefix string, client *oss.Client) {

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		fmt.Println("Bucket Error:", err)
		return
	}

	objectKey := filepath.ToSlash(filepath.Join(prefix, filepath.Base(src)))
	fmt.Println("objectKey:", objectKey)
	err = bucket.PutObjectFromFile(objectKey, src)
	if err != nil {
		fmt.Println("PutObjectFromFile Error:", err)
		return
	}
}
