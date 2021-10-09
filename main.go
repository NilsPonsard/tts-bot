package main

import (
	"os"
	"path"

	cli "github.com/jawher/mow.cli"
	"github.com/nilsponsard/tts-bot/internal/commands"
	"github.com/nilsponsard/tts-bot/pkg/files"
	"github.com/nilsponsard/tts-bot/pkg/verbosity"
)

// Version will be set by the script build.sh
var version string

func main() {

	app := cli.App("tts-bot", "starter project")
	app.Version("v version", version)

	defaultPath := files.ParsePath("~/.tts-bot/")

	// arguments

	var (
		verbose     = app.BoolOpt("d debug", false, "Debug mode, more verbose operations")
		appPath     = app.StringOpt("c config", path.Join(defaultPath, "config.json"), "Path to the config file")
		disableLogs = app.BoolOpt("disable-logs", false, "Disable the saving of logs")
	)

	// Executed befor the commands

	app.Before = func() {

		parsedConfigPath := *appPath
		files.EnsureFolder(parsedConfigPath)

		// create the folder for the logs

		files.EnsureFolder(path.Join(defaultPath, "test"))

		// Configure the logs

		verbosity.SetupLog(*verbose, path.Join(defaultPath, "logs.txt"), !*disableLogs)

	}

	// set subcommands

	commands.SetupCommands(app)

	// parse the args

	app.Run(os.Args)
}
