package logger

// LogLevel represents the severity level of a log entry
type LogLevel string

const (
	LevelDebug LogLevel = "debug"
	LevelInfo  LogLevel = "info"
	LevelWarn  LogLevel = "warn"
	LevelError LogLevel = "error"
)

// LogContextData holds the context information for logging
type LogContextData struct {
	Tags      map[string]bool
	Category  string
	Metadata  map[string]string
	SessionID string
}

// LogOutput represents the structured log output
type LogOutput struct {
	Level     LogLevel               `json:"level"`
	Message   interface{}            `json:"message,omitempty"`
	SessionID string                 `json:"sessionId,omitempty"`
	Details   map[string]interface{} `json:"details"`
}

// Logger defines the logging interface
type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
}
