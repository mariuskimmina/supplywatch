package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
)

type Config struct {
	Warehouse struct {
		ListenIP string `yaml:"listenIP" envconfig:"SW_LISTENIP"`
		UDPPort  int    `yaml:"udpPort" envconfig:"SW_UDP_PORT"`
		TCPPort  int    `yaml:"tcpPort" envconfig:"SW_TCP_PORT"`
	} `yaml:"warehouse"`
	SensorWarehouse struct {
		UDPPort int `yaml:"udpPort" envconfig:"SW_UDP_PORT"`
	} `yaml:"sensorWarehouse"`
}

// LoadConfig first gets all values from the config.yml file
// if the environment variables are set, they will overrider the config
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config.yml")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	err = envconfig.Process("", &config)
	return
}