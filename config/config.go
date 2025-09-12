package config

import (
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	AI       AIConfig
}

type ServerConfig struct {
	Port string
	Host string
}

type DatabaseConfig struct {
	Type     string // "sqlite" or "mongodb"
	Path     string // SQLite database file path
	MongoURI string // MongoDB connection URI
	MongoDB  string // MongoDB database name
}

type JWTConfig struct {
	SecretKey string
	ExpiresIn int // hours
}

// AIProvider 定义AI服务提供商类型
type AIProvider string

const (
	ProviderHuggingFace AIProvider = "huggingface"
	ProviderGroq        AIProvider = "groq"
	ProviderSimulation  AIProvider = "simulation"
	ProviderGemini      AIProvider = "gemini"
)

// AIConfig AI服务配置
type AIConfig struct {
	PrimaryProvider   AIProvider `json:"primary_provider"`
	FallbackProvider  AIProvider `json:"fallback_provider"`
	SimulationProvider AIProvider `json:"simulation_provider"`
	SwitchThreshold   float64    `json:"switch_threshold"`
	RequestTimeout    int        `json:"request_timeout"`
	HuggingFace       HuggingFaceConfig `json:"huggingface"`
	Groq              GroqConfig        `json:"groq"`
	Gemini            GeminiConfig      `json:"gemini"`
}

// HuggingFaceConfig Hugging Face API配置
type HuggingFaceConfig struct {
	APIURL     string `json:"api_url"`
	APIToken   string `json:"api_token"`
	Model      string `json:"model"`
	MaxTokens  int    `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
	DailyLimit int    `json:"daily_limit"`
}

// GroqConfig Groq API配置
type GroqConfig struct {
	APIURL     string `json:"api_url"`
	APIKey     string `json:"api_key"`
	Model      string `json:"model"`
	MaxTokens  int    `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
	DailyLimit int    `json:"daily_limit"`
}

// GeminiConfig Gemini API配置
type GeminiConfig struct {
	APIURL     string `json:"api_url"`
	APIKey     string `json:"api_key"`
	Model      string `json:"model"`
	MaxTokens  int    `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
	DailyLimit int    `json:"daily_limit"`
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Host: getEnv("HOST", "0.0.0.0"),
		},
		Database: DatabaseConfig{
			Type:     getEnv("DATABASE_TYPE", "sqlite"),
			Path:     getEnv("DB_PATH", "data/app.db"),
			MongoURI: getEnv("MONGODB_URI", ""),
			MongoDB:  getEnv("MONGODB_DATABASE", "chatbot"),
		},
		JWT: JWTConfig{
			SecretKey: getEnv("JWT_SECRET", "your-secret-key"),
			ExpiresIn: getEnvAsInt("JWT_EXPIRES_IN", 24),
		},
		AI: AIConfig{
			PrimaryProvider:   AIProvider(getEnv("AI_PRIMARY_PROVIDER", "groq")),
			FallbackProvider:  AIProvider(getEnv("AI_FALLBACK_PROVIDER", "gemini")),
			SimulationProvider: AIProvider(getEnv("AI_SIMULATION_PROVIDER", "simulation")),
			SwitchThreshold:   getEnvAsFloat("AI_SWITCH_THRESHOLD", 0.8),
			RequestTimeout:    getEnvAsInt("AI_REQUEST_TIMEOUT", 30),
			HuggingFace: HuggingFaceConfig{
				APIURL:      getEnv("HF_API_URL", "https://api-inference.huggingface.co/models/microsoft/DialoGPT-small"),
				APIToken:    getEnv("HF_API_TOKEN", ""),
				Model:       getEnv("HF_MODEL", "microsoft/DialoGPT-small"),
				MaxTokens:   getEnvAsInt("HF_MAX_TOKENS", 100),
				Temperature: getEnvAsFloat("HF_TEMPERATURE", 0.7),
				DailyLimit:  getEnvAsInt("HF_DAILY_LIMIT", 1000),
			},
			Groq: GroqConfig{
				APIURL:      getEnv("GROQ_API_URL", "https://api.groq.com/openai/v1/chat/completions"),
				APIKey:      getEnv("GROQ_API_KEY", ""),
				Model:       getEnv("GROQ_MODEL", "llama-3.1-8b-instant"),
				MaxTokens:   getEnvAsInt("GROQ_MAX_TOKENS", 100),
				Temperature: getEnvAsFloat("GROQ_TEMPERATURE", 0.7),
				DailyLimit:  getEnvAsInt("GROQ_DAILY_LIMIT", 10000),
			},
			Gemini: GeminiConfig{
				APIURL:      getEnv("GEMINI_API_URL", "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent"),
				APIKey:      getEnv("GEMINI_API_KEY", ""),
				Model:       getEnv("GEMINI_MODEL", "gemini-2.0-flash"),
				MaxTokens:   getEnvAsInt("GEMINI_MAX_TOKENS", 100),
				Temperature: getEnvAsFloat("GEMINI_TEMPERATURE", 0.7),
				DailyLimit:  getEnvAsInt("GEMINI_DAILY_LIMIT", 1500),
			},
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}
