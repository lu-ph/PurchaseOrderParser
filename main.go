package main

import (
	"purchase_order_parser/app/api"
	"purchase_order_parser/app/config"

	"github.com/gin-gonic/gin"
)

func main() {
	config := config.InitConfig()
	r := api.NewRouter(gin.Default())
	r.SetupRoutes(config)
	r.Run(":8086")
}
