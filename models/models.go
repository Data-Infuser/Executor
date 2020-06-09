package models

import (
	"time"
)

type Meta struct {
	ID               int          `gorm:"PRIMARY_KEY"`
	Title            string       `gorm:"column:title"`
	DateType         string       `gorm:"column:dataType"`
	OriginalFileName string       `gorm:"column:originalFileName"`
	FilePath         string       `gorm:"column:filePath"`
	Extension        string       `gorm:"column:extension"`
	Host             string       `gorm:"column:host"`
	Port             string       `gorm:"column:port"`
	Db               string       `gorm:"column:db"`
	DbUser           string       `gorm:"column:dbUser"`
	Pwd              string       `gorm:"column:pwd"`
	Table            string       `gorm:"column:table"`
	RowCounts        int          `gorm:"column:rowCounts"`
	Skip             int          `gorm:"column:skip"`
	Sheet            int          `gorm:"column:sheet"`
	IsActive         bool         `gorm:"column:isActive"`
	Service          Service      `gorm:"foreignkey:metaId"`
	MetaColumns      []MetaColumn `gorm:"foreignkey:metaId;association_foreignkey:id"`
}

type MetaColumn struct {
	ID                 int       `gorm:"PRIMARY_KEY"`
	OriginalColumnName string    `gorm:"column:originalColumnName"`
	ColumnName         string    `gorm:"column:columnName"`
	typ                string    `gorm:"column:type"`
	Size               int       `gorm:"column:size"`
	Order              int       `gorm:"column:order"`
	MetaID             int       `gorm:"column:metaId"`
	CreatedAt          time.Time `gorm:"column:createdAt"`
	UpdatedAt          time.Time `gorm:"column:updatedAt"`
}

type User struct {
	ID        int       `gorm:"PRIMARY_KEY"`
	Username  string    `gorm:"column:userName"`
	Password  string    `gorm:"column:password"`
	CreatedAt time.Time `gorm:"column:createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt"`
	Metas     []Meta    `gorm:"foreignkey:userId"`
}

type Application struct {
	ID             int       `gorm:"PRIMARY_KEY"`
	NameSpace      string    `gorm:"column:nameSpace"`
	Title          string    `gorm:"column:title"`
	Desc           string    `gorm:"column:description"`
	CreatedAt      time.Time `gorm:"column:createdAt"`
	UpdatedAt      time.Time `gorm:"column:updatedAt"`
	User           *User     `gorm:"foreignkey:userId"`
	UserID         int       `gorm:"column:userId"`
	Status         string    `gorm:"column:status"`
	ServiceColumns []Service `gorm:"foreignkey:applicationId;association_foreignkey:id"`
}

type Service struct {
	ID             int             `gorm:"PRIMARY_KEY"`
	Title          string          `gorm:"column:title"`
	EntityName     string          `gorm:"column:entityName"`
	Tn             string          `gorm:"column:tableName"`
	DataCounts     int             `gorm:"column:dataCounts"`
	User           *User           `gorm:"foreignkey:userId"`
	UserID         int             `gorm:"column:userId"`
	Meta           *Meta           `gorm:"foreignkey:metaId"`
	MetaID         int             `gorm:"column:metaId"`
	Status         string          `gorm:"column:status"`
	ServiceColumns []ServiceColumn `gorm:"foreignkey:serviceId;association_foreignkey:id"`
	Application    Application     `gorm:"foreignkey:applicationId"`
	ApplicationID  int             `gorm:"column:applicationId"`
	CreatedAt      time.Time       `gorm:"column:createdAt"`
	UpdatedAt      time.Time       `gorm:"column:updatedAt"`
}

type ServiceColumn struct {
	ID         int    `gorm:"PRIMARY_KEY"`
	ColumnName string `gorm:"column:columnName"`
	Typ        string `gorm:"column:type"`
	Hidden     bool   `gorm:"column:hidden"`
	ServiceID  int    `gorm:"column:serviceId"`
}

type CountRecord struct {
	Cnt int `gorm:"column:cnt"`
}

func (meta Meta) TableName() string {
	return "meta"
}

func (metaColumn MetaColumn) TableName() string {
	return "meta_column"
}

func (application Application) TableName() string {
	return "application"
}

func (service Service) TableName() string {
	return "service"
}

func (sericeColumn ServiceColumn) TableName() string {
	return "service_column"
}

func (user User) TableName() string {
	return "user"
}
