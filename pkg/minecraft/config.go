package minecraft

type HostConfig struct {
	ID        string                    `mapstructure:"id"`
	Instances map[string]InstanceConfig `mapstructure:"instances"`
}

type InstanceConfig struct {
	Server  string `mapstructure:"server"`
	Flavor  string `mapstructure:"flavor"`
	Version string `mapstructure:"version"`
}
