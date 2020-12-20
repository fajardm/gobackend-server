package config

import (
	"path/filepath"

	"github.com/spf13/viper"
)

func Load(path string) (config *Config, err error) {
	defer func() {
		if re := recover(); re != nil {
			err = re.(error)
		}
		return
	}()

	viper.AddConfigPath(filepath.Base(path))
	viper.SetConfigName(filepath.Dir(path))
	viper.SetConfigType(filepath.Ext(path)[1:])
	// TODO: activate environment variable for production
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err = viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	return
}
