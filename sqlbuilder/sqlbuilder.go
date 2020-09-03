package sqlbuilder

import (
	"fmt"
	"regexp"
	"strings"

	grpc_executor "queryprocessor/infuser-protobuf/gen/proto/executor"
	"queryprocessor/models"
	"queryprocessor/utils"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// Builder : builder 내의 함수를 사용하기 위한 구조체
type Builder struct{}

var defaultPage int32 = 1
var defaultPerPage int32 = 500

var colToOp *utils.ColTypeToOperation = utils.NewColTypeToOperation()

// GetMeta : 호출 받은 api의 메타 데이터를 가져오는 함수
func (builder *Builder) GetMeta(db *gorm.DB, application string, serviceNm string) *models.Service {
	var service models.Service

	db.Preload("ServiceColumns").Preload("Application").Where(
		"entityName = ?", serviceNm,
	).Joins(
		"JOIN application on application.id = service.applicationId",
	).Where("application.nameSpace = ?", application).First(&service)

	if service.Application.Status != "deployed" {
		err := new(utils.APIError)
		err.Status = 400
		err.Message = "This Application is Not Deployed"

		panic(err)
	} else if service.Status != "loaded" {
		err := new(utils.APIError)
		err.Status = 400
		err.Message = "This Service's Data is Not Loaded"

		panic(err)
	}

	return &service
}

// BuildSQL : API객체와 쿼리 파라미터를 받아 Data DB에서 실제 데이터를 가져올 SQL을 build하는 함수
func (builder *Builder) BuildSQL(service *models.Service, params *grpc_executor.ApiRequest) (string, string, string, map[string]string) {
	tableName := service.Tn
	cols := make([]string, len(service.ServiceColumns))
	colType := make(map[string]string)

	for i, col := range service.ServiceColumns {
		cols[i] = col.ColumnName
		colType[col.ColumnName] = col.Typ
	}

	page, perPage := GetPage(params)

	condition := buildCondition(params, service.ServiceColumns)

	searchQuery := fmt.Sprintf("SELECT %s FROM `%s` %s limit %d, %d", strings.Join(cols, ", "), tableName, condition, (page-1)*perPage, page*perPage)
	cntQuery := fmt.Sprintf("SELECT count(*) as cnt FROM `%s`", tableName)
	matchQuery := fmt.Sprintf("SELECT count(*) as cnt FROM `%s` %s", tableName, condition)

	return searchQuery, matchQuery, cntQuery, colType
}

// GetOperatorByType : 칼럼 타입 별 사용 가능한 Operator 리스트 반환
func GetOperatorByType() utils.ColTypeToOperation {
	return *colToOp
}

// GetPage : get page, perPage parameter to query param
func GetPage(params *grpc_executor.ApiRequest) (int32, int32) {
	var page, perPage *int32

	page = params.Page
	perPage = params.PerPage

	if page == nil {
		page = &defaultPage
	}

	if perPage == nil {
		perPage = &defaultPerPage
	}

	return *page, *perPage
}

func buildCondition(params *grpc_executor.ApiRequest, cols []models.ServiceColumn) string {
	condition := make([]string, 0)
	conditions := params.Cond

	for k, v := range conditions {
		splited := strings.Split(k, "::")
		// println(splited[0])
		// println(splited[1])
		col := arrayInAPIColumn(splited[0], cols)

		if !checkPossibleOperation(col.Typ, splited[1]) {
			err := new(utils.APIError)
			err.Status = 400
			err.Message = "Invalid Operator Error"

			panic(err)
		} else if col == nil {
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

func translateOperation(op string, col *models.ServiceColumn, val string) string {
	switch op {
	case "LT":
		val = wrapValueForType(val, col.Typ)
		return col.ColumnName + " < " + val
	case "LTE":
		val = wrapValueForType(val, col.Typ)
		return col.ColumnName + " <= " + val
	case "GT":
		val = wrapValueForType(val, col.Typ)
		return col.ColumnName + " > " + val
	case "GTE":
		val = wrapValueForType(val, col.Typ)
		return col.ColumnName + " >= " + val
	case "LIKE":
		val = wrapValueForType("%"+val+"%", col.Typ)
		return col.ColumnName + " like " + val
	case "EQ":
		val = wrapValueForType(val, col.Typ)
		return col.ColumnName + " = " + val
	case "NEQ":
		val = wrapValueForType(val, col.Typ)
		return col.ColumnName + " <> " + val
	default:
		return "1=1"
	}
}

func wrapValueForType(val string, colType string) string {
	r, _ := regexp.Compile("[^(]+")

	switch r.FindString(colType) {
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

	r, _ := regexp.Compile("[^(]+")

	switch r.FindString(colType) {
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

func arrayInAPIColumn(key string, arr []models.ServiceColumn) *models.ServiceColumn {
	for _, v := range arr {
		if v.ColumnName == key {
			return &v
		}
	}

	return nil
}
