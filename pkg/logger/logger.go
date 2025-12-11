package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	otelLog "go.opentelemetry.io/otel/log"
	otelLogSdk "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Key используется для извлечения значений из context
type Key string

const (
	// Context keys для извлечения метаданных из контекста
	traceIDKey Key = "trace_id"
	userIDKey  Key = "user_id"

	// Timeout для graceful shutdown
	shutdownTimeout = 5 * time.Second
)

var (
	globalLogger *logger                    // Глобальный экземпляр логгера
	initOnce     sync.Once                  // Гарантирует единственную инициализацию
	dynamicLevel zap.AtomicLevel            // Динамический уровень логирования
	otelProvider *otelLogSdk.LoggerProvider // OpenTelemetry provider для graceful shutdown
)

// LoggerConfig содержит конфигурацию для инициализации логгера
type LoggerConfig struct {
	ServiceName  string // Имя сервиса (например, "order-service")
	Environment  string // Окружение (dev, staging, production)
	OTLPEndpoint string // Endpoint OpenTelemetry Collector (например, "localhost:4317")
}

type logger struct {
	zapLogger *zap.Logger
}

// Init инициализирует глобальный логгер
// levelStr - уровень логирования (debug, info, warn, error)
// asJSON - использовать JSON формат вместо консольного
// config - конфигурация для OpenTelemetry (если nil, OTLP отключен)
func Init(levelStr string, asJSON bool, config *LoggerConfig) error {
	return InitWithFile(levelStr, asJSON, "", config)
}

// InitWithFile инициализирует логгер с записью в файл
// logFile - путь к файлу логов (если пустой - пишем в stdout)
func InitWithFile(levelStr string, asJSON bool, logFile string, config *LoggerConfig) error {
	var initErr error

	initOnce.Do(func() {
		dynamicLevel = zap.NewAtomicLevelAt(parseLevel(levelStr))
		cores := buildCoresWithFile(asJSON, logFile, config)

		if len(cores) == 0 {
			initErr = fmt.Errorf("failed to create any logger cores")
			return
		}

		zapLogger := zap.New(
			zapcore.NewTee(cores...),
			zap.AddCaller(),
			zap.AddCallerSkip(1), // Skip 1 уровень для правильного отображения caller
		)

		globalLogger = &logger{
			zapLogger: zapLogger,
		}
	})

	if initErr != nil {
		return initErr
	}

	if globalLogger == nil {
		return fmt.Errorf("logger init failed")
	}

	return nil
}

func buildCores(asJSON bool, config *LoggerConfig) []zapcore.Core {
	return buildCoresWithFile(asJSON, "", config)
}

func buildCoresWithFile(asJSON bool, logFile string, config *LoggerConfig) []zapcore.Core {
	var cores []zapcore.Core

	if logFile != "" {
		// Пишем в файл
		if fileCore := createFileCore(asJSON, logFile); fileCore != nil {
			cores = append(cores, fileCore)
		}
	} else {
		// Пишем в stdout
		cores = append(cores, createStdoutCore(asJSON))
	}

	if config != nil && config.OTLPEndpoint != "" {
		if otlpCore := createOTLPCore(config); otlpCore != nil {
			cores = append(cores, otlpCore)
		}
	}

	return cores
}

func createFileCore(asJSON bool, logFile string) zapcore.Core {
	// Создаем директорию если не существует
	dir := filepath.Dir(logFile)
	if dir != "" && dir != "." {
		_ = os.MkdirAll(dir, 0755)
	}

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Failed to open log file %s: %v\n", logFile, err)
		return createStdoutCore(asJSON) // Fallback to stdout
	}

	encoderConfig := buildProductionEncoderConfig()
	var encoder zapcore.Encoder
	if asJSON {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	return zapcore.NewCore(encoder, zapcore.AddSync(file), dynamicLevel)
}

func createOTLPCore(config *LoggerConfig) *SimpleOTLPCore {
	otlpLogger, err := createOTLPLogger(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Failed to initialize OTLP logger: %v\n", err)
		return nil
	}

	return NewSimpleOTLPCore(otlpLogger, dynamicLevel)
}

func createOTLPLogger(config *LoggerConfig) (otelLog.Logger, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	exporter, err := createOTLPExporter(ctx, config.OTLPEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	rs, err := createResource(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	provider := otelLogSdk.NewLoggerProvider(
		otelLogSdk.WithResource(rs),
		otelLogSdk.WithProcessor(otelLogSdk.NewBatchProcessor(exporter)),
	)

	otelProvider = provider

	return provider.Logger(config.ServiceName), nil
}

func createResource(ctx context.Context, config *LoggerConfig) (*resource.Resource, error) {
	return resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(config.ServiceName),
			attribute.String("deployment.environment", config.Environment),
		),
		resource.WithHost(),
		resource.WithOS(),
		resource.WithProcess(),
	)
}

