package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// 解析yml文件
type ConfigInfo struct {
	Endpoint        string `yaml:"endpoint"`
	AccessKeyID     string `yaml:"accessKeyID"`
	AccessKeySecret string `yaml:"accessKeySecret"`
	BucketName      string `yaml:"bucketName"`
	Prefix          string `yaml:"prefix"`
}

func (c *ConfigInfo) GetConf() (*ConfigInfo, bool) {
	yamlFile, err := os.ReadFile("config.yml")
	if err != nil {
		fmt.Println("ReadFile error", err.Error())
		return c, false
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Println("Unmarshal error:", err.Error())
		return c, false
	}
	return c, true
}
