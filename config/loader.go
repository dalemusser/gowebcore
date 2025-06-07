package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type loader struct {
	envPrefix string
	flagset   *pflag.FlagSet
	cfgFile   string
}

// Load merges flag → env → file values into the struct pointed to by dst.
func Load(dst any, opts ...Option) error {
	l := &loader{}
	for _, opt := range opts {
		opt(l)
	}

	v := viper.New()
	if l.envPrefix != "" {
		v.SetEnvPrefix(l.envPrefix)
	}
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := bindEnvs(v, dst); err != nil {
		return err
	}

	if l.cfgFile != "" {
		v.SetConfigFile(l.cfgFile)
		if err := v.ReadInConfig(); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("read config: %w", err)
		}
	}

	if l.flagset != nil {
		if err := v.BindPFlags(l.flagset); err != nil {
			return err
		}
	}

	if err := v.Unmarshal(dst); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	if v, ok := dst.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("config validation: %w", err)
		}
	}

	return nil
}

// --- helpers -------------------------------------------------------------

// bindEnvs walks the struct and calls v.BindEnv for every mapstructure tag.
func bindEnvs(v *viper.Viper, iface any, path ...string) error {
	val := reflect.ValueOf(iface)
	for val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)

		// recurse into embedded structs
		if f.Anonymous && f.Type.Kind() == reflect.Struct {
			if err := bindEnvs(v, val.Field(i).Interface(), path...); err != nil {
				return err
			}
			continue
		}

		tag := f.Tag.Get("mapstructure")
		if tag == "" || tag == "-" {
			continue
		}
		full := strings.Join(append(path, tag), ".")
		if err := v.BindEnv(full); err != nil {
			return fmt.Errorf("bind env %s: %w", full, err)
		}
	}
	return nil
}
