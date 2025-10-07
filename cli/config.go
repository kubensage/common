package gocli

import "flag"

// LogStdAndFileConfig holds configuration options for logging to both
// standard output and a rotating log file.
type LogStdAndFileConfig struct {
	LogLevel      string // Log verbosity level (e.g., "info", "debug", "error")
	LogFile       string // Path to the log file
	LogMaxSize    int    // Maximum size (in MB) before log file is rotated
	LogMaxBackups int    // Maximum number of old log files to retain
	LogMaxAge     int    // Maximum age (in days) to retain old log files
	LogCompress   bool   // Whether to compress old log files
}

// LogStdConfig holds configuration options for logging to standard output only.
type LogStdConfig struct {
	LogLevel string // Log verbosity level (e.g., "info", "debug", "error")
}

// RegisterLogStdAndFileFlags registers command-line flags for configuring
// both standard output and file-based logging with rotation settings.
//
// The provided FlagSet `fs` is used to define flags such as log level, log file path,
// max size, backup count, age, and compression behavior. The `appName` is used
// to generate a default log file path (e.g., /var/log/kubensage/myapp.log).
//
// Registered flags:
//
//	--log-level        string   Log verbosity level (default "info")
//	--log-file         string   Path to log file (default "/var/log/kubensage/<appName>.log")
//	--log-max-size     int      Max log file size in MB before rotation (default 10)
//	--log-max-backups  int      Max number of old log files to retain (default 5)
//	--log-max-age      int      Max age in days to retain old log files (default 30)
//	--log-compress     bool     Whether to compress old log files (default true)
//
// Parameters:
//   - fs       The flag set into which the flags will be registered.
//   - appName  The name of the application (used in the default log file path).
//
// Returns:
//
//	A closure that, when invoked, returns a populated *LogStdAndFileConfig
//	containing the values from the parsed flags.
func RegisterLogStdAndFileFlags(
	fs *flag.FlagSet,
	appName string,
) func() *LogStdAndFileConfig {
	logPath := "/var/log/kubensage/" + appName + ".log"

	logLevel := fs.String("log-level", "info", "Set log level")
	logFile := fs.String("log-file", logPath, "Path to log file")
	logMaxSize := fs.Int("log-max-size", 10, "Maximum log size (MB)")
	logMaxBackups := fs.Int("log-max-backups", 5, "Max backup files")
	logMaxAge := fs.Int("log-max-age", 30, "Max age in days")
	logCompress := fs.Bool("log-compress", true, "Compress logs")

	return func() *LogStdAndFileConfig {
		return &LogStdAndFileConfig{
			LogLevel:      *logLevel,
			LogFile:       *logFile,
			LogMaxSize:    *logMaxSize,
			LogMaxBackups: *logMaxBackups,
			LogMaxAge:     *logMaxAge,
			LogCompress:   *logCompress,
		}
	}
}

// RegisterLogStdFlags registers command-line flags for configuring logging
// to standard output only (without log file rotation).
//
// Registered flags:
//
//	--log-level string   Log verbosity level (default "info")
//
// Parameters:
//   - fs  The flag set into which the flags will be registered.
//
// Returns:
//
//	A closure that, when invoked, returns a populated *LogStdConfig
//	containing the value from the parsed flags.
func RegisterLogStdFlags(
	fs *flag.FlagSet,
) func() *LogStdConfig {
	logLevel := fs.String("log-level", "info", "Set log level")

	return func() *LogStdConfig {
		return &LogStdConfig{
			LogLevel: *logLevel,
		}
	}
}
