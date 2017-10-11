package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type Config struct {
	Username    string
	AccessToken string
	MainRepo    []string
	CoreRepos   [][]string
	OtherRepos  [][]string
}

//ConfigFromFile Loads a config from a json file
func ConfigFromFile(path string) (Config, error) {

	var cfg Config

	data, err := ioutil.ReadFile(path)

	if err != nil {
		return Config{}, errors.New("Config file not found at " + path)
	}

	unMarshalErr := json.Unmarshal(data, &cfg)
	if unMarshalErr != nil {
		return Config{}, errors.New("Invalid config file format")
	}

	if cfg.AccessToken == "" {
		return Config{}, errors.New("AccessToken must be set in the config file. \n You can create a \"Personal access token\" here https://github.com/settings/tokens")
	}

	if cfg.Username == "" {
		return Config{}, errors.New("Username must be set in the config file")
	}

	if len(cfg.MainRepo) < 1 {
		return Config{}, errors.New("You must provide the main repo used to manage the release. In the format [[owner][repo]]")
	}

	if len(cfg.CoreRepos) < 1 {
		return Config{}, errors.New("You must provide at least one repo in the config file. In the format [[owner][repo]]")
	}

	if len(cfg.OtherRepos) < 1 {
		cfg.OtherRepos = nil
	}

	return cfg, nil
}
