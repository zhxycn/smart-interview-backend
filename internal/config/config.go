package config

import (
	"github.com/joho/godotenv"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Config struct {
	Debug            bool
	Port             string
	Database         string
	XunfeiAppId      string
	XunfeiApiKey     string
	XunfeiApiSecret  string
	DifyEndpoint     string
	ResumeApiKey     string
	InterviewApiKey  string
	TencentAppId     string
	TencentSecretId  string
	TencentSecretKey string
	SiliconflowToken string
	SiliconflowModel string
	SiliconflowVoice string
}

func DefaultConfig() *Config {
	return &Config{
		Debug:            false,
		Port:             "8080",
		Database:         "",
		XunfeiAppId:      "",
		XunfeiApiKey:     "",
		XunfeiApiSecret:  "",
		DifyEndpoint:     "",
		ResumeApiKey:     "",
		InterviewApiKey:  "",
		TencentAppId:     "",
		TencentSecretId:  "",
		TencentSecretKey: "",
		SiliconflowToken: "",
		SiliconflowModel: "",
		SiliconflowVoice: "",
	}
}

func LoadConfig() *Config {
	// 获取默认配置
	config := DefaultConfig()

	// 从 .env 文件加载配置
	_ = godotenv.Load()

	// 自动处理环境变量
	return LoadConfigFromEnv(config)
}

func LoadConfigFromEnv(config *Config) *Config {
	val := reflect.ValueOf(config).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := typ.Field(i).Name
		envName := fieldNameToEnvName(fieldName)
		envVal := os.Getenv(envName)

		if envVal == "" {
			continue
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(envVal)
		case reflect.Bool:
			field.SetBool(envVal == "true")
		case reflect.Int, reflect.Int64:
			if intVal, err := strconv.Atoi(envVal); err == nil {
				field.SetInt(int64(intVal))
			}
		}

	}

	return config
}

func fieldNameToEnvName(name string) string {
	var result strings.Builder
	for i, char := range name {
		if i > 0 && 'A' <= char && char <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(char)
	}
	return strings.ToUpper(result.String())
}
