package core

// Struct for config flags
type Config struct {
	HTTP           HTTPConfig
	MaxSeverity    string
	ExportFormats  []string
	Urls           []string
	ExportFilename string
	SeverityFilter string
	PluginFilter   []string
	Threads        int
}

type HTTPConfig struct {
	Insecure bool
	Timeout  int
	UserAgent string
}
