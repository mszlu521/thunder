package einos

import (
	"context"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino-ext/components/embedding/dashscope"
	"github.com/cloudwego/eino-ext/components/embedding/gemini"
	"github.com/cloudwego/eino-ext/components/embedding/ollama"
	"github.com/cloudwego/eino-ext/components/embedding/openai"
	"github.com/cloudwego/eino-ext/components/embedding/qianfan"
	"github.com/cloudwego/eino-ext/components/embedding/tencentcloud"
	"github.com/cloudwego/eino/components/embedding"
)

const (
	EmbeddingOllama       = "ollama"
	EmbeddingDashscope    = "dashscope"
	EmbeddingOpenai       = "openai"
	EmbeddingArk          = "ark"
	EmbeddingQianfan      = "qianfan"
	EmbeddingTencentCloud = "tencentcloud"
	EmbeddingGemini       = "gemini"
)

type EmbeddingModelConfig struct {
	OllamaConfig       *ollama.EmbeddingConfig
	DashscopeConfig    *dashscope.EmbeddingConfig
	ArkConfig          *ark.EmbeddingConfig
	QianfanConfig      *qianfan.EmbeddingConfig
	TencentCloudConfig *tencentcloud.EmbeddingConfig
	OpenaiConfig       *openai.EmbeddingConfig
	GeminiConfig       *gemini.EmbeddingConfig
}

func LoadEmbedding(ctx context.Context, embeddingType string, config *EmbeddingModelConfig) (embedding.Embedder, error) {
	var embedder embedding.Embedder
	var err error
	switch embeddingType {
	case EmbeddingOllama:
		embedder, err = ollama.NewEmbedder(ctx, config.OllamaConfig)
	case EmbeddingDashscope:
		embedder, err = dashscope.NewEmbedder(ctx, config.DashscopeConfig)
	case EmbeddingArk:
		embedder, err = ark.NewEmbedder(ctx, config.ArkConfig)
	case EmbeddingQianfan:
		embedder, err = qianfan.NewEmbedder(ctx, config.QianfanConfig)
	case EmbeddingTencentCloud:
		embedder, err = tencentcloud.NewEmbedder(ctx, config.TencentCloudConfig)
	case EmbeddingOpenai:
		embedder, err = openai.NewEmbedder(ctx, config.OpenaiConfig)
	case EmbeddingGemini:
		embedder, err = gemini.NewEmbedder(ctx, config.GeminiConfig)
	default:
		embedder, err = openai.NewEmbedder(ctx, config.OpenaiConfig)
	}
	return embedder, err
}
