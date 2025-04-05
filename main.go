package main

import (
	"embed"
	"os"

	"github.com/AuruOS/abg/cmd"
	"github.com/AuruOS/abg/core"
	"github.com/AuruOS/orchid/cmdr"
)

var Version = "development"

//go:embed locales/*.yml
var fs embed.FS
var abgApp *cmdr.App

func main() {
	core.NewStandardAbg()

	abgApp = cmd.New(Version, fs)

	// Check if running as root, exit if so
	if core.RootCheck(false) {
		cmdr.Error.Println(abgApp.Trans("abg.errors.noRoot"))
		os.Exit(1)
	}

	// Root command
	root := cmd.NewRootCommand(Version)
	abgApp.CreateRootCommand(root, abgApp.Trans("abg.msg.help"), abgApp.Trans("abg.msg.version"))

	msgs := cmdr.UsageStrings{
		Usage:                abgApp.Trans("abg.msg.usage"),
		Aliases:              abgApp.Trans("abg.msg.aliases"),
		Examples:             abgApp.Trans("abg.msg.examples"),
		AvailableCommands:    abgApp.Trans("abg.msg.availableCommands"),
		AdditionalCommands:   abgApp.Trans("abg.msg.additionalCommands"),
		Flags:                abgApp.Trans("abg.msg.flags"),
		GlobalFlags:          abgApp.Trans("abg.msg.globalFlags"),
		AdditionalHelpTopics: abgApp.Trans("abg.msg.additionalHelpTopics"),
		MoreInfo:             abgApp.Trans("abg.msg.moreInfo"),
	}

	// Set usage strings for the app
	abgApp.SetUsageStrings(msgs)

	// Register commands
	registerCommands(root)

	// Run the app
	if err := abgApp.Run(); err != nil {
		cmdr.Error.Println(err)
		os.Exit(1)
	}
}

// registerCommands adds all available commands to the root command.
func registerCommands(root *cmdr.Command) {
	stacks := cmd.NewStacksCommand()
	root.AddCommand(stacks)

	subsystems := cmd.NewSubSystemsCommand()
	root.AddCommand(subsystems)

	pkgManagers := cmd.NewPkgManagersCommand()
	root.AddCommand(pkgManagers)

	runtimeCmds := cmd.NewRuntimeCommands()
	root.AddCommand(runtimeCmds...)
}
