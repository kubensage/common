package golog

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"runtime"
	"time"

	"github.com/kubensage/go-common/cli"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LogStartupInfo logs standard metadata at startup, including Go version,
// executable path, and current time. Optionally, any configuration structs passed
// are logged under their type name after sanitization.
//
// Parameters:
//   - logger: the zap.Logger to use for output.
//   - appName: the name of the application (used in the log message).
//   - configs: optional list of configuration structs to log.
//
// The function stringifies fields such as time.Duration for readability
// and excludes unexported fields.
func LogStartupInfo(
	logger *zap.Logger,
	appName string,
	configs ...any,
) {
	exePath, err := os.Executable()
	if err != nil {
		logger.Warn("Could not determine executable path", zap.Error(err))
		exePath = "unknown"
	}

	fields := []zap.Field{
		zap.String("go_version", runtime.Version()),
		zap.String("executable", exePath),
		zap.Time("start_time", time.Now()),
	}

	// Sanitize and log each config struct under its type name
	for _, cfg := range configs {
		fields = append(fields, zap.Any(getTypeName(cfg), sanitizeConfig(cfg)))
	}

	logger.Info(appName+" started", fields...)
}

// SetupStdLogger creates and returns a zap.Logger that writes logs to standard output.
// The logger format is JSON and the log level is determined by the given configuration.
//
// Parameters:
//   - cfg: the logging configuration (standard output only).
//
// Returns:
//   - *zap.Logger configured for stdout.
//
// Panics if the logger cannot be created.
func SetupStdLogger(
	cfg *gocli.LogStdConfig,
) *zap.Logger {
	logger, err := newStdLogger(cfg.LogLevel)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	return logger
}

// SetupStdAndFileLogger creates and returns a zap.Logger that writes logs to both
// standard output and a rotating file. File rotation settings are derived from the config.
//
// Parameters:
//   - cfg: the logging configuration, including file path and rotation policy.
//
// Returns:
//   - *zap.Logger configured for dual output.
//
// Panics if the logger cannot be created.
func SetupStdAndFileLogger(
	cfg *gocli.LogStdAndFileConfig,
) *zap.Logger {
	logger, err := newStdAndFileLogger(
		&cfg.LogLevel,
		&cfg.LogFile,
		&cfg.LogMaxSize,
		&cfg.LogMaxBackups,
		&cfg.LogMaxAge,
		&cfg.LogCompress,
	)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	return logger
}

// newStdAndFileLogger builds a zap.Logger that writes to both stdout and a file with log rotation.
//
// Parameters:
//   - logLevel: log verbosity level (e.g., "info", "debug").
//   - file: path to the log file.
//   - size: max size in MB before log rotation.
//   - backups: number of old logs to retain.
//   - age: max age in days for old logs.
//   - compress: whether to compress old logs.
//
// Returns:
//   - *zap.Logger configured with dual cores (file + stdout).
//   - error if log level is invalid.
func newStdAndFileLogger(
	logLevel *string,
	file *string,
	size *int,
	backups *int,
	age *int,
	compress *bool,
) (*zap.Logger, error) {
	level := zapcore.InfoLevel
	if err := (&level).UnmarshalText([]byte(*logLevel)); err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderCfg)

	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   *file,
		MaxSize:    *size,
		MaxBackups: *backups,
		MaxAge:     *age,
		Compress:   *compress,
	})

	stdoutWriter := zapcore.AddSync(os.Stdout)

	fileCore := zapcore.NewCore(encoder, fileWriter, level)
	stdoutCore := zapcore.NewCore(encoder, stdoutWriter, level)

	core := zapcore.NewTee(fileCore, stdoutCore)

	return zap.New(core), nil
}

// newStdLogger builds a zap.Logger that logs exclusively to stdout using the provided log level.
//
// Parameters:
//   - logLevel: string representation of the desired log level.
//
// Returns:
//   - *zap.Logger for stdout.
//   - error if the log level is invalid.
func newStdLogger(
	logLevel string,
) (*zap.Logger, error) {
	level := zapcore.InfoLevel
	if err := (&level).UnmarshalText([]byte(logLevel)); err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderCfg)

	stdoutWriter := zapcore.AddSync(os.Stdout)
	stdoutCore := zapcore.NewCore(encoder, stdoutWriter, level)

	return zap.New(stdoutCore), nil
}

// sanitizeConfig converts a struct to a map of field names to values,
// intended for structured logging. Fields of type time.Duration are
// converted to their string representation for readability.
//
// Parameters:
//   - cfg: any struct (pointer or value) to sanitize.
//
// Returns:
//   - map[string]any: a sanitized representation of the struct.
func sanitizeConfig(
	cfg any,
) any {
	val := reflect.ValueOf(cfg)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	out := make(map[string]any)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		var value any

		// Convert time.Duration fields to their string representation (e.g., "5s")
		switch field.Kind() {
		case reflect.Int64:
			if field.Type().PkgPath() == "time" && field.Type().Name() == "Duration" {
				value = field.Interface().(time.Duration).String()
			} else {
				value = field.Interface()
			}
		default:
			value = field.Interface()
		}

		out[fieldType.Name] = value
	}

	return out
}

// getTypeName returns the name of the type of the given value,
// automatically dereferencing pointers.
//
// Parameters:
//   - v: any value.
//
// Returns:
//   - string: the underlying type name.
func getTypeName(
	v any,
) string {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}
