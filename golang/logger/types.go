package logger

type LogLevel string

const (
	LevelDebug LogLevel = "debug"
	LevelInfo  LogLevel = "info"
	LevelWarn  LogLevel = "warn"
	LevelError LogLevel = "error"
)

type LogContextData struct {
	Tags      map[string]bool
	Category  string
	Metadata  map[string]string
	SessionID string
}

type LogOutput struct {
	Level     LogLevel               `json:"level"`
	Message   interface{}            `json:"message,omitempty"`
	SessionID string                 `json:"sessionId,omitempty"`
	Details   map[string]interface{} `json:"details"`
}
