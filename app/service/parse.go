package service

import (
	"context"
	"fmt"
	"log"
	"purchase_order_parser/app/config"
	"purchase_order_parser/app/dto"
)

func Parse(ocrService *OCRService, llmService *LLMService, base64 string, fileType int, needRawText bool) (dto.LLMOutput, string, error) {
	text, err := ocrService.RecognizeFileWithTableAndText(base64, fileType)
	if err != nil {
		return dto.LLMOutput{}, "", err
	}
	outputParser := dto.LLMOutput{}
	outputStruct, err := llmService.GenerateAIResponse(
		context.Background(),
		text,
		config.PARSE_PROMPT,
		0.3,
		outputParser,
	)
	if err != nil {
		return dto.LLMOutput{}, "", err
	}
	outputParser, err = AnyToStruct[dto.LLMOutput](outputStruct)
	if err != nil {
		log.Printf("json.Unmarshal错误: %v", err)
		return dto.LLMOutput{}, "", fmt.Errorf("json.Unmarshal错误: %v", err)
	}
	if !needRawText {
		return outputParser, "", nil
	}
	return outputParser, text, nil
}
