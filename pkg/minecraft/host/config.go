package host

import "github.com/bdreece/herobrian/pkg/minecraft/instance"

type Config struct {
	ID        string                     `mapstructure:"id"`
	Instances map[string]instance.Config `mapstructure:"instances"`
}
