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
	redisHost  string
	redisPort  string
	httpPort   string
	registries []recreate.RegistryConf
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
		redisHost:  os.Getenv("REDIS_HOST"),
		redisPort:  os.Getenv("REDIS_PORT"),
		httpPort:   os.Getenv("HTTP_PORT"),
		registries: getRegistries()}

	if envArgs.redisHost == "" {
		envArgs.redisHost = "127.0.0.1"
	}

	if envArgs.redisPort == "" {
		envArgs.redisPort = "6379"
	}

	if envArgs.httpPort == "" {
		envArgs.httpPort = "80"
	}

	return &envArgs
}
