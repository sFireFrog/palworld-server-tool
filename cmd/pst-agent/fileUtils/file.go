package fileUtils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/zaigie/palworld-server-tool/cmd/pst-agent/config"
	"github.com/zaigie/palworld-server-tool/service/component"
)

func CopyFile(src, dst string) bool {
	source, err := os.Open(src)
	if err != nil {
		log.Println(err)
		return false
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		log.Println(err)
		return false
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

// limitCacheFiles keeps only the latest `n` files in the cache directory
func LimitCacheFiles(cacheDir string, n int) {
	files, err := os.ReadDir(cacheDir)
	if err != nil {
		log.Println("Error reading cache directory:", err)
		return
	}

	if len(files) <= n {
		return
	}

	sort.Slice(files, func(i, j int) bool {
		infoI, _ := files[i].Info()
		infoJ, _ := files[j].Info()
		return infoI.ModTime().After(infoJ.ModTime())
	})

	// Delete files that exceed the limit
	for i := n; i < len(files); i++ {
		path := filepath.Join(cacheDir, files[i].Name())
		err = os.RemoveAll(path)
		if err != nil {
			fmt.Println("delete files path", path, err)
		}
	}
}

func UploadFileToOss(src string) string {
	info := config.ConfigInfo{}
	config, getConfStatus := info.GetConf()
	if !getConfStatus {
		fmt.Println("not config")
		return ""
	}
	client, err := component.NewOssClient(config.Endpoint, config.AccessKeyID, config.AccessKeySecret)
	if err != nil {
		fmt.Println("NewOssClient error")
		return ""
	}
	path := component.UploadFileToOss(src, config.BucketName, config.Prefix, client)
	return path
}

func GetFileFromOss(src, dest string) {
	info := config.ConfigInfo{}
	config, getConfStatus := info.GetConf()
	if !getConfStatus {
		fmt.Println("not config")
		return
	}
	client, err := component.NewOssClient(config.Endpoint, config.AccessKeyID, config.AccessKeySecret)
	if err != nil {
		fmt.Println("NewOssClient error")
		return
	}
	component.GetFileFromOss(src, dest, config.BucketName, config.Prefix, client)
}
