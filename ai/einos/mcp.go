package einos

import (
	"context"
	"fmt"

	mcpp "github.com/cloudwego/eino-ext/components/tool/mcp"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
)

type McpConfig struct {
	BaseUrl string
	Token   string
	Name    string
	Version string
}

func GetEinoBaseTools(ctx context.Context, config *McpConfig) ([]tool.BaseTool, error) {
	headers := make(map[string]string)
	if config.Token != "" {
		headers["Authorization"] = fmt.Sprintf("Bearer %s", config.Token)
	}
	options := transport.WithHeaders(headers)
	cli, err := client.NewSSEMCPClient(config.BaseUrl, options)
	if err != nil {
		return nil, err
	}
	err = cli.Start(ctx)
	if err != nil {
		return nil, err
	}

	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    config.Name,
		Version: config.Version,
	}

	_, err = cli.Initialize(ctx, initRequest)
	if err != nil {
		return nil, err
	}
	tools, err := mcpp.GetTools(ctx, &mcpp.Config{Cli: cli})
	if err != nil {
		return nil, err
	}

	return tools, nil
}

func GetMCPTool(ctx context.Context, config *McpConfig) ([]mcp.Tool, error) {
	headers := make(map[string]string)
	if config.Token != "" {
		headers["Authorization"] = fmt.Sprintf("Bearer %s", config.Token)
	}
	options := transport.WithHeaders(headers)
	cli, err := client.NewSSEMCPClient(config.BaseUrl, options)
	if err != nil {
		return nil, err
	}
	err = cli.Start(ctx)
	if err != nil {
		return nil, err
	}

	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    config.Name,
		Version: config.Version,
	}

	_, err = cli.Initialize(ctx, initRequest)
	if err != nil {
		return nil, err
	}
	tools, err := cli.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil {
		return nil, err
	}

	return tools.Tools, nil
}

func GetMCPToolAndCli(ctx context.Context, config *McpConfig) ([]mcp.Tool, *client.Client, error) {
	headers := make(map[string]string)
	if config.Token != "" {
		headers["Authorization"] = fmt.Sprintf("Bearer %s", config.Token)
	}
	options := transport.WithHeaders(headers)
	cli, err := client.NewSSEMCPClient(config.BaseUrl, options)
	if err != nil {
		return nil, nil, err
	}
	err = cli.Start(ctx)
	if err != nil {
		return nil, nil, err
	}

	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    config.Name,
		Version: config.Version,
	}

	_, err = cli.Initialize(ctx, initRequest)
	if err != nil {
		return nil, nil, err
	}
	tools, err := cli.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil {
		return nil, nil, err
	}

	return tools.Tools, cli, nil
}

// ConvertSchema converts mcp.ToolInputSchema to a map of ParameterInfo
func ConvertSchema(inputSchema mcp.ToolInputSchema) map[string]*schema.ParameterInfo {
	params := make(map[string]*schema.ParameterInfo)

	// 遍历schema中的属性
	for key, value := range inputSchema.Properties {
		paramInfo := &schema.ParameterInfo{}

		// 解析属性信息
		if propMap, ok := value.(map[string]any); ok {
			// 处理类型
			if typeVal, exists := propMap["type"]; exists {
				if typeStr, ok := typeVal.(string); ok {
					// 需要将MCP类型映射到Eino类型
					switch typeStr {
					case "string":
						paramInfo.Type = schema.String
					case "integer":
						paramInfo.Type = schema.Integer
					case "number":
						paramInfo.Type = schema.Number
					case "boolean":
						paramInfo.Type = schema.Boolean
					case "array":
						paramInfo.Type = schema.Array
					case "object":
						paramInfo.Type = schema.Object
					default:
						paramInfo.Type = schema.String // 默认为字符串类型
					}
				}
			}

			// 处理描述
			if descVal, exists := propMap["description"]; exists {
				if descStr, ok := descVal.(string); ok {
					paramInfo.Desc = descStr
				}
			}

			// 处理枚举值
			if enumVal, exists := propMap["enum"]; exists {
				if enumSlice, ok := enumVal.([]interface{}); ok {
					for _, enumItem := range enumSlice {
						if enumStr, ok := enumItem.(string); ok {
							paramInfo.Enum = append(paramInfo.Enum, enumStr)
						}
					}
				}
			}
		}

		// 检查是否为必填字段
		for _, requiredField := range inputSchema.Required {
			if requiredField == key {
				paramInfo.Required = true
				break
			}
		}

		params[key] = paramInfo
	}

	return params
}
