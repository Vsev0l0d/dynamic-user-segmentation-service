package config

import "github.com/spf13/viper"

type (
	Config struct {
		DB                                `mapstructure:"db"`
		HTTP                              `mapstructure:"http"`
		PeriodForDeletingInactiveSegments `mapstructure:"period_for_deleting_inactive_segments"`
	}

	DB struct {
		Host        string `mapstructure:"host"`
		Port        int    `mapstructure:"port"`
		DbName      string `mapstructure:"db_name"`
		SslMode     string `mapstructure:"sslmode"`
		EnvUser     string `mapstructure:"env_user"`
		EnvPassword string `mapstructure:"env_password"`
		DriverName  string `mapstructure:"driver_name"`
	}

	HTTP struct {
		Addr string `mapstructure:"addr"`
	}

	PeriodForDeletingInactiveSegments struct {
		CronExpression string `mapstructure:"cron_expression"`
	}
)

func NewConfig(configPath string) (*Config, error) {
	conf := &Config{}
	viper.SetConfigFile(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&conf)
	return conf, nil
}
