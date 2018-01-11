package main

import (
	"fmt"
	"io"
	"os"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/shen100/golang123/config"
	"github.com/shen100/golang123/model"
	"github.com/shen100/golang123/router"
	"github.com/shen100/golang123/utils"
)

func init() {
	db, err := gorm.Open(config.DBConfig.Dialect, config.DBConfig.URL)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	if config.ServerConfig.Env == model.DevelopmentMode {
		db.LogMode(true)
	}
	db.DB().SetMaxIdleConns(config.DBConfig.MaxIdleConns)
	db.DB().SetMaxOpenConns(config.DBConfig.MaxOpenConns)
	model.DB = db

	r, err := redis.Dial("tcp", config.RedisConfig.URL)
	if err != nil {
		fmt.Println("Connect to redis error", err.Error())
		os.Exit(-1)
	}
	utils.RedisConn = r
}

func main() {
	fmt.Println("gin.Version: ", gin.Version)
	if config.ServerConfig.Env != model.DevelopmentMode {
		// Disable Console Color, you don't need console color when writing the logs to file.
		gin.DisableConsoleColor()
		// Logging to a file.
		logFile, err := os.OpenFile(config.ServerConfig.LogFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			fmt.Printf(err.Error())
			os.Exit(-1)
		}
		gin.DefaultWriter = io.MultiWriter(logFile)
	}

	// Creates a router without any middleware by default
	app := gin.New()

	// Set a lower memory limit for multipart forms (default is 32 MiB)
	maxSize := int64(config.ServerConfig.MaxMultipartMemory)
	app.MaxMultipartMemory = maxSize << 20 // 3 MiB

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	app.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	app.Use(gin.Recovery())

	router.Route(app)

	app.Run(":" + fmt.Sprintf("%d", config.ServerConfig.Port))
}
