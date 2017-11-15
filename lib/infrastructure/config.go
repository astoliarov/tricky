package infrastructure

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Debug     bool
	ApiPort   string
	ProxyPort string
	CertPath  string
}

func (c *Config) SetDefaults() {
	c.Debug = true
	c.ApiPort = ":5050"
	c.ProxyPort = ":3030"
	c.CertPath = "cert/cert"
}

func (c *Config) FromYAML() error {
	data, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, c)
	if err != nil {
		return err
	}

	return nil
}

func NewConfig() *Config {
	config := &Config{}
	config.SetDefaults()

	return config
}
