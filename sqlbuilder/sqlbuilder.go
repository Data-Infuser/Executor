package sqlbuilder

import (
	"fmt"
	"queryprocessor/models"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// Builder : builder 내의 함수를 사용하기 위한 구조체
type Builder struct{}

const DEFAULT_PAGE int64 = 0
const DEFAULT_PER_PAGE int64 = 500

// GetMeta : 호출 받은 api의 메타 데이터를 가져오는 함수
func (builder *Builder) GetMeta(db *gorm.DB, api string) *models.API {
	var a models.API
	db.Preload("APIColumns").Where("api.tableName = ?", api).First(&a)

	return &a
}

// BuildSQL : API객체와 쿼리 파라미터를 받아 Data DB에서 실제 데이터를 가져올 SQL을 build하는 함수
func (builder *Builder) BuildSQL(api *models.API, params *gin.Context) (string, string, map[string]string) {
	tableName := api.Tn
	cols := make([]string, len(api.APIColumns))
	colType := make(map[string]string)

	for i, col := range api.APIColumns {
		cols[i] = col.ColumnName
		colType[col.ColumnName] = col.Typ
	}

	page, perPage := GetPage(params)

	searchQuery := fmt.Sprintf("SELECT %s FROM %s limit %d, %d", strings.Join(cols, ", "), tableName, page*perPage, (page+1)*perPage)
	cntQuery := fmt.Sprintf("SELECT count(*) as cnt FROM %s", tableName)

	return searchQuery, cntQuery, colType
}

// GetPage : get page, perPage parameter to query param
func GetPage(params *gin.Context) (int64, int64) {
	var page, perPage int64
	var err error

	pageStr := params.Query("page")
	perPageStr := params.Query("perPage")

	if pageStr == "" {
		page = DEFAULT_PAGE
	} else {
		page, err = strconv.ParseInt(pageStr, 0, 64)
		if err != nil {
			panic(err.Error())
		}
	}

	if perPageStr == "" {
		perPage = DEFAULT_PER_PAGE
	} else {
		perPage, err = strconv.ParseInt(perPageStr, 0, 64)
		if err != nil {
			panic(err.Error())
		}
	}

	return page, perPage
}
