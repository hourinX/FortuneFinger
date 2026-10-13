package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 日志打印模式
const (
	ModeConsole = "console" // stdout
	ModeFile    = "file"    // path
)

// 日志格式
const (
	EncodingJSON  = "json"
	EncodingPlain = "plain"
)

// 最大日志文件大小，单位MB
const defaultMaxSizeMB = 100

type Config struct {
	ServiceName string `yaml:"ServiceName"` // 服务名
	Mode        string `yaml:"Mode"`        // 日志模式
	Path        string `yaml:"Path"`        // 日志目录，file模式必填
	Encoding    string `yaml:"Encoding"`    // 日志格式
	Level       string `yaml:"Level"`       // 日志等级
	Compress    bool   `yaml:"Compress"`    // 是否压缩
	KeepDays    int    `yaml:"KeepDays"`    // 保留天数，file模式必填
}

// 兜底结构体
func DefaultConfig() Config {
	return Config{
		ServiceName: "service",
		Mode:        ModeConsole,
		Path:        "./logs",
		Encoding:    EncodingJSON,
		Level:       "info",
	}
}

// 构建logger
func New(conf Config) (*zap.Logger, error) {
	conf = conf.withDefaults()

	level, err := parseLevel(conf.Level)
	if err != nil {
		return nil, err
	}
	enc, err := newEncoder(conf.Encoding)
	if err != nil {
		return nil, err
	}
	sink, err := newSink(conf)
	if err != nil {
		return nil, err
	}

	core := newMaskCore(zapcore.NewCore(enc, sink, level))

	return zap.New(
		core,
		zap.AddCaller(),
		zap.Fields(zap.String("service", conf.ServiceName)),
	), nil
}

// 构建logger配置
func (c Config) withDefaults() Config {
	c.Mode = strings.ToLower(strings.TrimSpace(c.Mode))
	c.Encoding = strings.ToLower(strings.TrimSpace(c.Encoding))
	c.Level = strings.ToLower(strings.TrimSpace(c.Level))

	d := DefaultConfig()
	if c.ServiceName == "" {
		c.ServiceName = d.ServiceName
	}
	if c.Mode == "" {
		c.Mode = d.Mode
	}
	if c.Path == "" {
		c.Path = d.Path
	}
	if c.Encoding == "" {
		c.Encoding = d.Encoding
	}
	if c.Level == "" {
		c.Level = d.Level
	}
	return c
}

// 构建logger日志级别
func parseLevel(s string) (zapcore.Level, error) {
	switch s {
	case "debug":
		return zapcore.DebugLevel, nil
	case "info":
		return zapcore.InfoLevel, nil
	case "warn":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	default:
		return zapcore.InvalidLevel, fmt.Errorf("logger: 非法的日志级别「%s」，请使用 debug/info/warn/error", s)
	}
}

// 构建logger编码器
func newEncoder(encoding string) (zapcore.Encoder, error) {
	ec := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	switch encoding {
	case EncodingJSON:
		return zapcore.NewJSONEncoder(ec), nil
	case EncodingPlain:
		return zapcore.NewConsoleEncoder(ec), nil
	default:
		return nil, fmt.Errorf("logger: 非法的日志编码「%s」，请使用 json/plain", encoding)
	}
}

// 构建logger输出目标
func newSink(conf Config) (zapcore.WriteSyncer, error) {
	switch conf.Mode {
	case ModeConsole:
		return zapcore.AddSync(io.Writer(os.Stdout)), nil

	case ModeFile:
		if conf.KeepDays <= 0 {
			return nil, fmt.Errorf(
				"logger: file模式下KeepDays为0时旧日志将永不清理，当前值为「%d」", conf.KeepDays)
		}
		if err := os.MkdirAll(conf.Path, 0o750); err != nil {
			return nil, fmt.Errorf("logger: 创建日志目录失败「%s」: %w", conf.Path, err)
		}
		path := filepath.Join(conf.Path, conf.ServiceName+".log")

		probe, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
		if err != nil {
			return nil, fmt.Errorf("logger: 打开日志文件失败「%s」: %w", path, err)
		}
		if err := probe.Close(); err != nil {
			return nil, fmt.Errorf("logger: 关闭日志文件失败「%s」: %w", path, err)
		}

		return zapcore.AddSync(&lumberjack.Logger{
			Filename:  path,
			MaxSize:   defaultMaxSizeMB,
			MaxAge:    conf.KeepDays,
			Compress:  conf.Compress,
			LocalTime: true,
		}), nil

	default:
		return nil, fmt.Errorf("logger: 非法的日志模式「%s」，请使用 console/file", conf.Mode)
	}
}

// 信息脱敏档位
type maskKind int

const (
	maskNone maskKind = iota
	maskFull
	maskPartial
)

// 全脱敏屏蔽词
var fullMaskTerms = []string{
	"password", "passwd", "pwd", "secret", "token",
	"credential", "privatekey", "accesskey", "apikey", "authorization",
}

// 非全脱敏屏蔽词
var partialMaskTerms = []string{
	"mobile", "phone", "idcard", "bankcard",
}

// 脱敏档位判断
func keyMaskKind(key string) maskKind {
	k := normalizeKey(key)
	if k == "" {
		return maskNone
	}
	for _, term := range fullMaskTerms {
		if strings.Contains(k, term) {
			return maskFull
		}
	}
	for _, term := range partialMaskTerms {
		if strings.Contains(k, term) {
			return maskPartial
		}
	}
	return maskNone
}

func normalizeKey(key string) string {
	var b strings.Builder
	b.Grow(len(key))
	for _, r := range key {
		switch r {
		case '_', '-', '.', ' ', '\t':
			continue
		}
		b.WriteRune(unicode.ToLower(r))
	}
	return b.String()
}

// 脱敏值
func maskValue(kind maskKind, v string) string {
	if kind == maskFull {
		return "***"
	}
	r := []rune(v)
	if len(r) < 8 {
		return "***"
	}
	return string(r[:3]) + "****" + string(r[len(r)-4:])
}

// 字段脱敏处理
func maskFields(fields []zapcore.Field) []zapcore.Field {
	sensitive := false
	for i := range fields {
		if keyMaskKind(fields[i].Key) != maskNone {
			sensitive = true
			break
		}
	}
	if !sensitive {
		return fields
	}

	out := make([]zapcore.Field, len(fields))
	copy(out, fields)
	for i := range out {
		kind := keyMaskKind(out[i].Key)
		if kind == maskNone {
			continue
		}
		if out[i].Type == zapcore.StringType {
			out[i] = zapcore.Field{
				Key:    out[i].Key,
				Type:   zapcore.StringType,
				String: maskValue(kind, out[i].String),
			}
			continue
		}
		out[i] = zapcore.Field{Key: out[i].Key, Type: zapcore.StringType, String: "***"}
	}
	return out
}

type maskCore struct{ zapcore.Core }

func newMaskCore(c zapcore.Core) zapcore.Core { return maskCore{c} }

func (c maskCore) With(fields []zapcore.Field) zapcore.Core {
	return maskCore{c.Core.With(maskFields(fields))}
}

// 日志等级启用校验
func (c maskCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

func (c maskCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	return c.Core.Write(ent, maskFields(fields))
}

// 字段脱敏处理
func Sanitize(fields map[string]any) map[string]any {
	out := make(map[string]any, len(fields))
	for k, v := range fields {
		kind := keyMaskKind(k)
		if kind == maskNone {
			out[k] = v
			continue
		}
		if s, ok := v.(string); ok {
			out[k] = maskValue(kind, s)
			continue
		}
		out[k] = "***"
	}
	return out
}
