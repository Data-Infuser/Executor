package sqlexecutor

import (
	"database/sql"
	"queryprocessor/models"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// Executor : Executor 내의 함수를 사용하기 위한 구조체
type Executor struct{}

// Execute : Data DB에 SQL을 실행하여 데이터를 가져오는 함수
func (e *Executor) Execute(
	db *gorm.DB, query string, matchQuery string, cntQuery string, colType map[string]string) ([]map[string]interface{}, int, int) {
	if db == nil {
		panic("error")
	}

	totalCnt := make(chan int)
	go getCount(db, cntQuery, totalCnt)

	matchCnt := make(chan int)
	go getCount(db, matchQuery, matchCnt)

	rows, err := db.Raw(query).Rows()
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}

	values := make([]sql.RawBytes, len(columns))

	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	var results []map[string]interface{}
	results = make([]map[string]interface{}, 0)

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
		}

		var row map[string]interface{}
		row = make(map[string]interface{})

		var value interface{}
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = typeParse(colType[columns[i]], string(col))
			}
			row[columns[i]] = value
		}

		results = append(results, row)
	}
	if err = rows.Err(); err != nil {
		panic(err.Error())
	}

	return results, <-matchCnt, <-totalCnt
}

func getCount(db *gorm.DB, cntQuery string, cnt chan int) {
	var totalCnt models.CountRecord
	db.Raw(cntQuery).First(&totalCnt)

	cnt <- totalCnt.Cnt
}

func typeParse(colType string, data string) interface{} {
	var value interface{}
	var t time.Time
	var err error

	switch colType {
	case "text":
		fallthrough
	case "longtext":
		fallthrough
	case "varchar":
		value = data

	case "bigint":
		fallthrough
	case "int":
		value, err = strconv.ParseInt(data, 0, 64)

	case "float":
		value, err = strconv.ParseFloat(data, 32)
	case "double":
		value, err = strconv.ParseFloat(data, 64)

	case "boolean":
		value, err = strconv.ParseBool(data)

	case "date":
		t, err = time.Parse("2006-01-02T15:04:05Z", data)
		value = t.Format("2006-01-02")
	case "datetime":
		t, err = time.Parse("2006-01-02T15:04:05Z", data)
		value = t.Format("2006-01-02 15:04:05")

	default:
		value = data
	}

	if err != nil {
		panic(err.Error())
	}

	return value
}
