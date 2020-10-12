package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"

	ctx "queryprocessor/ctx"
	server "queryprocessor/grpc"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"gopkg.in/yaml.v2"
)

func dbConnect(config ctx.DBConfig) *gorm.DB {
	connectURL := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True",
		config.User, config.Password, config.Host, config.Port, config.DBName)

	db, err := gorm.Open(config.DBType, connectURL)
	if err != nil {
		panic(err.Error())
	}

	db.LogMode(true)

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	return db
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	// e := new(sqlexecutor.Executor)
	// b := new(sqlbuilder.Builder)

	// r.GET("/operators", func(c *gin.Context) {
	// 	c.JSON(http.StatusOK, sqlbuilder.GetOperatorByType())
	// })

	return r
}

func main() {
	// GC 호출을 줄이기 위한 방법
	ballast := make([]byte, 10<<24)
	_ = ballast

	filename, _ := filepath.Abs("config.yaml")
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err.Error())
	}

	var config ctx.Config
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

	ctx := new(ctx.Context)

	ctx.MetaDB = dbConnect(config.MetaDB)
	ctx.DataDB = dbConnect(config.DataDB)
	defer ctx.MetaDB.Close()
	defer ctx.DataDB.Close()

	var network = flag.String("network", "tcp", `one of "tcp" or "unix". Must be consistent to -endpoint`)

	s := server.New(ctx)
	if err := s.Run(*network, fmt.Sprintf(":%d", config.Server.Port)); err != nil {
		println("Service Run failed")
		println(err.Error())
	}
}
