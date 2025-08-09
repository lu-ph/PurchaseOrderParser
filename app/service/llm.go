package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/outputparser"
)

type LLMService struct {
	llm *openai.LLM
}

func InitLLM(modelName string, apiUrl string, apiKey string) (*LLMService, error) {
	llm, err := openai.New(
		openai.WithModel(modelName),
		openai.WithBaseURL(apiUrl),
		openai.WithToken(apiKey),
	)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	llmService := &LLMService{llm: llm}
	return llmService, nil
}

func (llmService *LLMService) GenerateAIResponse(
	ctx context.Context,
	input string,
	systemPrompt string,
	temperature float64,
	T any) (any, error) {

	parser, _ := outputparser.NewDefined(T)

	systemPromptWithParserStr := systemPrompt + parser.GetFormatInstructions()

	messages := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, systemPromptWithParserStr),
		llms.TextParts(llms.ChatMessageTypeHuman, input),
	}
	content, err := llmService.llm.GenerateContent(ctx, messages, llms.WithOptions(llms.CallOptions{Temperature: temperature}))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("AI响应内容: %s", content.Choices[0].Content)
	log.Printf("停止原因: %s", content.Choices[0].StopReason)

	result, err := parser.Parse(content.Choices[0].Content)
	if err != nil {
		log.Printf("parser解析JSON错误: %v", err)
		return "", fmt.Errorf("parser解析JSON错误: %v", err)
	}

	return result, nil
}

func AnyToStruct[T any](any interface{}) (T, error) {
	var t T
	jsonBytes, err := json.Marshal(any)
	if err != nil {
		return t, err
	}
	err = json.Unmarshal(jsonBytes, &t)
	return t, err
}
