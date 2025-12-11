package env

import "os"

type LoggerConfig struct {
	level   string
	asJson  bool
	logFile string
}

func (l *LoggerConfig) Level() string   { return l.level }
func (l *LoggerConfig) AsJson() bool    { return l.asJson }
func (l *LoggerConfig) LogFile() string { return l.logFile }

func LoadLoggerConfig() *LoggerConfig {
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "info"
	}
	return &LoggerConfig{
		level:   level,
		asJson:  os.Getenv("LOG_AS_JSON") == "true",
		logFile: os.Getenv("LOG_FILE"),
	}
}
