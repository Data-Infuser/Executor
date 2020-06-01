# Api Query Processor

## 본 프로젝트는 API GEN Sub Project

### 프로젝트 목적
* 외부 사용자의 API 호출 시, 해당 API에 대한 정보를 가져와서 Query를 Build, Execute한 뒤 그 결과를 리턴
* 본 프로젝트는 고부하 환경에서 성능 이슈 발생 가능성이 예상되어 Golang으로 작성됨

### 프로젝트 환경
* Golang 1.13.4
* MySQL 5.7


### Go 프로젝트 Dependency 설치
> go get -u -v -f all

### 프로젝트 실행
> go run main.go

### 프로젝트 빌드 (프로젝트 루트 폴더에 실행파일이 생성됨)
> go build
