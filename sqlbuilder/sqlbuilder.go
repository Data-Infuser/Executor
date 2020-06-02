package sqlbuilder

import (
	"fmt"
	"strconv"
	"strings"

	"queryprocessor/models"
	"queryprocessor/utils"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// Builder : builder 내의 함수를 사용하기 위한 구조체
type Builder struct{}

const defaultPage int64 = 1
const defaultPerPage int64 = 500

var colToOp *utils.ColTypeToOperation = utils.NewColTypeToOperation()

// GetMeta : 호출 받은 api의 메타 데이터를 가져오는 함수
func (builder *Builder) GetMeta(db *gorm.DB, api string) *models.API {
	var a models.API
	db.Preload("APIColumns").Where("api.tableName = ?", api).First(&a)

	return &a
}

// BuildSQL : API객체와 쿼리 파라미터를 받아 Data DB에서 실제 데이터를 가져올 SQL을 build하는 함수
func (builder *Builder) BuildSQL(api *models.API, params *gin.Context) (string, string, string, map[string]string) {
	tableName := api.Tn
	cols := make([]string, len(api.APIColumns))
	colType := make(map[string]string)

	for i, col := range api.APIColumns {
		cols[i] = col.ColumnName
		colType[col.ColumnName] = col.Typ
	}

	page, perPage := GetPage(params)

	condition := buildCondition(params, api.APIColumns)

	searchQuery := fmt.Sprintf("SELECT %s FROM %s %s limit %d, %d", strings.Join(cols, ", "), tableName, condition, (page-1)*perPage, page*perPage)
	cntQuery := fmt.Sprintf("SELECT count(*) as cnt FROM %s", tableName)
	matchQuery := fmt.Sprintf("SELECT count(*) as cnt FROM %s %s", tableName, condition)

	return searchQuery, matchQuery, cntQuery, colType
}

// GetPage : get page, perPage parameter to query param
func GetPage(params *gin.Context) (int64, int64) {
	var page, perPage int64
	var err error

	pageStr := params.Query("page")
	perPageStr := params.Query("perPage")

	if pageStr == "" {
		page = defaultPage
	} else {
		page, err = strconv.ParseInt(pageStr, 0, 64)
		if err != nil {
			panic(err.Error())
		}
	}

	if perPageStr == "" {
		perPage = defaultPerPage
	} else {
		perPage, err = strconv.ParseInt(perPageStr, 0, 64)
		if err != nil {
			panic(err.Error())
		}
	}

	return page, perPage
}

func buildCondition(params *gin.Context, cols []models.ApiColumn) string {
	condition := make([]string, 0)
	conditions := params.QueryMap("cond")

	for k, v := range conditions {
		splited := strings.Split(k, "::")
		col := arrayInAPIColumn(splited[0], cols)
		if col == nil {
			err := new(utils.APIError)
			err.Status = 400
			err.Message = "Invalid Parameter Error"

			panic(err)
		}

		condition = append(condition, translateOperation(splited[1], col, v))
	}

	result := strings.Join(condition, " AND ")

	if len(condition) != 0 {
		result = "WHERE " + result
	}

	return result
}

func translateOperation(op string, col *models.ApiColumn, val string) string {
	switch op {
	case "lt":
		val = wrapValueForType(val, col.Typ)
		return col.ColumnName + " < " + val
	case "lte":
		val = wrapValueForType(val, col.Typ)
		return col.ColumnName + " <= " + val
	case "gt":
		val = wrapValueForType(val, col.Typ)
		return col.ColumnName + " > " + val
	case "gte":
		val = wrapValueForType(val, col.Typ)
		return col.ColumnName + " >= " + val
	case "like":
		val = wrapValueForType("%"+val+"%", col.Typ)
		return col.ColumnName + " like " + val
	case "eq":
		val = wrapValueForType(val, col.Typ)
		return col.ColumnName + " = " + val
	case "neq":
		val = wrapValueForType(val, col.Typ)
		return col.ColumnName + " <> " + val
	default:
		return "1=1"
	}
}

func wrapValueForType(val string, colType string) string {
	switch colType {
	case "text":
		fallthrough
	case "longtext":
		fallthrough
	case "varchar":
		val = fmt.Sprintf("'%s'", val)

	case "date":
		val = fmt.Sprintf("STR_TO_DATE('%s', '%%Y-%%m-%%d')", val)
	case "datetime":
		val = fmt.Sprintf("STR_TO_DATE('%s', '%%Y-%%m-%%d %%h:%%i:%%s')", val)
	}

	return val
}

func checkPossibleOperation(colType string, operation string) bool {
	var arr []string

	switch colType {
	case "text":
		fallthrough
	case "longtext":
		fallthrough
	case "varchar":
		arr = colToOp.Str

	case "bigint":
		fallthrough
	case "int":
		fallthrough
	case "float":
		fallthrough
	case "double":
		arr = colToOp.Number

	case "boolean":
		arr = colToOp.Boolean

	case "date":
		fallthrough
	case "datetime":
		arr = colToOp.Dt
	}

	return arrayInStr(operation, arr) != nil
}

func arrayInStr(key string, arr []string) *string {
	for _, v := range arr {
		if v == key {
			return &v
		}
	}

	return nil
}

func arrayInAPIColumn(key string, arr []models.ApiColumn) *models.ApiColumn {
	for _, v := range arr {
		if v.ColumnName == key {
			return &v
		}
	}

	return nil
}
