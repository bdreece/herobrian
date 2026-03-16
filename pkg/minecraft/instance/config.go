package instance

type Config struct {
	Server  string `mapstructure:"server"`
	Flavor  string `mapstructure:"flavor"`
	Version string `mapstructure:"version"`
}
