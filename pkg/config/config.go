package config

import "github.com/spf13/viper"

type WarehouseConfig struct {
    ListenIP    string `mapstructure:"listenIP"`
    UdpPort     int `mapstructure:"udpPort"`
    TcpPort     int `mapstructure:"tcpPort"`
}

func LoadWarehouseConfig(path string) (config WarehouseConfig, err error) {
    viper.AddConfigPath(path)
    viper.SetConfigName("warehouse.yml")
    viper.SetConfigType("yaml")

    viper.AutomaticEnv()

    err = viper.ReadInConfig()
    if err != nil {
        return
    }

    err = viper.Unmarshal(&config)
    return
}
