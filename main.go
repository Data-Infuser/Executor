package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"queryprocessor/sqlbuilder"
	"queryprocessor/sqlexecutor"

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
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
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

	r.GET("/api/:api", func(c *gin.Context) {
		api := b.GetMeta(metaDB, c.Param("api"))
		searchSQL, countSQL, colType := b.BuildSQL(api, c.QueryMap("param"))
		data, cnt := e.Execute(dataDB, searchSQL, countSQL, colType)

		c.JSON(http.StatusOK, gin.H{
			"currentCount": len(data),
			"totalCount":   cnt,
			"data":         data,
		})
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

	metaDB = dbConnect(config.MetaDB)
	dataDB = dbConnect(config.DataDB)

	r := setupRouter()

	r.Run(fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port))

	defer metaDB.Close()
	defer dataDB.Close()
}
