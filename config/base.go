package config

// Base holds only the truly cross-cutting knobs shared by all services.
// Each service embeds this and adds its own database, cache, etc.
type Base struct {
	AppName         string   `mapstructure:"app_name"`
	Env             string   `mapstructure:"env"` // dev|staging|prod
	HTTPPort        int      `mapstructure:"http_port"`
	HTTPSPort       int      `mapstructure:"https_port"`
	Domain          string   `mapstructure:"domain"`
	EnableTLS       bool     `mapstructure:"enable_tls"`
	LogLevel        string   `mapstructure:"log_level"` // debug|info|warn|error
	CORSOrigins     []string `mapstructure:"cors_origins"`
	CertFile        string   `mapstructure:"cert_file"`                        // fullchain.pem
	KeyFile         string   `mapstructure:"key_file"`                         // privkey.pem
	MaxBodyBytes    int      `mapstructure:"max_body_bytes" default:"4194304"` // 4 MB
	RateLimitEvents int      `mapstructure:"rate_limit_events" default:"10"`   // tokens per second
	RateLimitBurst  int      `mapstructure:"rate_limit_burst"  default:"20"`   // bucket size
}
