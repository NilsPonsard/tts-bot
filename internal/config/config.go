package config

import (
	"encoding/json"
	"os"
	"path"

	cli "github.com/jawher/mow.cli"
	"github.com/nilsponsard/tts-bot/pkg/verbosity"
)

type Config struct {
	Token string `json:"token"`
}

var (
	confPath string
	config   Config
)

func Load(configPath string) {

	confPath = configPath

	err := os.MkdirAll(path.Base(configPath), 0700)

	if err != nil {
		verbosity.Error(err)
		cli.Exit(1)
	}

	content, err := os.ReadFile(configPath)

	if os.IsNotExist(err) {
		// TODO : save config
	}

	if err != nil {
		verbosity.Error(err)
		cli.Exit(1)
	}

	err = json.Unmarshal(content, &config)

	if err != nil {
		verbosity.Error(err)
		cli.Exit(1)
	}

	verbosity.Debug(config)

}
