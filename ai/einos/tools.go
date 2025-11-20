package einos

import (
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

type InvokeParamTool interface {
	tool.InvokableTool
	Params() map[string]*schema.ParameterInfo
}
