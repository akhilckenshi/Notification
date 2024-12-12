package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var Config Configuration

type Configuration struct {
	Environment            string `mapstructure:"environment"`
	Logger                 LoggerConfig
	Database               DatabaseConfig
	WhatsApp               WhatsAppConfig
	App                    AppConfig
	Email                  EmailConfig
	DBURI                  string `mapstructure:"DBURI"`
	DBName                 string `mapstructure:"DBNAME"`
	DBConnCount            int    `mapstructure:"DBCONNCNT"`
	DBHost                 string `mapstructure:"DB_HOST"`
	AppPort                string `mapstructure:"PORT"`
	KafkaPort              string `mapstructure:"KAFKA_PORT"`
	KafkaTopic             string `mapstructure:"KAFKA_TOPIC"`
	AppEmailID             string `mapstructure:"APP_EMILID"`
	AppEmailPassword       string `mapstructure:"APP_EMAIL_PWD"`
	SMTPHost               string `mapstructure:"SMTP_HOST"`
	SMTPPort               int    `mapstructure:"SMTP_PORT"`
	WhatsProviderURL       string `mapstructure:"WHATS_PROVIDER_URL"`
	WhatsAppProviderKey    string `mapstructure:"WHATSAPP_PROVIDER_KEY"`
	WhatsAppProviderSecret string `mapstructure:"WHATSAPP_PROVIDER_SECRET"`
	WhatsAppFromNumber     string `mapstructure:"WHATSAPP_FROM_NUMBER"`
}

type LoggerConfig struct {
	FileName     string `mapstructure:"fileName"`
	FileSize     int    `mapstructure:"fileSize"`
	MaxLogFile   int    `mapstructure:"maxLogFile"`
	MaxRetention int    `mapstructure:"maxRetention"`
	CompressLog  bool   `mapstructure:"compressLog"`
	Level        string `mapstructure:"level"`
}

type DatabaseConfig struct {
	DBTimeout int         `mapstructure:"dbTimeout"`
	DBType    string      `mapstructure:"dbType"`
	MongoDb   MongoConfig `mapstructure:"mongo"`
}

type MongoConfig struct {
	DBUri    string `mapstructure:"dbUri"`
	DBName   string `mapstructure:"dbName"`
	MaxLimit int    `mapstructure:"maxLimit"`
}

type WhatsAppConfig struct {
	Provider     string `mapstructure:"provider"`
	AuthRequired bool   `mapstructure:"authRequired"`
	Key          string `mapstructure:"key"`
	Secret       string `mapstructure:"secret"`
	Number       string `mapstructure:"number"`
}

type AppConfig struct {
	WithSSL bool `mapstructure:"withssl"`
}

type EmailConfig struct {
	Id       string `mapstructure:"id"`
	Username string `mapstructure:"username"`
	Pwd      string `mapstructure:"password"`
	SmtpHost string `mapstructure:"smtpHost"`
	SmtpPort int    `mapstructure:"smtpPort"`
}

func InitConfig() (Configuration, error) {
	env := os.Getenv("APP_ENV")
	var configDir, envDir string

	if env == "production" {
		configDir = "/etc/config"
		envDir = "/etc/env"
	} else {
		currentDir, err := os.Getwd()
		if err != nil {
			return Config, err
		}
		configDir = filepath.Join(currentDir, "../../config")
		if _, err := os.Stat(configDir); os.IsNotExist(err) {
			configDir = "./"
		}
		envDir = filepath.Join(currentDir, "../../env")
		if _, err := os.Stat(envDir); os.IsNotExist(err) {
			envDir = "./"
		}
	}

	// Load `config.yaml`
	viper.AddConfigPath(configDir)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		return Config, err
	}

	// Load `.env`
	viper.AddConfigPath(envDir)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	if err := viper.MergeInConfig(); err != nil {
		return Config, err
	}

	// Bind sensitive environment variables
	envVars := []string{
		"DBURI", "DBNAME", "PORT", "KAFKA_PORT", "KAFKA_TOPIC",
		"APP_EMILID", "APP_USERNAME", "APP_PWD", "SMTP_HOST", "SMTP_PORT",
		"WHATS_PROVIDER_URL", "WHATSAPP_PROVIDER_KEY", "WHATSAPP_PROVIDER_SECRET", "WHATSAPP_FROM_NUMBER",
	}
	for _, envVar := range envVars {
		if err := viper.BindEnv(envVar); err != nil {
			return Config, err
		}
	}

	// Unmarshal the configuration
	if err := viper.Unmarshal(&Config); err != nil {
		return Config, err
	}
	fmt.Println("config-----------", Config)
	return Config, nil
}
