package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/AuruOS/abg/core"
	"github.com/AuruOS/orchid/cmdr"
)

const (
	PkgManagerCmdAutoRemove = "autoRemove"
	PkgManagerCmdClean      = "clean"
	PkgManagerCmdInstall    = "install"
	PkgManagerCmdList       = "list"
	PkgManagerCmdPurge      = "purge"
	PkgManagerCmdRemove     = "remove"
	PkgManagerCmdSearch     = "search"
	PkgManagerCmdShow       = "show"
	PkgManagerCmdUpdate     = "update"
	PkgManagerCmdUpgrade    = "upgrade"
)

var PkgManagerCmdSetOrder = []string{
	PkgManagerCmdInstall,
	PkgManagerCmdUpdate,
	PkgManagerCmdRemove,
	PkgManagerCmdPurge,
	PkgManagerCmdAutoRemove,
	PkgManagerCmdClean,
	PkgManagerCmdList,
	PkgManagerCmdSearch,
	PkgManagerCmdShow,
	PkgManagerCmdUpgrade,
}

func NewPkgManagersCommand() *cmdr.Command {
	cmd := cmdr.NewCommand(
		"pkgmanagers",
		abg.Trans("pkgmanagers.description"),
		abg.Trans("pkgmanagers.description"),
		nil,
	)

	// List sub command
	listCmd := cmdr.NewCommand(
		"list",
		abg.Trans("pkgmanagers.list.description"),
		abg.Trans("pkgmanagers.list.description"),
		listPkgManagers,
	)
	listCmd.WithBoolFlag(
		cmdr.NewBoolFlag(
			"json",
			"j",
			abg.Trans("pkgmanagers.list.options.json.description"),
			false,
		),
	)

	// Show sub command
	showCmd := cmdr.NewCommand(
		"show",
		abg.Trans("pkgmanagers.show.description"),
		abg.Trans("pkgmanagers.show.description"),
		showPkgManager,
	)
	showCmd.Args = cobra.MinimumNArgs(1)

	// New sub command
	newCmd := cmdr.NewCommand(
		"new",
		abg.Trans("pkgmanagers.new.description"),
		abg.Trans("pkgmanagers.new.description"),
		newPkgManager,
	)
	setupPkgManagerFlags(newCmd)

	// Remove subcommand
	rmCmd := cmdr.NewCommand(
		"rm",
		abg.Trans("pkgmanagers.rm.description"),
		abg.Trans("pkgmanagers.rm.description"),
		rmPkgManager,
	)
	rmCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"name",
			"n",
			abg.Trans("pkgmanagers.rm.options.name.description"),
			"",
		),
	)
	rmCmd.WithBoolFlag(
		cmdr.NewBoolFlag(
			"force",
			"f",
			abg.Trans("pkgmanagers.rm.options.force.description"),
			false,
		),
	)

	// Export subcommand
	exportCmd := cmdr.NewCommand(
		"export",
		abg.Trans("pkgmanagers.export.description"),
		abg.Trans("pkgmanagers.export.description"),
		exportPkgmanager,
	)
	exportCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"name",
			"n",
			abg.Trans("pkgmanagers.export.options.name.description"),
			"",
		),
	)
	exportCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"output",
			"o",
			abg.Trans("stacks.export.options.output.description"),
			"",
		),
	)

	// Import subcommand
	importCmd := cmdr.NewCommand(
		"import",
		abg.Trans("pkgmanagers.import.description"),
		abg.Trans("pkgmanagers.import.description"),
		importPkgmanager,
	)
	importCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"input",
			"i",
			abg.Trans("pkgmanagers.import.options.input.description"),
			"",
		),
	)

	// Update subcommand
	updateCmd := cmdr.NewCommand(
		"update",
		abg.Trans("pkgmanagers.update.description"),
		abg.Trans("pkgmanagers.update.description"),
		updatePkgManager,
	)
	setupPkgManagerFlags(updateCmd)

	// Add subcommands
	cmd.AddCommand(listCmd)
	cmd.AddCommand(showCmd)
	cmd.AddCommand(newCmd)
	cmd.AddCommand(rmCmd)
	cmd.AddCommand(exportCmd)
	cmd.AddCommand(importCmd)
	cmd.AddCommand(updateCmd)

	return cmd
}

