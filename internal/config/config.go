package config

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/Lanworm/image-previewe/internal/validation"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Logger  LoggerConf
	Server  ServerConf
	Cache   CacheConf
	Storage StorageConf
}

type ServerConf struct {
	HTTP ServerHTTPConf
}

type ServerHTTPConf struct {
	Host     string `validate:"required"`
	Port     int    `validate:"required"`
	Protocol string
	Timeout  time.Duration
}

func (s *ServerHTTPConf) GetFullAddress() string {
	address := net.JoinHostPort(s.Host, strconv.Itoa(s.Port))

	return address
}

type LoggerConf struct {
	Level string `validate:"required,oneof=DEBUG INFO WARNING ERROR"`
}
type CacheConf struct {
	Capacity int `validate:"required,gt=1,lte=99"`
}
type StorageConf struct {
	Path string `validate:"required,dirpath"`
}

func NewConfig(configFile string) (*Config, error) {
	fileData, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	c := &Config{}
	err = yaml.Unmarshal(fileData, c)
	if err != nil {
		log.Fatalf("parse congig file: %v", err)
	}

	err = validation.Validate(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
