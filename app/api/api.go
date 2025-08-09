package api

import (
	"purchase_order_parser/app/config"
	"purchase_order_parser/app/handler"

	"github.com/gin-gonic/gin"
)

type Router struct {
	*gin.Engine

	handler *handler.Handler
}

func NewRouter(engine *gin.Engine) *Router {
	return &Router{Engine: engine}
}

func (r *Router) SetupRoutes(config config.Config) {
	r.handler = handler.NewHandler(config.BaiDuOCR.OCRApiKey, config.BaiDuOCR.OCRSecretKey, config.LLM.ModelName, config.LLM.ApiKey, config.LLM.ApiUrl)

	r.POST("/parse", r.handler.ParseHandler)
}
