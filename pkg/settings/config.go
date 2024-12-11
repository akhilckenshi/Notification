/*
config.go
Author: Bipin Kumar Ojha
Description: This file handles the configuration setup for the application.
It reads configuration from a config.yaml file and stores it in structured types for easier access.
*/

package cfg

import (
	"log"
	"os"
	"smartsme-notificationservice/pkg/utils"

	"github.com/spf13/viper"
)

var Config Configuration

// Configuration struct holds all the configurations from the config.yaml file
type Configuration struct {
	Server         ServerConfig   `mapstructure:"server"`   // Server configuration
	Logger         LoggerConfig   `mapstructure:"logger"`   // Logger configuration
	Database       DatabaseConfig `mapstructure:"database"` // Database configuration
	Kafka          KafkaConfig    `mapstructure:"kafka"`    // Kafka configuration
	EMail          Email          `mapstructure:"email"`    // Email configuration
	WhatsappConfig WhatsAppConfig `mapstructure:"whatsppp"` // Whatsapp config
}

type Email struct {
	Id       string `mapstructure:"id"`       // EMail ID
	Username string `mapstructure:"username"` //Username
	Pwd      string `mapstructure:"password"` // Password
	SmtpHost string `mapstructure:"smtpHost"` // smtp host
	SmtpPort int    `mapstructure:"smtpPort"` // smtp port
}

type WhatsAppConfig struct {
	Provider     string `mapstructure:"provider"`
	AuthRequired bool   `mapstructure:"authRequired"`
	Key          string `mapstructure:"key"`
	Secret       string `mapstructure:"secret"`
	Number       string `mapstructure:"number"`
}

// ServerConfig holds server-specific configurations
type ServerConfig struct {
	Host string `mapstructure:"host"` // Server host
	Port int    `mapstructure:"port"` // Server port
}

// LoggerConfig holds logger-specific configurations
type LoggerConfig struct {
	FileName     string `mapstructure:"fileName"`     // Log file name
	FileSize     int    `mapstructure:"fileSize"`     // Log file size in MB
	MaxLogFile   int    `mapstructure:"maxLogFile"`   // Maximum number of log files
	MaxRetention int    `mapstructure:"maxRetention"` // Maximum retention period for logs in days
	CompressLog  bool   `mapstructure:"compressLog"`  // Whether to compress log files
	Level        string `mapstructure:"level"`
}

// DatabaseConfig holds database-specific configurations
type DatabaseConfig struct {
	DBTimeout int         `mapstructure:"dbTimeout"` // Database request timeout
	DBType    string      `mapstructure:"dbType"`    // Database Type
	MongoDb   MongoConfig `mapstructure:"mongo"`     // Mongo specific configuration
}

// MongoConfig holds Mongo-specific configurations
type MongoConfig struct {
	DBUri    string `mapstructure:"dbUri"`    // Database URI
	DBName   string `mapstructure:"dbName"`   // Database Name
	MaxLimit int    `mapstructure:"maxLimit"` // No of concurrent connection
}

// KafkaConfig holds Kafka-specific configurations
type KafkaConfig struct {
	KafkaPort  string `mapstructure:"KAFKA_PORT"`  // Kafka Port
	KafkaTopic string `mapstructure:"KAFKA_TOPIC"` // Kafka Topic
}

// InitConfig initializes the application configuration by reading the config.yaml file
func InitConfig() {
	// Get the current working directory
	currentWorkDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err.Error()) // Exit if there's an error getting the working directory
	}

	// Add configuration paths
	viper.AddConfigPath(currentWorkDirectory + "/../../config") // Relative path for local development
	viper.AddConfigPath("/etc/viper/config")                    // Path for production environment
	viper.AddConfigPath(".")                                    // Current directory

	// Set the configuration file name and type
	viper.SetConfigName("config") // Configuration file name (without extension)
	viper.SetConfigType("yaml")   // Configuration file type

	// Read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		log.Println("Could not read config file, using default values!") // Log a warning if the config file isn't found
	}

	// Unmarshal the configuration into the Config struct
	if err := viper.Unmarshal(&Config); err != nil {
		log.Fatalf("Could not unmarshal config: %v", err) // Exit if there's an error unmarshalling the config
	}

	utils.LoadEnv(currentWorkDirectory)

	viper.AutomaticEnv()

	Config.Database.MongoDb.DBUri = viper.GetString("DBURI")
	Config.Database.MongoDb.DBName = viper.GetString("DBNAME")
	Config.Server.Host = viper.GetString("HOST")
	Config.Server.Port = viper.GetInt("PORT")
	Config.EMail.Id = viper.GetString("APP_EMILID")
	Config.EMail.Username = viper.GetString("APP_USERNAME")
	Config.EMail.Pwd = viper.GetString("APP_PWD")
	Config.EMail.SmtpHost = viper.GetString("SMTP_HOST")
	Config.EMail.SmtpPort = viper.GetInt("SMTP_PORT")
	Config.WhatsappConfig.Key = viper.GetString("WHATSAPP_PROVIDER_KEY")
	Config.WhatsappConfig.Secret = viper.GetString("WHATSAPP_PROVIDER_SECRET")
	Config.WhatsappConfig.Number = viper.GetString("WHATSAPP_FROM_NUMBER")
	Config.Kafka.KafkaPort = viper.GetString("KAFKA_PORT")
	Config.Kafka.KafkaTopic = viper.GetString("KAFKA_TOPIC")

}
