package config

import "github.com/spf13/viper"

type Config struct {
	XToken        string `mapstructure:"XTOKEN"`
	TonURL        string `mapstructure:"TON_URL"`
	MainAddress   string `mapstructure:"MAIN_ADDRESS"`
	TelegramToken string `mapstructure:"TELEGRAM_TOKEN"`

	ConfigTON string `mapstructure:"CONFIG_TON"`

	NotificationDest string `mapstructure:"NOTIFICATION_DESTINATION"`
}

func NewConfig() (*Config, error) {

	var (
		conf *Config
	)

	err := viperSetup()
	if err != nil {
		return &Config{}, err

	}

	err = viper.ReadInConfig()
	if err != nil {
		return &Config{}, err

	}
	err = viper.Unmarshal(&conf)
	if err != nil {
		return &Config{}, err
	}
	return conf, nil

}

func viperSetup() error {
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	return nil
}
