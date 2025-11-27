package logs

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/mszlu521/thunder/config"
)

// 定义一个私有的全局 logger 实例
var defaultLogger *slog.Logger

// 为了在 context 中传递 logger，我们定义一个私有的 key 类型
type loggerKey struct{}

// 彩色输出常量
var (
	reset = "\033[0m"
	red   = "\033[31m"
	green = "\033[32m"
	yellow = "\033[33m"
	blue  = "\033[34m"
	magenta = "\033[35m"
	cyan  = "\033[36m"
)

// 彩色文本结构
type coloredText struct {
	color string
	text  string
}

// 自定义处理器，用于美化本地日志输出
type prettyHandler struct {
	opts slog.HandlerOptions
	w    io.Writer
}

// 实现 slog.Handler 接口
func (h *prettyHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

func (h *prettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &prettyHandler{opts: h.opts, w: h.w}
}

func (h *prettyHandler) WithGroup(name string) slog.Handler {
	return &prettyHandler{opts: h.opts, w: h.w}
}

func (h *prettyHandler) Handle(ctx context.Context, r slog.Record) error {
	// 时间戳
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")

	// 获取源代码位置
	var source string
	if h.opts.AddSource && r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		// 简化文件路径，只显示相对路径
		file := f.File
		if idx := strings.Index(file, "agents-workflows"); idx != -1 {
			file = file[idx:]
		} else {
			// 如果找不到项目根目录，尝试简化路径
			parts := strings.Split(file, string(os.PathSeparator))
			if len(parts) > 2 {
				file = strings.Join(parts[len(parts)-2:], string(os.PathSeparator))
			}
		}
		source = fmt.Sprintf("%s:%d", file, f.Line)
	}

	// 根据日志级别设置颜色和标签
	var levelText coloredText
	switch r.Level {
	case slog.LevelDebug:
		levelText = coloredText{color: cyan, text: "DEBUG"}
	case slog.LevelInfo:
		levelText = coloredText{color: green, text: "INFO "}
	case slog.LevelWarn:
		levelText = coloredText{color: yellow, text: "WARN "}
	case slog.LevelError:
		levelText = coloredText{color: red, text: "ERROR"}
	default:
		levelText = coloredText{color: reset, text: r.Level.String()}
	}

	// 构建输出行
	var sb strings.Builder
	sb.WriteString(levelText.color)
	sb.WriteString("[")
	sb.WriteString(timestamp)
	sb.WriteString("] ")
	sb.WriteString(levelText.text)
	sb.WriteString(reset)
	sb.WriteString(" ")
	sb.WriteString(r.Message)

	// 添加源代码位置（如果启用）
	if source != "" {
		sb.WriteString(" ")
		sb.WriteString(magenta)
		sb.WriteString("(")
		sb.WriteString(source)
		sb.WriteString(")")
		sb.WriteString(reset)
	}

	// 添加键值对
	r.Attrs(func(attr slog.Attr) bool {
		sb.WriteString(" ")
		sb.WriteString(blue)
		sb.WriteString(attr.Key)
		sb.WriteString(reset)
		sb.WriteString("=")
		sb.WriteString(formatValue(attr.Value))
		return true
	})

	sb.WriteString("\n")
	_, err := h.w.Write([]byte(sb.String()))
	return err
}

// 格式化 slog.Value 为字符串
func formatValue(v slog.Value) string {
	switch v.Kind() {
	case slog.KindString:
		return v.String()
	case slog.KindInt64:
		return strconv.FormatInt(v.Int64(), 10)
	case slog.KindUint64:
		return strconv.FormatUint(v.Uint64(), 10)
	case slog.KindFloat64:
		return strconv.FormatFloat(v.Float64(), 'g', -1, 64)
	case slog.KindBool:
		return strconv.FormatBool(v.Bool())
	case slog.KindDuration:
		return v.Duration().String()
	case slog.KindTime:
		return v.Time().Format(time.RFC3339)
	default:
		return fmt.Sprintf("%+v", v.Any())
	}
}

