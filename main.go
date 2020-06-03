package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"runtime"

	"queryprocessor/sqlbuilder"
	"queryprocessor/sqlexecutor"
	"queryprocessor/utils"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"gopkg.in/yaml.v2"
)

// DBConfig : Database Config
type DBConfig struct {
	DBName   string `yaml:"dbName"`
	DBType   string `yaml:"dbType"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// ServerConfig : Server side Config
type ServerConfig struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	MaxUseCPU int    `yaml:"maxUseCPU"`
	Env       string `yaml:"env"`
}

// Config : Whole Config Information
type Config struct {
	MetaDB DBConfig     `yaml:"metaDB"`
	DataDB DBConfig     `yaml:"dataDB"`
	Server ServerConfig `yaml:"server"`
}

var metaDB *gorm.DB
var dataDB *gorm.DB

func dbConnect(config DBConfig) *gorm.DB {
	connectURL := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True",
		config.User, config.Password, config.Host, config.Port, config.DBName)

	db, err := gorm.Open(config.DBType, connectURL)
	if err != nil {
		panic(err.Error())
	}

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	return db
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	e := new(sqlexecutor.Executor)
	b := new(sqlbuilder.Builder)

	r.GET("/rest/:application/:api", func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				switch err.(type) {
				case *utils.APIError:
					apiError := err.(*utils.APIError)
					c.JSON(http.StatusBadRequest, apiError)
				default:
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
				}
			}
		}()

		api := b.GetMeta(metaDB, c.Param("application"), c.Param("api"))
		searchSQL, matchSQL, countSQL, colType := b.BuildSQL(api, c)
		data, matchCnt, totalCnt := e.Execute(dataDB, searchSQL, matchSQL, countSQL, colType)

		page, perPage := sqlbuilder.GetPage(c)

		c.JSON(http.StatusOK, gin.H{
			"page":         page,
			"perPage":      perPage,
			"currentCount": len(data),
			"matchCount":   matchCnt,
			"totalCount":   totalCnt,
			"data":         data,
		})
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page Not Found"})
	})

	return r
}

func main() {
	// GC 호출을 줄이기 위한 방법
	ballast := make([]byte, 10<<30)
	_ = ballast

	filename, _ := filepath.Abs("config.yaml")
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err.Error())
	}

	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		panic(err.Error())
	}

	var env string

	switch config.Server.Env {
	case "development":
		env = "debug"
	case "production":
		env = "release"
	}

	gin.SetMode(env)
	runtime.GOMAXPROCS(config.Server.MaxUseCPU)

	metaDB = dbConnect(config.MetaDB)
	dataDB = dbConnect(config.DataDB)
	defer metaDB.Close()
	defer dataDB.Close()

	r := setupRouter()
	r.Run(fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port))
}
