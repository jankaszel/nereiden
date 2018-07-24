package main

import (
	"encoding/json"
	"github.com/fallafeljan/docker-recreate/lib"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Args describe arguments we expect from the environment
type Args struct {
	httpPort    string
	rateLimit   string
	redisHost   string
	redisPort   string
	redisPrefix string
	registries  []recreate.RegistryConf
}

var defaultSettings = Args{
	httpPort:    "80",
	rateLimit:   "30-M",
	redisHost:   "127.0.0.1",
	redisPort:   "6379",
	redisPrefix: "token",
}

func getRegistries() (registries []recreate.RegistryConf) {
	emptyRegistries := []recreate.RegistryConf{}
	ex, err := os.Executable()
	if err != nil {
		return emptyRegistries
	}

	cwd := filepath.Dir(ex)
	filePath := strings.Join([]string{
		cwd,
		"registries.json"},
		"/")

	file, err := os.Open(filePath)
	if err != nil {
		return emptyRegistries
	}

	defer file.Close()

	var parsedRegistries []recreate.RegistryConf
	byteValue, _ := ioutil.ReadAll(file)
	err = json.Unmarshal(byteValue, &parsedRegistries)

	if err != nil {
		return emptyRegistries
	}

	return parsedRegistries
}

func getArgs() (args *Args) {
	envArgs := Args{
		httpPort:    os.Getenv("HTTP_PORT"),
		rateLimit:   os.Getenv("RATE_LIMIT"),
		redisHost:   os.Getenv("REDIS_HOST"),
		redisPort:   os.Getenv("REDIS_PORT"),
		redisPrefix: os.Getenv("REDIS_PREFIX"),
		registries:  getRegistries(),
	}

	if envArgs.httpPort == "" {
		envArgs.httpPort = defaultSettings.httpPort
	}

	if envArgs.rateLimit == "" {
		envArgs.rateLimit = defaultSettings.rateLimit
	}

	if envArgs.redisHost == "" {
		envArgs.redisHost = defaultSettings.redisHost
	}

	if envArgs.redisPort == "" {
		envArgs.redisPort = defaultSettings.redisPort
	}

	if envArgs.redisPrefix == "" {
		envArgs.redisPrefix = defaultSettings.redisPrefix
	}

	return &envArgs
}
