package config

import (
	"AuthService/pkg/secure"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	System        SystemConfig
	Server        ServerConfig   `yaml:"Server"`
	Postgres      PostgresConfig `yaml:"Postgres"`
	LoggerService LoggerServiceConfig
	Logger        Logger
	Redis         RedisConfig
	EmailService  EmailServiceConfig
	Kyc           Kyc
}

type RedisConfig struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
}

type SystemConfig struct {
	MaxGoRoutines int64
}

//type ServerConfig struct {
//	AppVersion                  string `json:"appVersion"`
//	Host                        string `json:"host" validate:"required"`
//	Port                        string `json:"port" validate:"required"`
//	ShowUnknownErrorsInResponse bool   `json:"showUnknownErrorsInResponse"`
//}

type ServerConfig struct {
	AppVersion                  string `yaml:"AppVersion"`
	Host                        string `yaml:"Host" validate:"required"`
	Port                        string `yaml:"Port" validate:"required"`
	ShowUnknownErrorsInResponse bool   `yaml:"ShowUnknownErrorsInResponse"`
}

//type PostgresConfig struct {
//	Host     string `json:"host"`
//	Port     string `json:"port"`
//	User     string `json:"user"`
//	Password string `json:"-"`
//	DBName   string `json:"DBName"`
//	SSLMode  string `json:"sslMode"`
//	PgDriver string `json:"pgDriver"`
//}

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbName"`
	SSLMode  string `yaml:"sslMode"`
	PgDriver string `yaml:"pgDriver"`
}

type LoggerServiceConfig struct {
	Url        string `json:"url"`
	ServiceId  int64  `json:"serviceId"`
	DevPublic  string `json:"dev_public"`
	DevPrivate string `json:"dev_private"`
}

type Kyc struct {
	Username string `json:"username"`
	Pass     string `json:"pass"`
	Url      string `json:"url"`
}

type EmailServiceConfig struct {
	Url        string `json:"url"`
	ServiceId  int64  `json:"serviceId"`
	DevPublic  string `json:"dev_public"`
	DevPrivate string `json:"dev_private"`
}

type Logger struct {
	Level          string   `json:"level"`
	SkipFrameCount int      `json:"skipFrameCount"`
	InFile         bool     `json:"inFile"`
	FilePath       string   `json:"filePath"`
	InTG           bool     `json:"inTg"`
	ChatID         int64    `json:"chatID"`
	TGToken        string   `json:"-"`
	AlertUsers     []string `json:"alertUsers"`
}

func NewConfig() (*Config, error) {
	var c Config

	configPath := filepath.Join("C:\\Users\\user\\GolandProjects\\sso-service-refresh-token\\config\\", "cfg.yml")
	f, err := os.Open(configPath)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&c)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	return &c, nil
}

func LoadConfig() (*viper.Viper, error) {

	viperInstance := viper.New()

	if _, ok := os.LookupEnv("LOCAL"); ok {
		viperInstance.AddConfigPath("config")
	} else {
		viperInstance.AddConfigPath("./config")
	}
	viperInstance.SetConfigName("config")
	viperInstance.SetConfigType("yml")

	err := viperInstance.ReadInConfig()
	if err != nil {
		return nil, err
	}
	return viperInstance, nil
}

func ParseConfig(v *viper.Viper) (*Config, error) {

	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		log.Fatalf("unable to decode config into struct, %v", err)
		return nil, err
	}
	return &c, nil
}

func DecryptConfig(cfg *Config, s *secure.Shield) {

	cfg.Postgres.Host = s.DecryptMessage(cfg.Postgres.Host)
	cfg.Postgres.Port = s.DecryptMessage(cfg.Postgres.Port)
	cfg.Postgres.User = s.DecryptMessage(cfg.Postgres.User)
	cfg.Postgres.Password = s.DecryptMessage(cfg.Postgres.Password)
	cfg.Postgres.DBName = s.DecryptMessage(cfg.Postgres.DBName)

	cfg.LoggerService.Url = s.DecryptMessage(cfg.LoggerService.Url)
	cfg.LoggerService.DevPublic = s.DecryptMessage(cfg.LoggerService.DevPublic)
	cfg.LoggerService.DevPrivate = s.DecryptMessage(cfg.LoggerService.DevPrivate)

	cfg.Redis.Addr = s.DecryptMessage(cfg.Redis.Addr)
	cfg.Redis.Password = s.DecryptMessage(cfg.Redis.Password)
}
