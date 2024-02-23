package tool

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/viper"
	"github.com/zaigie/palworld-server-tool/internal/auth"
	"github.com/zaigie/palworld-server-tool/internal/database"
	"github.com/zaigie/palworld-server-tool/internal/logger"
	"github.com/zaigie/palworld-server-tool/internal/source"
	"github.com/zaigie/palworld-server-tool/service/component"
)

type Sturcture struct {
	Players []database.Player `json:"players"`
	Guilds  []database.Guild  `json:"guilds"`
}

type fileResponse struct {
	Path string `json:"path"`
}

func getSavCli() (string, error) {
	savCliPath := viper.GetString("save.decode_path")
	if _, err := os.Stat(savCliPath); err != nil {
		return "", err
	}
	return savCliPath, nil
}

func ConversionLoading(file string) error {
	var tmpFile string
	var err error

	savCli, err := getSavCli()
	if err != nil {
		return errors.New("error getting executable path: " + err.Error())
	}

	if strings.HasPrefix(file, "http://") || strings.HasPrefix(file, "https://") {
		// http(s)://url
		logger.Infof("downloading Level.sav from %s\n", file)
		useOss := viper.GetBool("save.use_oss")
		if useOss {
			fileRes := fileResponse{}
			err = source.GetRequest(file, fileRes)
			if err != nil {
				tmpFile, err = useOssDownload(fileRes.Path)
			}
		} else {
			tmpFile, err = source.DownloadFromHttp(file)
		}
		if err != nil {
			return errors.New("error downloading file: " + err.Error())
		}
		logger.Info("Level.sav downloaded\n")
	} else if strings.HasPrefix(file, "k8s://") {
		// k8s://namespace/pod/container:remotePath
		logger.Infof("copy Level.sav from %s\n", file)
		namespace, podName, container, remotePath, err := source.ParseK8sAddress(file)
		if err != nil {
			return errors.New("error parsing k8s address: " + err.Error())
		}
		tmpFile, err = source.CopyFromPod(namespace, podName, container, remotePath)
		if err != nil {
			return errors.New("error copying file from pod: " + err.Error())
		}
	} else if strings.HasPrefix(file, "docker://") {
		// docker://containerID(Name):remotePath
		logger.Infof("copy Level.sav from %s\n", file)
		containerId, remotePath, err := source.ParseDockerAddress(file)
		if err != nil {
			return errors.New("error parsing docker address: " + err.Error())
		}
		tmpFile, err = source.CopyFromContainer(containerId, remotePath)
		if err != nil {
			return errors.New("error copying file from container: " + err.Error())
		}
	} else {
		// local file
		tmpFile, err = source.CopyFromLocal(file)
		if err != nil {
			return errors.New("error copying file to temporary directory: " + err.Error())
		}
	}
	defer os.Remove(tmpFile)

	baseUrl := "http://127.0.0.1"
	if viper.GetBool("web.tls") {
		baseUrl = "https://127.0.0.1"
	}

	requestUrl := fmt.Sprintf("%s:%d/api/", baseUrl, viper.GetInt("web.port"))
	tokenString, err := auth.GenerateToken()
	if err != nil {
		return errors.New("error generating token: " + err.Error())
	}
	execArgs := []string{"-f", tmpFile, "--request", requestUrl, "--token", tokenString, "--clear"}
	cmd := exec.Command(savCli, execArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		return errors.New("error starting command: " + err.Error())
	}
	err = cmd.Wait()
	if err != nil {
		return errors.New("error waiting for command: " + err.Error())
	}

	return nil
}

func useOssDownload(src string) (string, error) {
	endpoint := viper.GetString("oss.endpoint")
	accessKeyID := viper.GetString("oss.accessKeyID")
	accessKeySecret := viper.GetString("oss.accessKeySecret")
	bucketName := viper.GetString("oss.bucketName")
	client, err := component.NewOssClient(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		return "", errors.New("init oss client error: " + err.Error())
	}
	tmpFile, err := os.CreateTemp("", "Level.sav")
	destPath := tmpFile.Name()
	getStatus := component.GetFileFromOss(src, destPath, bucketName, "", client)
	if !getStatus {
		return "", errors.New("get file error: " + err.Error())
	}
	return destPath, nil
}
