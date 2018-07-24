package main

import (
	"fmt"
	"os"
)

// Args describe arguments we expect from the environment
type Args struct {
	authUser     string
	authPassword string
	httpPort     string
	inProduction bool
	rateLimit    string
}

var defaultSettings = Args{
	inProduction: false,
	httpPort:     "80",
	rateLimit:    "30-M",
}

func getArgs() (args *Args) {
	envArgs := Args{
		authUser:     os.Getenv("AUTH_USER"),
		authPassword: os.Getenv("AUTH_PASSWORD"),
		inProduction: os.Getenv("PRODUCTION") == "true",
		httpPort:     os.Getenv("HTTP_PORT"),
		rateLimit:    os.Getenv("RATE_LIMIT"),
	}

	if envArgs.inProduction && (envArgs.authUser == "" || envArgs.authPassword == "") {
		fmt.Println("IMPORTANT: When in production, you should secure the " +
			"service by enforcing HTTP authentication (`AUTH_USER`, " +
			"`AUTH_PASSWORD`). Please refer to the documentation.")
	}

	if envArgs.httpPort == "" {
		envArgs.httpPort = defaultSettings.httpPort
	}

	if envArgs.rateLimit == "" {
		envArgs.rateLimit = defaultSettings.rateLimit
	}

	return &envArgs
}
