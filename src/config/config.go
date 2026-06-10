package config

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Addr                 string
	EvictionLoopInterval int
	MaxKeysAmount           int
	InitialSize          uint32
}

func loadConfig(path string) (map[string]string, error) {
	config := make(map[string]string)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			config[key] = value
		}
	}
	return config, scanner.Err()
}

func GetConfig() *Config {
	keyValues, err := loadConfig(".conf")
	if err != nil {
		panic(err.Error())
	}

	Addr := keyValues["Addr"]
	if Addr == "" {
		Addr = "127.0.0.1:9990"
	}
	EvictLoopInterval, err := strconv.ParseInt(keyValues["EvictLoopInterval"], 10, 32)
	if err != nil {
		EvictLoopInterval = 5 // default
	}
	MaxKeysAmount, err := strconv.ParseInt(keyValues["MaxKeysAmount"], 10, 32)
	if err != nil {
		MaxKeysAmount = 500 // default
	}
	InitialSize, err := strconv.ParseUint(keyValues["InitialSize"], 10, 32)
	if err != nil {
		InitialSize = 100 // default
	}

	var config Config

	config.Addr = Addr
	config.EvictionLoopInterval = int(EvictLoopInterval)
	config.MaxKeysAmount = int(MaxKeysAmount)
	config.InitialSize = uint32(InitialSize)

	return &config
}
