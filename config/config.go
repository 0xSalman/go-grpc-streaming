package config

import (
  "os/user"
  "strings"
  
  "github.com/kelseyhightower/envconfig"
)

type Config struct {
  Port             string `envconfig:"PORT" default:"7000"`
  PrivateKey       string `envconfig:"PRIVATE_KEY" default:"~/.ssh/maxnumber_rsa_private.pem"`
  PublicKey        string `envconfig:"PUBLIC_KEY" default:"~/.ssh/maxnumber_rsa_public.pem"`
  NumbersToSend    int    `envconfig:"TOTAL_NUMBERS" default:"15"`
  NumberMultiplier int    `envconfig:"NUMBER_MULTIPLIER" default:"100"`
}

func LoadConfig() (*Config, error) {
  conf := &Config{}
  err := envconfig.Process("grpc", conf)
  if err != nil {
    return nil, err
  }
  return conf, nil
}

func UserHomeDir() (string, error) {
  usr, err := user.Current()
  if err != nil {
    return "", err
  }
  return usr.HomeDir, nil
}

func AbsolutePath(path string) (string, error) {
  if strings.HasPrefix(path, "~") {
    userHome, err := UserHomeDir()
    absPath := userHome + strings.TrimPrefix(path, "~")
    return absPath, err
  }
  return path, nil
}
