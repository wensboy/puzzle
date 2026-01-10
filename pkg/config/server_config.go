package config

type (
	ServerConfig struct {
		Addr              string `yaml:"addr"`
		Host              string `yaml:"host"`
		Port              int    `yaml:"port"`
		EnableOptions     bool   `yaml:"enableOptions"`
		ReadTimeout       int    `yaml:"readTimeout"`
		ReadHeaderTimeout int    `yaml:"readHeaderTimeout"`
		WriteTimeout      int    `yaml:"writeTimeout"`
		IdleTimeout       int    `yaml:"idleTimeout"`
		MaxHeaderBytes    []int  `yaml:"maxHeaderBytes"`
	}
)

func initServerConfig() ServerConfig {
	return ServerConfig{
		EnableOptions:  true,
		MaxHeaderBytes: []int{1, 20}, // default -> 1Mib
	}
}
