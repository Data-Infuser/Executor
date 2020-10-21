package models

import (
	"time"
)

type Meta struct {
	ID               int          `gorm:"PRIMARY_KEY"`
	Status           string       `gorm:"Column:status"`
	Samples          string       `gorm:"Column:samples"`
	Title            string       `gorm:"column:title"`
	DateType         string       `gorm:"column:dataType"`
	OriginalFileName string       `gorm:"column:originalFileName"`
	RemoteFilePath   string       `gorm:"column:remoteFilePath"`
	FilePath         string       `gorm:"column:filePath"`
	Encoding         string       `gorm:"column:encoding"`
	Extension        string       `gorm:"column:extension"`
	Host             string       `gorm:"column:host"`
	Port             string       `gorm:"column:port"`
	Db               string       `gorm:"column:db"`
	DbUser           string       `gorm:"column:dbUser"`
	Pwd              string       `gorm:"column:pwd"`
	Table            string       `gorm:"column:table"`
	Dbms             string       `gorm:"column:dbms"`
	RowCounts        int          `gorm:"column:rowCounts"`
	Skip             int          `gorm:"column:skip"`
	Sheet            int          `gorm:"column:sheet"`
	UserId           int          `gorm:"column:userId"`
	IsActive         bool         `gorm:"column:isActive"`
	Service          Service      `gorm:"foreignkey:metaId"`
	Stage            Stage        `gorm:"foreignkey:stageId"`
	StageId          int          `gorm:"column:stageId"`
	MetaColumns      []MetaColumn `gorm:"foreignkey:metaId;association_foreignkey:id"`
}

type Stage struct {
	ID     int    `gorm:"PRIMARY_KEY"`
	Status string `gorm:"column:Status"`
	Metas  []Meta `gorm:"foreignkey:stageId;association_foreignkey:id"`
	Name   string `gorm:"column:name"`
}

type MetaColumn struct {
	ID                 int         `gorm:"PRIMARY_KEY"`
	OriginalColumnName string      `gorm:"column:originalColumnName"`
	ColumnName         string      `gorm:"column:columnName"`
	Type               string      `gorm:"column:type"`
	OriginalType       string      `form:"column:originalType"`
	Size               int         `gorm:"column:size"`
	Order              int         `gorm:"column:order"`
	IsHidden           bool        `gorm:"column:isHidden"`
	IsSearchable       bool        `gorm:"column:isSearchable"`
	IsNullable         bool        `gorm:"column:isNullable"`
	DateFormat         string      `gorm:"column:dateFormat"`
	MetaID             int         `gorm:"column:metaId"`
	Params             []MetaParam `gorm:"foreignkey:metaColumnId;association_foreignkey:id"`
	CreatedAt          time.Time   `gorm:"column:createdAt"`
	UpdatedAt          time.Time   `gorm:"column:updatedAt"`
}

type Service struct {
	ID         int       `gorm:"PRIMARY_KEY"`
	Title      string    `gorm:"column:title"`
	EntityName string    `gorm:"column:entityName"`
	UserID     int       `gorm:"column:userId"`
	Meta       *Meta     `gorm:"foreignkey:metaId"`
	MetaID     int       `gorm:"column:metaId"`
	Status     string    `gorm:"column:status"`
	CreatedAt  time.Time `gorm:"column:createdAt"`
	UpdatedAt  time.Time `gorm:"column:updatedAt"`
}

type MetaParam struct {
	ID           int         `gorm:"PRIMARY_KEY"`
	Operator     string      `gorm:"column:operator"`
	Description  string      `gorm:"column:description"`
	IsRequired   bool        `gorm:"column:isRequired"`
	MetaColumn   *MetaColumn `gorm:"foreignkey:metaColumnId"`
	MetaColumnId int         `gorm:"column:metaColumnId"`
	CreatedAt    time.Time   `gorm:"column:createdAt"`
	UpdatedAt    time.Time   `gorm:"column:updatedAt"`
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

func (service Service) TableName() string {
	return "service"
}

func (stage Stage) TableName() string {
	return "stage"
}

func (metaParam MetaParam) TableName() string {
	return "meta_param"
}
