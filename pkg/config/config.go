package config

import "github.com/spf13/viper"

type Config struct {
    Warehouse struct {
        ListenIP    string `yaml:"listenIP"`
        UDPPort     int `yaml:"udpPort"`
        TCPPort     int `yaml:"tcpPort"`
    }
    SensorWarehouse struct {
        UDPPort     int `yaml:"udpPort"`
    }
}

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
    return
}
