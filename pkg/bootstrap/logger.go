package bootstrap

import (
	"os"

	fluentLogger "github.com/devexps/go-micro/log/fluent/v2"
	logrusLogger "github.com/devexps/go-micro/log/logrus/v2"
	zapLogger "github.com/devexps/go-micro/log/zap/v2"

	"github.com/devexps/go-bootstrap/api/gen/go/common/conf"

	"github.com/devexps/go-micro/v2/log"
	"github.com/devexps/go-micro/v2/middleware/tracing"

	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LoggerType string

const (
	LoggerTypeStd    LoggerType = "std"
	LoggerTypeFluent LoggerType = "fluent"
	LoggerTypeLogrus LoggerType = "logrus"
	LoggerTypeZap    LoggerType = "zap"
)

// NewLoggerProvider creates a new logger provider
func NewLoggerProvider(cfg *conf.Logger, serviceInfo *ServiceInfo) log.Logger {
	l := NewLogger(cfg)

	return log.With(
		l,
		"service.id", serviceInfo.Id,
		"service.name", serviceInfo.Name,
		"service.version", serviceInfo.Version,
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"trace_id", tracing.TraceID(),
		"span_id", tracing.SpanID(),
	)
}

// NewLogger creates a new logger
func NewLogger(cfg *conf.Logger) log.Logger {
	if cfg == nil {
		return NewStdLogger()
	}

	switch LoggerType(cfg.Type) {
	default:
		fallthrough
	case LoggerTypeStd:
		return NewStdLogger()
	case LoggerTypeFluent:
		return NewFluentLogger(cfg)
	case LoggerTypeLogrus:
		return NewLogrusLogger(cfg)
	case LoggerTypeZap:
		return NewZapLogger(cfg)
	}
}

// NewStdLogger creates a new logger - GoMicro built-in, console output
func NewStdLogger() log.Logger {
	return log.NewStdLogger(os.Stdout)
}

// NewFluentLogger creates a new logger - Fluent
func NewFluentLogger(cfg *conf.Logger) log.Logger {
	wrapped, err := fluentLogger.NewLogger(cfg.Fluent.Endpoint)
	if err != nil {
		panic("create fluent logger failed")
		return nil
	}
	return wrapped
}

// NewLogrusLogger creates a new logger - Logrus
func NewLogrusLogger(cfg *conf.Logger) log.Logger {
	loggerLevel, err := logrus.ParseLevel(cfg.Logrus.Level)
	if err != nil {
		loggerLevel = logrus.InfoLevel
	}

	var loggerFormatter logrus.Formatter
	switch cfg.Logrus.Formatter {
	default:
		fallthrough
	case "text":
		loggerFormatter = &logrus.TextFormatter{
			DisableColors:    cfg.Logrus.DisableColors,
			DisableTimestamp: cfg.Logrus.DisableTimestamp,
			TimestampFormat:  cfg.Logrus.TimestampFormat,
		}
		break
	case "json":
		loggerFormatter = &logrus.JSONFormatter{
			DisableTimestamp: cfg.Logrus.DisableTimestamp,
			TimestampFormat:  cfg.Logrus.TimestampFormat,
		}
		break
	}

	logger := logrus.New()
	logger.Level = loggerLevel
	logger.Formatter = loggerFormatter

	wrapped := logrusLogger.NewLogger(logger)
	return wrapped
}

// NewZapLogger creates a new logger - Zap
func NewZapLogger(cfg *conf.Logger) log.Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)

	lumberJackLogger := &lumberjack.Logger{
		Filename:   cfg.Zap.Filename,
		MaxSize:    int(cfg.Zap.MaxSize),
		MaxBackups: int(cfg.Zap.MaxBackups),
		MaxAge:     int(cfg.Zap.MaxAge),
	}
	writeSyncer := zapcore.AddSync(lumberJackLogger)

	var lvl = new(zapcore.Level)
	if err := lvl.UnmarshalText([]byte(cfg.Zap.Level)); err != nil {
		return nil
	}

	core := zapcore.NewCore(jsonEncoder, writeSyncer, lvl)
	logger := zap.New(core).WithOptions()

	wrapped := zapLogger.NewLogger(logger)
	return wrapped
}