func setupPkgManagerFlags(cmd *cmdr.Command) {
	cmd.WithBoolFlag(
		cmdr.NewBoolFlag(
			"no-prompt",
			"y",
			abg.Trans("pkgmanagers.new.options.noPrompt.description"),
			false,
		),
	)
	cmd.WithStringFlag(
		cmdr.NewStringFlag(
			"name",
			"n",
			abg.Trans("pkgmanagers.new.options.name.description"),
			"",
		),
	)
	cmd.WithBoolFlag(
		cmdr.NewBoolFlag(
			"need-sudo",
			"S",
			abg.Trans("pkgmanagers.new.options.needSudo.description"),
			false,
		),
	)
	cmd.WithStringFlag(
		cmdr.NewStringFlag(
			"autoremove",
			"a",
			abg.Trans("pkgmanagers.new.options.autoremove.description"),
			"",
		),
	)
	cmd.WithStringFlag(
		cmdr.NewStringFlag(
			"clean",
			"c",
			abg.Trans("pkgmanagers.new.options.clean.description"),
			"",
		),
	)
	cmd.WithStringFlag(
		cmdr.NewStringFlag(
			"install",
			"i",
			abg.Trans("pkgmanagers.new.options.install.description"),
			"",
		),
	)
	cmd.WithStringFlag(
		cmdr.NewStringFlag(
			"list",
			"l",
			abg.Trans("pkgmanagers.new.options.list.description"),
			"",
		),
	)
	cmd.WithStringFlag(
		cmdr.NewStringFlag(
			"purge",
			"p",
			abg.Trans("pkgmanagers.new.options.purge.description"),
			"",
		),
	)
	cmd.WithStringFlag(
		cmdr.NewStringFlag(
			"remove",
			"r",
			abg.Trans("pkgmanagers.new.options.remove.description"),
			"",
		),
	)
	cmd.WithStringFlag(
		cmdr.NewStringFlag(
			"search",
			"s",
			abg.Trans("pkgmanagers.new.options.search.description"),
			"",
		),
	)
	cmd.WithStringFlag(
		cmdr.NewStringFlag(
			"show",
			"w",
			abg.Trans("pkgmanagers.new.options.show.description"),
			"",
		),
	)
	cmd.WithStringFlag(
		cmdr.NewStringFlag(
			"update",
			"u",
			abg.Trans("pkgmanagers.new.options.update.description"),
			"",
		),
	)
	cmd.WithStringFlag(
		cmdr.NewStringFlag(
			"upgrade",
			"U",
			abg.Trans("pkgmanagers.new.options.upgrade.description"),
			"",
		),
	)
}

func listPkgManagers(cmd *cobra.Command, args []string) error {
	jsonFlag, _ := cmd.Flags().GetBool("json")
	pkgManagers := core.ListPkgManagers()

	if !jsonFlag {
		if len(pkgManagers) == 0 {
			cmdr.Info.Printfln(abg.Trans("pkgmanagers.list.info.noPkgManagers"))
			return nil
		}

		cmdr.Info.Printfln(abg.Trans("pkgmanagers.list.info.foundPkgManagers"), len(pkgManagers))

		table := core.CreateApxTable(os.Stdout)
		table.SetHeader([]string{abg.Trans("pkgmanagers.labels.name"), abg.Trans("pkgmanagers.labels.builtIn")})

		for _, pm := range pkgManagers {
			builtIn := abg.Trans("terminal.no")
			if pm.BuiltIn {
				builtIn = abg.Trans("terminal.yes")
			}
			table.Append([]string{pm.Name, builtIn})
		}
		table.Render()
		return nil
	}

	jsonData, err := json.MarshalIndent(pkgManagers, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal package managers: %w", err)
	}
	fmt.Println(string(jsonData))
	return nil
}

