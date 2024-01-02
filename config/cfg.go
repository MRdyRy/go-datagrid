package config

import "github.com/spf13/viper"

type ConfigureEnv struct {
	DatagridUrl      string `mapstructure:"URL"`
	DatagridPort     string `mapstructure:"PORT"`
	DatagridUser     string `mapstructure:"USER"`
	DatagridPass     string `mapstructure:"PASS"`
	DatagridProtocol string `mapstructure:"PROTOCOL"`
}

func LoadConfig(path string) (config ConfigureEnv, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("dev")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	return
}
