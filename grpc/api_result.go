package server

import (
	"context"
	"fmt"
	handler "queryprocessor/handler"
	grpc_executor "queryprocessor/infuser-protobuf/gen/proto/executor"
	"queryprocessor/sqlbuilder"
)

type apiResultServer struct {
	handler *handler.ApiResultHandler
}

func newApiResultServer(handler *handler.ApiResultHandler) grpc_executor.ApiResultServiceServer {
	return &apiResultServer{
		handler: handler,
	}
}

// GetApiResult : 실질적인 Api 호출 처리 기능
func (s *apiResultServer) GetApiResult(ctx context.Context, req *grpc_executor.ApiRequest) (*grpc_executor.ApiResult, error) {
	//e := new(sqlexecutor.Executor)
	b := new(sqlbuilder.Builder)

	service := b.GetMeta(s.handler.Ctx.MetaDB, req.StageId, req.ServiceId)
	fmt.Printf("%+v", service)
	fmt.Printf("213213213231")
	// searchSQL, matchSQL, countSQL, colType := b.BuildSQL(service, req)
	// data, matchCnt, totalCnt := e.Execute(s.handler.Ctx.DataDB, searchSQL, matchSQL, countSQL, colType)

	// page, perPage := sqlbuilder.GetPage(req)

	// apiResult := new(grpc_executor.ApiResult)
	// apiResult.Data = data
	// apiResult.Page = page
	// apiResult.PerPage = perPage
	// apiResult.CurrentCount = int32(len(data))
	// apiResult.MatchCount = int32(matchCnt)
	// apiResult.TotalCount = int32(totalCnt)

	return nil, nil
}
