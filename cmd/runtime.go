package cmd

import (
	"fmt"
	"slices"

	"github.com/spf13/cobra"
	"github.com/AuruOS/abg/core"
	"github.com/AuruOS/orchid/cmdr"
)

func NewRuntimeCommands() []*cmdr.Command {
	var commands []*cmdr.Command

	subSystems, err := core.ListSubSystems(false, false)
	if err != nil {
		return []*cmdr.Command{}
	}

	handleFunc := func(subSystem *core.SubSystem, reqFunc func(*core.SubSystem, string, *cobra.Command, []string) error) func(cmd *cobra.Command, args []string) error {
		return func(cmd *cobra.Command, args []string) error {
			return reqFunc(subSystem, cmd.Name(), cmd, args)
		}
	}

	for _, subSystem := range subSystems {
		subSystemCmd := cmdr.NewCommand(
			subSystem.Name,
			abg.Trans("runtimeCommand.description"),
			abg.Trans("runtimeCommand.description"),
			nil,
		)

		autoRemoveCmd := cmdr.NewCommand(
			"autoremove",
			abg.Trans("runtimeCommand.autoremove.description"),
			abg.Trans("runtimeCommand.autoremove.description"),
			handleFunc(subSystem, runPkgCmd),
		)
		cleanCmd := cmdr.NewCommand(
			"clean",
			abg.Trans("runtimeCommand.clean.description"),
			abg.Trans("runtimeCommand.clean.description"),
			handleFunc(subSystem, runPkgCmd),
		)

		installCmd := cmdr.NewCommand(
			"install",
			abg.Trans("runtimeCommand.install.description"),
			abg.Trans("runtimeCommand.install.description"),
			handleFunc(subSystem, runPkgCmd),
		)
		installCmd.WithBoolFlag(
			cmdr.NewBoolFlag(
				"no-export",
				"n",
				abg.Trans("runtimeCommand.install.options.noExport.description"),
				false,
			),
		)

		listCmd := cmdr.NewCommand(
			"list",
			abg.Trans("runtimeCommand.list.description"),
			abg.Trans("runtimeCommand.list.description"),
			handleFunc(subSystem, runPkgCmd),
		)
		purgeCmd := cmdr.NewCommand(
			"purge",
			abg.Trans("runtimeCommand.purge.description"),
			abg.Trans("runtimeCommand.purge.description"),
			handleFunc(subSystem, runPkgCmd),
		)
		removeCmd := cmdr.NewCommand(
			"remove",
			abg.Trans("runtimeCommand.remove.description"),
			abg.Trans("runtimeCommand.remove.description"),
			handleFunc(subSystem, runPkgCmd),
		)
		searchCmd := cmdr.NewCommand(
			"search",
			abg.Trans("runtimeCommand.search.description"),
			abg.Trans("runtimeCommand.search.description"),
			handleFunc(subSystem, runPkgCmd),
		)
		showCmd := cmdr.NewCommand(
			"show",
			abg.Trans("runtimeCommand.show.description"),
			abg.Trans("runtimeCommand.show.description"),
			handleFunc(subSystem, runPkgCmd),
		)
		updateCmd := cmdr.NewCommand(
			"update",
			abg.Trans("runtimeCommand.update.description"),
			abg.Trans("runtimeCommand.update.description"),
			handleFunc(subSystem, runPkgCmd),
		)
		upgradeCmd := cmdr.NewCommand(
			"upgrade",
			abg.Trans("runtimeCommand.upgrade.description"),
			abg.Trans("runtimeCommand.upgrade.description"),
			handleFunc(subSystem, runPkgCmd),
		)
		runCmd := cmdr.NewCommand(
			"run",
			abg.Trans("runtimeCommand.run.description"),
			abg.Trans("runtimeCommand.run.description"),
			handleFunc(subSystem, runPkgCmd),
		)
		enterCmd := cmdr.NewCommand(
			"enter",
			abg.Trans("runtimeCommand.enter.description"),
			abg.Trans("runtimeCommand.enter.description"),
			handleFunc(subSystem, runPkgCmd),
		)
		exportCmd := cmdr.NewCommand(
			"export",
			abg.Trans("runtimeCommand.export.description"),
			abg.Trans("runtimeCommand.export.description"),
			handleFunc(subSystem, runPkgCmd),
		)
		exportCmd.WithStringFlag(
			cmdr.NewStringFlag(
				"app-name",
				"a",
				abg.Trans("runtimeCommand.export.options.appName.description"),
				"",
			),
		)
		exportCmd.WithStringFlag(
			cmdr.NewStringFlag(
				"bin",
				"b",
				abg.Trans("runtimeCommand.export.options.bin.description"),
				"",
			),
		)
		exportCmd.WithStringFlag(
			cmdr.NewStringFlag(
				"bin-output",
				"o",
				abg.Trans("runtimeCommand.export.options.binOutput.description"),
				"",
			),
		)
		unexportCmd := cmdr.NewCommand(
			"unexport",
			abg.Trans("runtimeCommand.unexport.description"),
			abg.Trans("runtimeCommand.unexport.description"),
			handleFunc(subSystem, runPkgCmd),
		)
		unexportCmd.WithStringFlag(
			cmdr.NewStringFlag(
				"app-name",
				"a",
				abg.Trans("runtimeCommand.unexport.options.appName.description"),
				"",
			),
		)
		unexportCmd.WithStringFlag(
			cmdr.NewStringFlag(
				"bin",
				"b",
				abg.Trans("runtimeCommand.unexport.options.bin.description"),
				"",
			),
		)
		unexportCmd.WithStringFlag(
			cmdr.NewStringFlag(
				"bin-output",
				"o",
				abg.Trans("runtimeCommand.unexport.options.binOutput.description"),
				"",
			),
		)

		startCmd := cmdr.NewCommand(
			"start",
			abg.Trans("runtimeCommand.start.description"),
			abg.Trans("runtimeCommand.start.description"),
			handleFunc(subSystem, runPkgCmd),
		)

		stopCmd := cmdr.NewCommand(
			"stop",
			abg.Trans("runtimeCommand.stop.description"),
			abg.Trans("runtimeCommand.stop.description"),
			handleFunc(subSystem, runPkgCmd),
		)

		subSystemCmd.AddCommand(autoRemoveCmd)
		subSystemCmd.AddCommand(cleanCmd)
		subSystemCmd.AddCommand(installCmd)
		subSystemCmd.AddCommand(listCmd)
		subSystemCmd.AddCommand(purgeCmd)
		subSystemCmd.AddCommand(removeCmd)
		subSystemCmd.AddCommand(searchCmd)
		subSystemCmd.AddCommand(showCmd)
		subSystemCmd.AddCommand(updateCmd)
		subSystemCmd.AddCommand(upgradeCmd)
		subSystemCmd.AddCommand(runCmd)
		subSystemCmd.AddCommand(enterCmd)
		subSystemCmd.AddCommand(exportCmd)
		subSystemCmd.AddCommand(unexportCmd)
		subSystemCmd.AddCommand(startCmd)
		subSystemCmd.AddCommand(stopCmd)

		commands = append(commands, subSystemCmd)
	}

	return commands
}

