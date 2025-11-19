package config

type (
	ServerConfig struct {
		Host              string `yaml:"host"`
		Port              int    `yaml:"port"`
		EnableOptions     bool   `yaml:"enableOptions"`
		ReadTimeout       int    `yaml:"readTimeout"`
		ReadHeaderTimeout int    `yaml:"readHeaderTimeout"`
		WriteTimeout      int    `yaml:"writeTimeout"`
		IdleTimeout       int    `yaml:"idleTimeout"`
		MaxHeaderBytes    []int  `yaml:"maxHeaderBytes"`
	}
	serverConfigOption func(c *ServerConfig)
)

func initServerConfig(host string, port int) ServerConfig {
	return ServerConfig{
		Host:           host,
		Port:           port,
		EnableOptions:  true,
		MaxHeaderBytes: []int{1, 20}, // default -> 1Mib
	}
}

func (c *ServerConfig) SetupConfig(opts ...serverConfigOption) {
	for _, fn := range opts {
		fn(c)
	}
}