func showPkgManager(cmd *cobra.Command, args []string) error {
	pkgManager, err := core.LoadPkgManager(args[0])
	if err != nil {
		return fmt.Errorf("failed to load package manager: %w", err)
	}

	table := core.CreateApxTable(os.Stdout)
	table.Append([]string{abg.Trans("pkgmanagers.labels.name"), pkgManager.Name})
	table.Append([]string{"NeedSudo", fmt.Sprintf("%t", pkgManager.NeedSudo)})
	table.Append([]string{"AutoRemove", pkgManager.CmdAutoRemove})
	table.Append([]string{"Clean", pkgManager.CmdClean})
	table.Append([]string{"Install", pkgManager.CmdInstall})
	table.Append([]string{"List", pkgManager.CmdList})
	table.Append([]string{"Purge", pkgManager.CmdPurge})
	table.Append([]string{"Remove", pkgManager.CmdRemove})
	table.Append([]string{"Search", pkgManager.CmdSearch})
	table.Append([]string{"Show", pkgManager.CmdShow})
	table.Append([]string{"Update", pkgManager.CmdUpdate})
	table.Append([]string{"Upgrade", pkgManager.CmdUpgrade})
	table.Render()

	return nil
}

func newPkgManager(cmd *cobra.Command, args []string) error {
	var (
		noPrompt, _   = cmd.Flags().GetBool("no-prompt")
		name, _       = cmd.Flags().GetString("name")
		needSudo, _   = cmd.Flags().GetBool("need-sudo")
		autoRemove, _ = cmd.Flags().GetString("autoremove")
		clean, _      = cmd.Flags().GetString("clean")
		install, _    = cmd.Flags().GetString("install")
		list, _       = cmd.Flags().GetString("list")
		purge, _      = cmd.Flags().GetString("purge")
		remove, _     = cmd.Flags().GetString("remove")
		search, _     = cmd.Flags().GetString("search")
		show, _       = cmd.Flags().GetString("show")
		update, _     = cmd.Flags().GetString("update")
		upgrade, _    = cmd.Flags().GetString("upgrade")
	)

	reader := bufio.NewReader(os.Stdin)

	// Validate and get name
	if name == "" {
		if noPrompt {
			return fmt.Errorf(abg.Trans("pkgmanagers.new.error.noName"))
		}
		name = promptForInput(reader, abg.Trans("pkgmanagers.new.info.askName"), "")
		if name == "" {
			return fmt.Errorf(abg.Trans("pkgmanagers.new.error.emptyName"))
		}
	}

	// Get sudo requirement if not specified
	if !needSudo && !noPrompt {
		needSudo = promptForConfirmation(reader, abg.Trans("pkgmanagers.new.info.askSudo"))
	}

	// Collect commands
	commands := map[string]*string{
		PkgManagerCmdAutoRemove: &autoRemove,
		PkgManagerCmdClean:      &clean,
		PkgManagerCmdInstall:    &install,
		PkgManagerCmdList:       &list,
		PkgManagerCmdPurge:      &purge,
		PkgManagerCmdRemove:     &remove,
		PkgManagerCmdSearch:     &search,
		PkgManagerCmdShow:       &show,
		PkgManagerCmdUpdate:     &update,
		PkgManagerCmdUpgrade:    &upgrade,
	}

	for _, cmdName := range PkgManagerCmdSetOrder {
		cmdPtr := commands[cmdName]
		if *cmdPtr == "" {
			if noPrompt {
				return fmt.Errorf(abg.Trans("pkgmanagers.new.error.noCommand"), cmdName)
			}

			defaultValue := ""
			if cmdName == PkgManagerCmdPurge || cmdName == PkgManagerCmdAutoRemove {
				defaultValue = remove
			}

			*cmdPtr = promptForInput(
				reader,
				fmt.Sprintf(abg.Trans("pkgmanagers.new.info.askCommandWithDefault"), cmdName, defaultValue),
				defaultValue,
			)
			if *cmdPtr == "" && cmdName != PkgManagerCmdPurge && cmdName != PkgManagerCmdAutoRemove {
				return fmt.Errorf(abg.Trans("pkgmanagers.new.error.emptyCommand"), cmdName)
			}
		}
	}

	// Check if package manager exists
	if core.PkgManagerExists(name) {
		if noPrompt {
			return fmt.Errorf(abg.Trans("pkgmanagers.new.error.alreadyExists"), name)
		}

		if !promptForConfirmation(reader, fmt.Sprintf(abg.Trans("pkgmanagers.new.info.askOverwrite"), name)) {
			cmdr.Info.Println(abg.Trans("terminal.info.aborting"))
			return nil
		}
	}

	// Create and save package manager
	pkgManager := core.NewPkgManager(
		name, needSudo,
		autoRemove, clean, install, list, purge, remove, search, show, update, upgrade,
		false,
	)

	if err := pkgManager.Save(); err != nil {
		return fmt.Errorf("failed to save package manager: %w", err)
	}

	cmdr.Success.Printfln(abg.Trans("pkgmanagers.new.success"), name)
	return nil
}

