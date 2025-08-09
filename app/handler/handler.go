package handler

import (
	"log"
	"net/http"
	"purchase_order_parser/app/dto"
	"purchase_order_parser/app/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	ocrApiKey    string
	ocrSecretKey string
	ocrService   *service.OCRService
	llmApiKey    string
	llmApiURL    string
	llmService   *service.LLMService
}

func NewHandler(ocrApiKey, ocrSecretKey, modelName, llmApiKey, llmApiURL string) *Handler {
	ocrService := service.NewOCRService(ocrApiKey, ocrSecretKey)
	llmService, err := service.InitLLM(modelName, llmApiURL, llmApiKey)
	if err != nil {
		log.Fatal(err)
	}
	return &Handler{ocrApiKey, ocrSecretKey, ocrService, llmApiKey, llmApiURL, llmService}
}

func (h *Handler) ParseHandler(c *gin.Context) {
	var req dto.Request
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Output{LLMOutput: dto.LLMOutput{}, Message: "Invalid request"})
		return
	}
	if len(req.Base64) == 0 {
		c.JSON(http.StatusBadRequest, dto.Output{LLMOutput: dto.LLMOutput{}, Message: "缺少base64字段"})
		return
	}
	if req.FileType == 0 {
		c.JSON(http.StatusBadRequest, dto.Output{LLMOutput: dto.LLMOutput{}, Message: "缺少fileType字段"})
		return
	}
	output, rawText, err := service.Parse(h.ocrService, h.llmService, req.Base64, req.FileType, req.NeedRawText)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Output{RawText: "", Message: "failed"})
		return
	}
	c.JSON(http.StatusOK, dto.Output{RawText: rawText, LLMOutput: output, Message: "success"})

}
