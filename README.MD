
# Executor
> ``Data Infuser`` / Executor 프로젝트 입니다.

외부 사용자의 API 호출 시, 해당 API에 대한 정보를 가져와서 Query를 Build, Execute한 뒤 그 결과를 리턴
본 프로젝트는 고부하 환경에서 성능 이슈 발생 가능성이 예상되어 Golang으로 작성됨

## Environment
 * Golang 1.13.4
 * MySQL 5.7
 * GIN (go http web framework, https://github.com/gin-gonic/gin)
 * GORM (go orm library, https://github.com/jinzhu/gorm)

## Installation

> go get -u -v -f all

## Build

> go build

## Usage

> go run main.go

### API 호출 시 제한 조건 설정 방법
* page : Pagination의 페이지
* perPage : 1페이지당 레코드 수
* cond[column::operation]=value : 호출하는 칼럼 별 제한조건(<>= like 등)과 해당 값을 설정
* ex) /api/:api_nm?page=1&perPage=100&cond[ADDR::like]=XXX

## Meta

Promptechnology - [@Homepage](http://www.promptech.co.kr/) - [promptech@promptech.co.kr](promptech@promptech.co.kr)

프로젝트는 아파치 2.0 라이센스로 배포되었습니다. 자세한 사항은 ``LICENSE`` 파일을 확인해주세요.

Distributed under the Apache 2.0 license. See ``LICENSE`` for more information.

## Support
![alt text](http://wisepaip.org/assets/home/promptech-d8574a0910561aaea077bc759b1cf94c07baecc551f034ee9c7e830572d671de.png "Title Text")