func rmPkgManager(cmd *cobra.Command, args []string) error {
	var (
		name, _  = cmd.Flags().GetString("name")
		force, _ = cmd.Flags().GetBool("force")
	)

	if name == "" {
		return fmt.Errorf(abg.Trans("pkgmanagers.rm.error.noName"))
	}

	pkgManager, err := core.LoadPkgManager(name)
	if err != nil {
		return fmt.Errorf("failed to load package manager: %w", err)
	}

	// Check if package manager is in use
	stacks := core.ListStackForPkgManager(pkgManager.Name)
	if len(stacks) > 0 {
		cmdr.Error.Printf(abg.Trans("pkgmanagers.rm.error.inUse"), len(stacks))
		table := core.CreateApxTable(os.Stdout)
		table.SetHeader([]string{
			abg.Trans("pkgmanagers.labels.name"),
			"Base",
			"Packages",
			"PkgManager",
			abg.Trans("pkgmanagers.labels.builtIn"),
		})

		for _, stack := range stacks {
			builtIn := abg.Trans("terminal.no")
			if stack.BuiltIn {
				builtIn = abg.Trans("terminal.yes")
			}
			table.Append([]string{
				stack.Name,
				stack.Base,
				strings.Join(stack.Packages, ", "),
				stack.PkgManager,
				builtIn,
			})
		}
		table.Render()
		return nil
	}

	// Confirm deletion if not forced
	if !force {
		reader := bufio.NewReader(os.Stdin)
		force = promptForConfirmation(
			reader,
			fmt.Sprintf(abg.Trans("pkgmanagers.rm.info.askConfirmation"), pkgManager.Name),
		)
	}

	if !force {
		cmdr.Info.Printfln(abg.Trans("pkgmanagers.rm.info.aborting"), pkgManager.Name)
		return nil
	}

	if err := pkgManager.Remove(); err != nil {
		return fmt.Errorf("failed to remove package manager: %w", err)
	}

	cmdr.Info.Printfln(abg.Trans("pkgmanagers.rm.info.success"), pkgManager.Name)
	return nil
}

func exportPkgmanager(cmd *cobra.Command, args []string) error {
	var (
		name, _   = cmd.Flags().GetString("name")
		output, _ = cmd.Flags().GetString("output")
	)

	if name == "" {
		return fmt.Errorf(abg.Trans("pkgmanagers.export.error.noName"))
	}
	if output == "" {
		return fmt.Errorf(abg.Trans("pkgmanagers.export.error.noOutput"))
	}

	pkgManager, err := core.LoadPkgManager(name)
	if err != nil {
		return fmt.Errorf("failed to load package manager: %w", err)
	}

	if err := pkgManager.Export(output); err != nil {
		return fmt.Errorf("failed to export package manager: %w", err)
	}

	cmdr.Info.Printfln(abg.Trans("pkgmanagers.export.info.success"), pkgManager.Name, output)
	return nil
}

func importPkgmanager(cmd *cobra.Command, args []string) error {
	input, _ := cmd.Flags().GetString("input")
	if input == "" {
		return fmt.Errorf(abg.Trans("pkgmanagers.import.error.noInput"))
	}

	pkgmanager, err := core.LoadPkgManagerFromPath(input)
	if err != nil {
		return fmt.Errorf("failed to load package manager from %s: %w", input, err)
	}

	if err := pkgmanager.Save(); err != nil {
		return fmt.Errorf("failed to save package manager: %w", err)
	}

	cmdr.Info.Printfln(abg.Trans("pkgmanagers.import.info.success"), pkgmanager.Name)
	return nil
}

