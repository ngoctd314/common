package env

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

type option func(*config)

func withEnvPrefix(envPrefix string) option {
	return func(c *config) {
		// if we use envPrefix => make sure it is not empty
		if strings.TrimSpace(envPrefix) != "" {
			c.SetEnvPrefix(envPrefix)
		}
	}
}

// WithFile inject config from config file into *config
func WithFile(configFile string) option {
	return func(c *config) {
		_, err := os.Stat(configFile)
		if err != nil {
			slog.Warn(fmt.Sprintf("cannot apply configuration settings from %s. Please check permissions, ensure the file exists, or ignore if reading from environment variables.", configFile))
			return
		}

		configFile = strings.TrimSpace(configFile)
		if configFile != "" {
			c.SetConfigFile(configFile)
		}

		c.AddConfigPath(".") // set current working directory as config path

		// Find and read the config file
		if err := c.ReadInConfig(); err != nil {
			slog.Warn(fmt.Sprintf("error occur when read config from file, err: '%v', config file: '%s'", err, configFile))
		}
	}
}
