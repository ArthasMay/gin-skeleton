package main

import (
	"goskeleton/app/utils/yml_config"
	_ "goskeleton/bootstrap"
	"goskeleton/routers"
)

func main() {
	router := routers.InitApiRouter()
	_ = router.Run(yml_config.CreateYamlFactory().GetString("HttpServer.Web.Port"))
}