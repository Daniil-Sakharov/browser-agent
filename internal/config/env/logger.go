package env

import "os"

type LoggerConfig struct {
	level   string
	asJson  bool
	logFile string
}

func (l *LoggerConfig) Level() string {
	return l.level
}

func (l *LoggerConfig) AsJson() bool {
	return l.asJson
}

func (l *LoggerConfig) LogFile() string {
	return l.logFile
}

// NewLoggerConfig создает конфигурацию логгера из ENV
func NewLoggerConfig() (*LoggerConfig, error) {
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "info"
	}

	asJson := os.Getenv("LOG_AS_JSON") == "true"
	logFile := os.Getenv("LOG_FILE")

	return &LoggerConfig{
		level:   level,
		asJson:  asJson,
		logFile: logFile,
	}, nil
}
