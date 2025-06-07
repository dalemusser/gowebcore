package config

import "github.com/spf13/pflag"

type Option func(*loader)

func WithEnvPrefix(p string) Option     { return func(l *loader) { l.envPrefix = p } }
func WithConfigFile(path string) Option { return func(l *loader) { l.cfgFile = path } }
func WithFlagSet(fs *pflag.FlagSet) Option {
	return func(l *loader) { l.flagset = fs }
}
