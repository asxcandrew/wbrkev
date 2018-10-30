package config

import (
	"github.com/kelseyhightower/envconfig"
)

type conf struct {
	DBHost string `required:"true"`
	DBPort string `default:"5432"`
	DBUser string `required:"true"`
	DBPass string `required:"true"`
	DBName string `required:"true"`
	TLS    bool
}

var C conf

func init() {
	C = conf{}
	err := envconfig.Process("keeper", &C)

	if err != nil {
		panic(err)
	}
}