func updatePkgManager(cmd *cobra.Command, args []string) error {
	var (
		name, _       = cmd.Flags().GetString("name")
		needSudo, _   = cmd.Flags().GetBool("need-sudo")
		noPrompt, _   = cmd.Flags().GetBool("no-prompt")
		autoRemove, _ = cmd.Flags().GetString("autoremove")
		clean, _      = cmd.Flags().GetString("clean")
		install, _    = cmd.Flags().GetString("install")
		list, _       = cmd.Flags().GetString("list")
		purge, _      = cmd.Flags().GetString("purge")
		remove, _     = cmd.Flags().GetString("remove")
		search, _     = cmd.Flags().GetString("search")
		show, _       = cmd.Flags().GetString("show")
		update, _     = cmd.Flags().GetString("update")
		upgrade, _    = cmd.Flags().GetString("upgrade")
	)

	if name == "" && (len(args) == 0 || args[0] == "") {
		return fmt.Errorf(abg.Trans("pkgmanagers.update.error.noName"))
	}
	if name == "" {
		name = args[0]
	}

	pkgmanager, err := core.LoadPkgManager(name)
	if err != nil {
		return fmt.Errorf("failed to load package manager: %w", err)
	}

	if pkgmanager.BuiltIn {
		return fmt.Errorf(abg.Trans("pkgmanagers.update.error.builtIn"))
	}

	// Update commands
	commands := map[string]*string{
		PkgManagerCmdAutoRemove: &autoRemove,
		PkgManagerCmdClean:      &clean,
		PkgManagerCmdInstall:    &install,
		PkgManagerCmdList:       &list,
		PkgManagerCmdPurge:      &purge,
		PkgManagerCmdRemove:     &remove,
		PkgManagerCmdSearch:     &search,
		PkgManagerCmdShow:       &show,
		PkgManagerCmdUpdate:     &update,
		PkgManagerCmdUpgrade:    &upgrade,
	}

	reader := bufio.NewReader(os.Stdin)
	for cmdName, cmdPtr := range commands {
		if *cmdPtr == "" && !noPrompt {
			defaultValue := pkgmanager.GetCommand(cmdName)
			*cmdPtr = promptForInput(
				reader,
				fmt.Sprintf(abg.Trans("pkgmanagers.update.info.askNewCommand"), cmdName, defaultValue),
				defaultValue,
			)
		} else if *cmdPtr == "" {
			return fmt.Errorf(abg.Trans("pkgmanagers.update.error.missingCommand"), cmdName)
		}
	}

	// Update package manager properties
	pkgmanager.NeedSudo = needSudo
	pkgmanager.CmdAutoRemove = autoRemove
	pkgmanager.CmdClean = clean
	pkgmanager.CmdInstall = install
	pkgmanager.CmdList = list
	pkgmanager.CmdPurge = purge
	pkgmanager.CmdRemove = remove
	pkgmanager.CmdSearch = search
	pkgmanager.CmdShow = show
	pkgmanager.CmdUpdate = update
	pkgmanager.CmdUpgrade = upgrade

	if err := pkgmanager.Save(); err != nil {
		return fmt.Errorf("failed to save package manager: %w", err)
	}

	cmdr.Info.Printfln(abg.Trans("pkgmanagers.update.info.success"), name)
	return nil
}

// Helper functions
func promptForInput(reader *bufio.Reader, prompt string, defaultValue string) string {
	cmdr.Info.Println(prompt)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultValue
	}
	return input
}

func promptForConfirmation(reader *bufio.Reader, prompt string) bool {
	for {
		cmdr.Info.Printf("%s [y/N]: ", prompt)
		answer, _ := reader.ReadString('\n')
		answer = strings.ToLower(strings.TrimSpace(answer))

		switch answer {
		case "y", "yes":
			return true
		case "n", "no", "":
			return false
		default:
			cmdr.Warning.Println(abg.Trans("errors.invalidChoice"))
		}
	}
}
