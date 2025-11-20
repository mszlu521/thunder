package einos

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino-ext/components/model/claude"
	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino-ext/components/model/gemini"
	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino-ext/components/model/qianfan"
	"github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/cloudwego/eino/components/model"
)

const (
	OLLAMA   = "ollama"
	OPENAI   = "openai"
	CLAUDE   = "claude"
	DEEPSEEK = "deepseek"
	GEMINI   = "gemini"
	QIANFAN  = "qianfan"
	QWEN     = "qwen"
	ARK      = "ark"
)

type ChatModelConfig struct {
	OllamaConfig   *ollama.ChatModelConfig
	ArkConfig      *ark.ChatModelConfig
	ClaudeConfig   *claude.Config
	OpenaiConfig   *openai.ChatModelConfig
	DeepseekConfig *deepseek.ChatModelConfig
	QianfanConfig  *qianfan.ChatModelConfig
	QwenConfig     *qwen.ChatModelConfig
	GeminiConfig   *gemini.Config
}

func LoadChatModel(ctx context.Context, chatType string, config *ChatModelConfig) (model.ToolCallingChatModel, error) {
	var chatModel model.ToolCallingChatModel
	var err error
	if chatType == OLLAMA {
		chatModel, err = ollama.NewChatModel(ctx, config.OllamaConfig)
		return chatModel, err
	}
	if chatType == ARK {
		chatModel, err = ark.NewChatModel(ctx, config.ArkConfig)
		return chatModel, err
	}
	if chatType == CLAUDE {
		chatModel, err = claude.NewChatModel(ctx, config.ClaudeConfig)
		return chatModel, err
	}
	if chatType == OPENAI {
		chatModel, err = openai.NewChatModel(ctx, config.OpenaiConfig)
		return chatModel, err
	}
	if chatType == DEEPSEEK {
		chatModel, err = deepseek.NewChatModel(ctx, config.DeepseekConfig)
		return chatModel, err
	}
	if chatType == GEMINI {
		chatModel, err = gemini.NewChatModel(ctx, config.GeminiConfig)
		return chatModel, err
	}
	if chatType == QIANFAN {
		chatModel, err = qianfan.NewChatModel(ctx, config.QianfanConfig)
		return chatModel, err
	}
	if chatType == QWEN {
		chatModel, err = qwen.NewChatModel(ctx, config.QwenConfig)
		return chatModel, err
	}
	chatModel, err = openai.NewChatModel(ctx, config.OpenaiConfig)
	return chatModel, err
}
