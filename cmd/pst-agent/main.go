package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/zaigie/palworld-server-tool/cmd/pst-agent/fileUtils"
)

var (
	port int
	file string
)

type fileResponse struct {
	Path string `json:"path"`
}

func main() {
	flag.IntVar(&port, "port", 8081, "port")
	flag.StringVar(&file, "f", "", "Level.sav file location")
	flag.Parse()

	viper.BindEnv("sav_file", "SAV_FILE")

	viper.SetDefault("port", port)
	viper.SetDefault("sav_file", file)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/sync", func(c *gin.Context) {
		uuid := uuid.New().String()
		cacheDir := filepath.Join(os.TempDir(), "pst", uuid)
		os.MkdirAll(cacheDir, os.ModePerm)

		destFile := filepath.Join(cacheDir, "Level.sav")
		copyStatus := fileUtils.CopyFile(viper.GetString("sav_file"), destFile)
		if !copyStatus {
			c.Redirect(http.StatusFound, "/404")
			return
		}

		c.File(destFile)
	})

	r.GET("/sync/v2", func(c *gin.Context) {
		uuid := uuid.New().String()
		cacheDir := filepath.Join(os.TempDir(), "pst", uuid)
		os.MkdirAll(cacheDir, os.ModePerm)

		destFile := filepath.Join(cacheDir, uuid)
		copyStatus := fileUtils.CopyFile(viper.GetString("sav_file"), destFile)
		if !copyStatus {
			c.JSON(http.StatusBadRequest, gin.H{"error": "not get"})
			return
		}
		path := fileUtils.UploadFileToOss(destFile)
		if path == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "upload error"})
			return
		}
		c.JSON(http.StatusOK, fileResponse{Path: path})
	})

	s, err := gocron.NewScheduler()
	if err != nil {
		fmt.Println(err)
	}
	_, err = s.NewJob(
		gocron.DurationJob(60*time.Second),
		gocron.NewTask(fileUtils.LimitCacheFiles, filepath.Join(os.TempDir(), "pst"), 5),
	)
	if err != nil {
		fmt.Println(err)
	}
	s.Start()

	fmt.Println("pst-agent is running, Listening on port", port)

	r.Run(":" + strconv.Itoa(port))
}
