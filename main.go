package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/tmkontra/s3-explorer/filesystem"
)

type PathRequest struct {
	Path string `uri:"Path" binding:"required"`
}

type ServerConfig struct {
	Env        string `mapstructure:"ENV"`        // "development", "production"
	Filesystem string `mapstructure:"FILESYSTEM"` // "local", "s3"
	Port       int    `mapstructure:"PORT"`

	// s3 configs
	BucketName   string `mapstructure:"BUCKET_NAME"`
	BucketRegion string `mapstructure:"BUCKET_REGION"`
}

func (c *ServerConfig) GetFilesystem() (filesystem.Filesystem, error) {
	if c.Filesystem == "local" {
		return filesystem.NewLocalFilesystem(), nil
	} else if c.Filesystem == "s3" {
		if c.BucketName == "" || c.BucketRegion == "" {
			return nil, fmt.Errorf(
				"BUCKET_NAME and BUCKET_REGION must be specified, got '%s' and '%s', respectively",
				c.BucketName, c.BucketRegion)
		}
		return filesystem.NewS3Filesystem(c.BucketName, c.BucketRegion), nil
	} else {
		return nil, fmt.Errorf("FILESYSTEM must be one of 'local' or 's3', got '%s'", c.Filesystem)
	}
}

func main() {
	config, err := loadConfig()
	if err != nil {
		panic(err)
	}
	if config.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	fs, err := config.GetFilesystem()
	if err != nil {
		panic(err)
	}
	r.GET("/*Path", func(c *gin.Context) {
		var request PathRequest
		if err := c.ShouldBindUri(&request); err != nil {
			c.JSON(400, gin.H{"error": err})
			return
		}
		result, err := filesystem.GetPath(fs, request.Path)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, result.GetData())
	})
	r.Run(fmt.Sprintf(":%v", config.Port)) // listen and serve on 0.0.0.0:PORT
}

func loadConfig() (config ServerConfig, err error) {
	viper.SetDefault("Port", 8080)
	viper.SetDefault("Filesystem", "local")
	viper.SetDefault("Env", "development")
	viper.BindEnv("ENV")
	viper.BindEnv("BUCKET_NAME")
	viper.BindEnv("BUCKET_REGION")
	viper.AutomaticEnv()

	err = viper.Unmarshal(&config)
	return
}
