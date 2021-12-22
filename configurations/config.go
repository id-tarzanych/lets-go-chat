package configurations

import (
	"errors"
	"github.com/ilyakaznacheev/cleanenv"
)

type Configuration struct {
	Database Database
	Server   Server
}

type Database struct {
	Type     string `yaml:"type" env:"LETS_GO_CHAT_DATABASE__TYPE"`
	Host     string `yaml:"host" env:"LETS_GO_CHAT_DATABASE__HOST"`
	Port     int    `yaml:"port" env:"LETS_GO_CHAT_DATABASE__PORT"`
	Protocol string `yaml:"protocol" env:"LETS_GO_CHAT_DATABASE_PROTOCOL"`
	User     string `yaml:"user" env:"LETS_GO_CHAT_DATABASE__USER"`
	Password string `yaml:"password" env:"LETS_GO_CHAT_DATABASE__PASSWORD"`
	Database string `yaml:"database" env:"LETS_GO_CHAT_DATABASE__DATABASE"`
	Ssl      bool   `yaml:"ssl" env:"LETS_GO_CHAT_DATABASE__SSL"`
	SslCert  string `yaml:"sslCert" env:"LETS_GO_CHAT_DATABASE__SSLCA" env-default:"db/cert/ca-certificate.crt"`
}

type Server struct {
	Port int `yaml:"port"  env:"LETS_GO_CHAT_SERVER__PORT" env-default:"8080"`
}

func New() (*Configuration, error) {
	cfg := Configuration{Database{}, Server{}}

	err := cleanenv.ReadConfig("config.yml", &cfg)
	if err != nil {
		return nil, errors.New("could not read configs")
	}

	return &cfg, nil
}
