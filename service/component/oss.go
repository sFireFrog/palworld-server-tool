package component

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

func UploadFileToOss(src string, bucketName string, prefix string, client *oss.Client) string {

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		fmt.Println("Bucket Error:", err)
		return ""
	}

	objectKey := filepath.ToSlash(filepath.Join(prefix, filepath.Base(src)))
	fmt.Println("objectKey:", objectKey)
	err = bucket.PutObjectFromFile(objectKey, src)
	if err != nil {
		fmt.Println("PutObjectFromFile Error:", err)
		return ""
	}
	return objectKey
}

func GetFileFromOss(src, dest, bucketName, prefix string, client *oss.Client) bool {

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		fmt.Println("Bucket Error:", err)
		return false
	}
	slashPath := src
	if prefix != "" {
		slashPath = filepath.Join(prefix, src)
	}
	objectKey := filepath.ToSlash(slashPath)
	fmt.Println("objectKey:", objectKey)
	if objectKey == "" {
		return false
	}
	err = bucket.GetObjectToFile(objectKey, dest)
	if err != nil {
		fmt.Println("GetObjectToFile Error:", err)
		return false
	}
	return true
}
