package config

import (
	"github.com/spf13/viper"
)

// type Config struct {
// Warehouse struct {
// ListenIP   string `yaml:"listenIP" envconfig:"SW_LISTENIP"`
// UDPPort    int    `yaml:"udpPort" envconfig:"SW_UDP_PORT"`
// TCPPort    int    `yaml:"tcpPort" envconfig:"SW_TCP_PORT"`
// LogFileDir string `yaml:"logFileBase" envconfig:"SW_LOG_FILE_DIR"`
// } `yaml:"warehouse"`
// SensorWarehouse struct {
// UDPPort      int `yaml:"udpPort" envconfig:"SW_UDP_PORT"`
// Delay        int `yaml:"delay" envconfig:"SW_DELAY"`
// NumOfPackets int `yaml:"numOfPackets" envconfig:"SW_NUMOFPACKETS"`
// } `yaml:"sensorWarehouse"`
// }

type WarehouseConfig struct {
	ListenIP   string `mapstructure:"SW_LISTEN_IP"`
	UDPPort    string `mapstructure:"SW_UDP_PORT"`
	TCPPort    string `mapstructure:"SW_TCP_PORT"`
	GRPCPort   string `mapstructure:"SW_GRPC_PORT"`
	DBUser     string `mapstructure:"SW_DATABASE_USER"`
	DBDatabase string `mapstructure:"SW_DATABASE_DB"`
	DBPassword string `mapstructure:"SW_DATABASE_PASSWORD"`
	DBPort     int    `mapstructure:"SW_DATABASE_PORT"`
	LogFileDir string `mapstructure:"SW_LOG_FILE_DIR"`
}

type SensorConfig struct {
	Delay           int
	NumberOfPackets int
}

type SupplywatchConfig struct {
	NumOfWarehouses   int `mapstructure:"SW_NUMBER_OF_WAREHOUSES"`
}

// LoadConfig first gets all values from the config.yml file
// if the environment variables are set, they will overrider the config
func LoadWarehouseConfig(path string) (config WarehouseConfig, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("warehouse")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

func LoadSensorConfig(path string) (config SensorConfig, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("sensor")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

func LoadSupplywatchConfig(path string) (config SupplywatchConfig, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("supplywatch")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
