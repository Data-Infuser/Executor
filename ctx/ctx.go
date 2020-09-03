package ctx

import "github.com/jinzhu/gorm"

// Context : Shared Object
type Context struct {
	MetaDB *gorm.DB
	DataDB *gorm.DB
}

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
