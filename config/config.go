package config

import (
	"github.com/spf13/viper"
)

// This struct will store all the configration related parameters.
type Config struct {
	DatabaseHost              string `mapstructure:"DATABASE_HOST"`
	DatabasePort              string `mapstructure:"DATABASE_PORT"`
	DatabaseName              string `mapstructure:"DATABASE_NAME"`
	DatabaseUser              string `mapstructure:"DATABASE_USER"`
	DatabasePassword          string `mapstructure:"DATABASE_PASSWORD"`
	NumberOfTestPersonEntries int    `mapstructure:"NUMBER_OF_TEST_PERSON_ENTRIES"`
	ShouldSeedData            bool   `mapstructure:"SHOULD_SEED_DATA"`
}

//LoadConfig Reads the env variables and returns a config struct.
func LoadConfig(environment string) (config Config, err error) {
	viper.AddConfigPath(getEnvironmentFilePath(environment))
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	// This helps read environment variables if they have been set.
	viper.AutomaticEnv()

	// Read the configuration file
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	// Unmarshal the read config to a struct.
	// Unmarshal indicates conversion of external files to a go struct.
	// An example of named return.
	err = viper.Unmarshal(&config)
	return
}

//getEnvironmentFilePath Returns the relative path to the env file.
func getEnvironmentFilePath(environment string) string {

	return "./config/env/"
}