var baseCmds = []string{"run", "enter", "export", "unexport", "start", "stop"}

// isBaseCommand informs whether the command is a subsystem-base command
// (e.g. run, enter) instead of a subsystem-specific one (e.g. install, update)
func isBaseCommand(command string) bool {
	return slices.Contains(baseCmds, command)
}

// pkgManagerCommands maps command line arguments into package manager commands
func pkgManagerCommands(pkgManager *core.PkgManager, command string) (string, error) {
	switch command {
	case "autoremove":
		return pkgManager.CmdAutoRemove, nil
	case "clean":
		return pkgManager.CmdClean, nil
	case "install":
		return pkgManager.CmdInstall, nil
	case "list":
		return pkgManager.CmdList, nil
	case "purge":
		return pkgManager.CmdPurge, nil
	case "remove":
		return pkgManager.CmdRemove, nil
	case "search":
		return pkgManager.CmdSearch, nil
	case "show":
		return pkgManager.CmdShow, nil
	case "update":
		return pkgManager.CmdUpdate, nil
	case "upgrade":
		return pkgManager.CmdUpgrade, nil
	default:
		return "", fmt.Errorf(abg.Trans("abg.errors.unknownCommand"), command)
	}
}

func runPkgCmd(subSystem *core.SubSystem, command string, cmd *cobra.Command, args []string) error {
	if !isBaseCommand(command) {
		pkgManager, err := subSystem.Stack.GetPkgManager()
		if err != nil {
			return fmt.Errorf(abg.Trans("runtimeCommand.error.cantAccessPkgManager"), err)
		}

		realCommand, err := pkgManagerCommands(pkgManager, command)
		if err != nil {
			return err
		}

		if command == "remove" {
			exportedN, err := subSystem.UnexportDesktopEntries(args...)
			if err == nil {
				cmdr.Info.Printfln(abg.Trans("runtimeCommand.info.unexportedApps"), exportedN)
			}
		}

		finalArgs := pkgManager.GenCmd(realCommand, args...)
		_, err = subSystem.Exec(false, false, finalArgs...)
		if err != nil {
			return fmt.Errorf(abg.Trans("runtimeCommand.error.executingCommand"), err)
		}

		if command == "install" && !cmd.Flag("no-export").Changed {
			exportedN, err := subSystem.ExportDesktopEntries(args...)
			if err == nil {
				cmdr.Info.Printfln(abg.Trans("runtimeCommand.info.exportedApps"), exportedN)
			}
		}

		return nil
	}

	if command == "run" {
		_, err := subSystem.Exec(false, false, args...)
		if err != nil {
			return fmt.Errorf(abg.Trans("runtimeCommand.error.executingCommand"), err)
		}

		return nil
	}

	if command == "enter" {
		err := subSystem.Enter()
		if err != nil {
			return fmt.Errorf(abg.Trans("runtimeCommand.error.enteringContainer"), err)
		}

		return nil
	}

	if command == "export" || command == "unexport" {
		appName, _ := cmd.Flags().GetString("app-name")
		bin, _ := cmd.Flags().GetString("bin")
		binOutput, _ := cmd.Flags().GetString("bin-output")

		return handleExport(subSystem, command, appName, bin, binOutput)
	}

	if command == "start" {
		cmdr.Info.Printfln(abg.Trans("runtimeCommand.info.startingContainer"), subSystem.Name)
		err := subSystem.Start()
		if err != nil {
			return fmt.Errorf(abg.Trans("runtimeCommand.error.startingContainer"), err)
		}

		cmdr.Info.Printfln(abg.Trans("runtimeCommand.info.startedContainer"))
	}

	if command == "stop" {
		cmdr.Info.Printfln(abg.Trans("runtimeCommand.info.stoppingContainer"), subSystem.Name)
		err := subSystem.Stop()
		if err != nil {
			return fmt.Errorf(abg.Trans("runtimeCommand.error.stoppingContainer"), err)
		}

		cmdr.Info.Printfln(abg.Trans("runtimeCommand.info.stoppedContainer"))
	}

	return nil
}

