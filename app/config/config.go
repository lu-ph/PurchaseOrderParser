package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	LLM struct {
		ModelName string `yaml:"model_name"`
		ApiUrl    string `yaml:"api_url"`
		ApiKey    string `yaml:"api_key"`
	} `yaml:"llm"`
	BaiDuOCR struct {
		OCRApiKey    string `yaml:"ocr_api_key"`
		OCRSecretKey string `yaml:"ocr_secret_key"`
	} `yaml:"ocr"`
}

func InitConfig() Config {
	configData, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("配置文件读取失败: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		log.Fatalf("配置文件解析失败: %v", err)
	}

	return config
}
