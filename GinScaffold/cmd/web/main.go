package main

import (
	"goskeleton/app/utils/yml_config"
	_ "goskeleton/bootstrap"
	"goskeleton/routers"
)

func main() {
	router := routers.InitWebRouter()
	_ = router.Run(yml_config.CreateYamlFactory().GetString("HttpServer.Web.Port"))
}