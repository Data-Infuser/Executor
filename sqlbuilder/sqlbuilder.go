package sqlbuilder

import (
	"fmt"
	"queryprocessor/models"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// Builder : builder 내의 함수를 사용하기 위한 구조체
type Builder struct{}

// GetMeta : 호출 받은 api의 메타 데이터를 가져오는 함수
func (builder *Builder) GetMeta(db *gorm.DB, api string) *models.API {
	var a models.API
	db.Preload("APIColumns").Where("api.tableName = ?", api).First(&a)

	return &a
}

// BuildSQL : API객체와 쿼리 파라미터를 받아 Data DB에서 실제 데이터를 가져올 SQL을 build하는 함수
func (builder *Builder) BuildSQL(api *models.API, params map[string]string) (string, string, map[string]string) {
	tableName := api.Tn
	cols := make([]string, len(api.APIColumns))
	colType := make(map[string]string)

	for i, col := range api.APIColumns {
		cols[i] = col.ColumnName
		colType[col.ColumnName] = col.Typ
	}

	searchQuery := fmt.Sprintf("SELECT %s FROM %s limit %d, %d", strings.Join(cols, ", "), tableName, 0, 500)
	cntQuery := fmt.Sprintf("SELECT count(*) as cnt FROM %s", tableName)

	return searchQuery, cntQuery, colType
}
