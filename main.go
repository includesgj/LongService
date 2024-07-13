package main

import (
	sdb "GinProject12/databases"
	utssh "GinProject12/util/ssh"
	"github.com/gin-gonic/gin"
)

//	@title			龙芯服务 API
//	@version		1.0
//	@description	这是软件杯B1赛题
//	@termsOfService	http://swagger.io/terms/ 使用条款

//	@contact.name	LoveSong.
//	@contact.url	无
//	@contact.email	3130250166@qq.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		39.99.139.249:8080

//	@securityDefinitions.basic	BasicAuth

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	sdb.GetDm()
	defer sdb.CloseDm()
	go utssh.Heartbeat()

	r := gin.Default()
	r = CollectRoute(r)
	panic(r.Run())
}

// put /Users/songguanju/RearEndCode/GolandProjects/GinProject12/hello

// put /Users/songguanju/RearEndCode/GolandProjects