func createOTLPExporter(ctx context.Context, endpoint string) (*otlploggrpc.Exporter, error) {
	return otlploggrpc.New(ctx,
		otlploggrpc.WithEndpoint(endpoint),
		otlploggrpc.WithInsecure(), // [PRODUCTION] Включить TLS
		otlploggrpc.WithTimeout(10*time.Second),
	)
}

func createStdoutCore(asJSON bool) zapcore.Core {
	config := buildProductionEncoderConfig()
	var encoder zapcore.Encoder
	if asJSON {
		encoder = zapcore.NewJSONEncoder(config)
	} else {
		encoder = zapcore.NewConsoleEncoder(config)
	}

	return zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), dynamicLevel)
}

func parseLevel(levelStr string) zapcore.Level {
	switch strings.ToLower(levelStr) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func buildProductionEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "timestamp",
		NameKey:        "name",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
}

func SetLevel(levelStr string) {
	if dynamicLevel == (zap.AtomicLevel{}) {
		return
	}
	dynamicLevel.SetLevel(parseLevel(levelStr))
}

// Logger возвращает глобальный экземпляр логгера
func Logger() *logger {
	return globalLogger
}

func SetNopLogger() {
	globalLogger = &logger{
		zapLogger: zap.NewNop(),
	}
}

func Shutdown(ctx context.Context) error {
	// Сначала синхронизируем zap logger
	if globalLogger != nil {
		if err := globalLogger.zapLogger.Sync(); err != nil {
			// Игнорируем ошибки sync для stdout/stderr
			fmt.Fprintf(os.Stderr, "Warning: logger sync failed: %v\n", err)
		}
	}

	// Затем shutdown OpenTelemetry provider
	if otelProvider != nil {
		shutdownCtx, cancel := context.WithTimeout(ctx, shutdownTimeout)
		defer cancel()

		if err := otelProvider.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("failed to shutdown OTLP logger provider: %w", err)
		}
	}

	return nil
}

func Sync() error {
	if globalLogger != nil {
		return globalLogger.zapLogger.Sync()
	}
	return nil
}

func With(fields ...zap.Field) *logger {
	if globalLogger == nil {
		return &logger{zapLogger: zap.NewNop()}
	}

	return &logger{zapLogger: globalLogger.zapLogger.With(fields...)}
}

func WithContext(ctx context.Context) *logger {
	if globalLogger == nil {
		return &logger{zapLogger: zap.NewNop()}
	}

	return &logger{
		zapLogger: globalLogger.zapLogger.With(fieldsFromContext(ctx)...),
	}
}

// ==================== Глобальные функции логирования ====================

func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	if globalLogger == nil {
		return
	}
	globalLogger.Debug(ctx, msg, fields...)
}

func Info(ctx context.Context, msg string, fields ...zap.Field) {
	if globalLogger == nil {
		return
	}
	globalLogger.Info(ctx, msg, fields...)
}

func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	if globalLogger == nil {
		return
	}
	globalLogger.Warn(ctx, msg, fields...)
}

func Error(ctx context.Context, msg string, fields ...zap.Field) {
	if globalLogger == nil {
		return
	}
	globalLogger.Error(ctx, msg, fields...)
}

func Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	if globalLogger == nil {
		return
	}
	globalLogger.Fatal(ctx, msg, fields...)
}

func (l *logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Debug(msg, allFields...)
}

func (l *logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Info(msg, allFields...)
}

func (l *logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Warn(msg, allFields...)
}

func (l *logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Error(msg, allFields...)
}

func (l *logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	allFields := append(fieldsFromContext(ctx), fields...)
	l.zapLogger.Fatal(msg, allFields...)
}

func fieldsFromContext(ctx context.Context) []zap.Field {
	fields := make([]zap.Field, 0, 3)

	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		traceID := span.SpanContext().TraceID().String()
		if traceID != "" {
			fields = append(fields, zap.String("trace_id", traceID))
		}

		spanID := span.SpanContext().SpanID().String()
		if spanID != "" {
			fields = append(fields, zap.String("span_id", spanID))
		}
	}

	// Также поддерживаем старый способ через context.Value для обратной совместимости
	if traceID, ok := ctx.Value(traceIDKey).(string); ok && traceID != "" {
		// Добавляем только если не было добавлено из span
		found := false
		for _, f := range fields {
			if f.Key == "trace_id" {
				found = true
				break
			}
		}
		if !found {
			fields = append(fields, zap.String("trace_id", traceID))
		}
	}

	// Извлекаем userID из контекста
	if userID, ok := ctx.Value(userIDKey).(string); ok && userID != "" {
		fields = append(fields, zap.String("user_id", userID))
	}

	return fields
}
