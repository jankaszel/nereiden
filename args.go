package main

import (
	"fmt"
	"github.com/caarlos0/env"
)

// Args describe arguments we expect from the environment
type Args struct {
	AllowedOrigins   []string `env:"ALLOWED_ORIGINS" envSeparator:"," envDefault:"*"`
	HTTPPort         string   `env:"HTTP_PORT" envDefault:"80"`
	InProduction     bool     `env:"PRODUCTION" envDefault:"false"`
	LetsEncryptEmail string   `env:"LETS_ENCRYPT_EMAIL"`
	ProxyNetworkName string   `env:"PROXY_NETWORK_NAME" envDefault:"nginx-proxy"`
	RateLimit        string   `env:"RATE_LIMIT" envDefault:"30-M"`
}

func getArgs() Args {
	args := Args{}
	err := env.Parse(&args)

	if err != nil {
		panic(fmt.Sprintf("Error while parsing configuration: %+v\n", err))
	}

	if args.LetsEncryptEmail == "" {
		panic("You must specify an email address in order to obtain certificates " +
			"from Let's Encrypt. Please refer to the documentation.")
	}

	return args
}