func handleExport(subSystem *core.SubSystem, command, appName, bin, binOutput string) error {
	if appName == "" && bin == "" {
		return fmt.Errorf(abg.Trans("runtimeCommand.error.noAppNameOrBin"))
	}

	if appName != "" && bin != "" {
		return fmt.Errorf(abg.Trans("runtimeCommand.error.sameAppOrBin"))
	}

	if command == "export" {
		if appName != "" {
			err := subSystem.ExportDesktopEntry(appName)
			if err != nil {
				return fmt.Errorf(abg.Trans("runtimeCommand.error.exportingApp"), err)
			}

			cmdr.Info.Printfln(abg.Trans("runtimeCommand.info.exportedApp"), appName)
		} else {
			err := subSystem.ExportBin(bin, binOutput)
			if err != nil {
				return fmt.Errorf(abg.Trans("runtimeCommand.error.exportingBin"), err)
			}

			cmdr.Info.Printfln(abg.Trans("runtimeCommand.info.exportedBin"), bin)
		}
	} else {
		if appName != "" {
			err := subSystem.UnexportDesktopEntry(appName)
			if err != nil {
				return fmt.Errorf(abg.Trans("runtimeCommand.error.unexportingApp"), err)
			}

			cmdr.Info.Printfln(abg.Trans("runtimeCommand.info.unexportedApp"), appName)
		} else {
			err := subSystem.UnexportBin(bin, binOutput)
			if err != nil {
				return fmt.Errorf(abg.Trans("runtimeCommand.error.unexportingBin"), err)
			}

			cmdr.Info.Printfln(abg.Trans("runtimeCommand.info.unexportedBin"), bin)
		}
	}

	return nil
}
