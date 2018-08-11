package main

import (
	"os"
)

// Args describe arguments we expect from the environment
type Args struct {
	httpPort         string
	inProduction     bool
	letsEncryptEmail string
	rateLimit        string
}

var defaultSettings = Args{
	inProduction: false,
	httpPort:     "80",
	rateLimit:    "30-M",
}

func getArgs() (args *Args) {
	envArgs := Args{
		httpPort:         os.Getenv("HTTP_PORT"),
		inProduction:     os.Getenv("PRODUCTION") == "true",
		letsEncryptEmail: os.Getenv("LETS_ENCRYPT_EMAIL"),
		rateLimit:        os.Getenv("RATE_LIMIT"),
	}

	if envArgs.letsEncryptEmail == "" {
		panic("You must specify an email address in order to obtain certificates " +
			"from Let's Encrypt. Please refer to the documentation.")
	}

	if envArgs.httpPort == "" {
		envArgs.httpPort = defaultSettings.httpPort
	}

	if envArgs.rateLimit == "" {
		envArgs.rateLimit = defaultSettings.rateLimit
	}

	return &envArgs
}
