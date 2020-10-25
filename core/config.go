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
}

type HTTPConfig struct {
	Insecure bool
	Timeout  int
}