// Init 初始化全局日志记录器
// 这是应用启动时应该调用的第一个函数
func Init(c *config.LogConfig) {
	if c == nil {
		return
	}
	var level slog.Level
	// 解析日志级别字符串
	switch c.GetLevel() {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo // 默认为 info
	}
	opts := &slog.HandlerOptions{
		AddSource: c.GetAddSource(),
		Level:     level,
	}

	var handler slog.Handler
	output := c.Output
	if output == nil {
		output = os.Stdout // 默认为标准输出
	}

	// 选择日志格式
	if c.GetFormat() == "json" {
		handler = slog.NewJSONHandler(output, opts)
	} else if c.GetFormat() == "pretty" {
		// 使用我们自定义的美化处理器
		handler = &prettyHandler{opts: *opts, w: output}
	} else {
		// 默认使用文本处理器
		handler = slog.NewTextHandler(output, opts)
	}

	// 创建并设置默认 logger
	defaultLogger = slog.New(handler)
	slog.SetDefault(defaultLogger)

	// 将标准库 logs 的输出重定向到 slog
	// 这可以捕获那些使用标准 `logs` 包的第三方库的日志
	log.SetFlags(0)
	log.SetOutput(slog.NewLogLogger(handler, slog.LevelInfo).Writer())
}

// ----- 包级别的便捷函数 -----

// Debug 记录 debug 级别的日志
func Debug(msg string, args ...any) {
	defaultLogger.Debug(msg, args...)
}

// Info 记录 info 级别的日志
func Info(msg string, args ...any) {
	defaultLogger.Info(msg, args...)
}

// Warn 记录 warn 级别的日志
func Warn(msg string, args ...any) {
	defaultLogger.Warn(msg, args...)
}

// Error 记录 error 级别的日志
func Error(msg string, args ...any) {
	defaultLogger.Error(msg, args...)
}

// Debugf 记录格式化的 debug 级别日志
func Debugf(format string, args ...any) {
	defaultLogger.Debug(fmt.Sprintf(format, args...))
}

// Infof 记录格式化的 info 级别日志
func Infof(format string, args ...any) {
	defaultLogger.Info(fmt.Sprintf(format, args...))
}

// Warnf 记录格式化的 warn 级别日志
func Warnf(format string, args ...any) {
	defaultLogger.Warn(fmt.Sprintf(format, args...))
}

// Errorf 记录格式化的 error 级别日志
func Errorf(format string, args ...any) {
	defaultLogger.Error(fmt.Sprintf(format, args...))
}

// ----- 上下文处理函数 -----

// WithContext 将一个带有附加属性的 logger 存入 context 并返回新的 context
// 在中间件中使用，为每个请求添加 request_id 等信息
func WithContext(ctx context.Context, args ...any) context.Context {
	l := FromContext(ctx).With(args...)
	return context.WithValue(ctx, loggerKey{}, l)
}

// FromContext 从 context 中获取 logger
// 如果 context 中没有，则返回全局默认 logger
func FromContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(loggerKey{}).(*slog.Logger); ok {
		return l
	}
	return defaultLogger
}

// CtxDebug 从 context 中获取 logger 并记录 debug 日志
func CtxDebug(ctx context.Context, msg string, args ...any) {
	FromContext(ctx).Debug(msg, args...)
}

// CtxInfo 从 context 中获取 logger 并记录 info 日志
func CtxInfo(ctx context.Context, msg string, args ...any) {
	FromContext(ctx).Info(msg, args...)
}

// CtxWarn 从 context 中获取 logger 并记录 warn 日志
func CtxWarn(ctx context.Context, msg string, args ...any) {
	FromContext(ctx).Warn(msg, args...)
}

// CtxError 从 context 中获取 logger 并记录 error 日志
func CtxError(ctx context.Context, msg string, args ...any) {
	FromContext(ctx).Error(msg, args...)
}