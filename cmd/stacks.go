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

func NewStacksCommand() *cmdr.Command {
	// Root command
	cmd := cmdr.NewCommand(
		"stacks",
		abg.Trans("stacks.description"),
		abg.Trans("stacks.description"),
		nil,
	)

	// List subcommand
	listCmd := cmdr.NewCommand(
		"list",
		abg.Trans("stacks.list.description"),
		abg.Trans("stacks.list.description"),
		listStacks,
	)

	listCmd.WithBoolFlag(
		cmdr.NewBoolFlag(
			"json",
			"j",
			abg.Trans("stacks.list.options.json.description"),
			false,
		),
	)

	// Show subcommand
	showCmd := cmdr.NewCommand(
		"show",
		abg.Trans("stacks.show.description"),
		abg.Trans("stacks.show.description"),
		showStack,
	)
	showCmd.Args = cobra.MinimumNArgs(1)

	// New subcommand
	newCmd := cmdr.NewCommand(
		"new",
		abg.Trans("stacks.new.description"),
		abg.Trans("stacks.new.description"),
		newStack,
	)
	newCmd.WithBoolFlag(
		cmdr.NewBoolFlag(
			"no-prompt",
			"y",
			abg.Trans("stacks.new.options.noPrompt.description"),
			false,
		),
	)
	newCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"name",
			"n",
			abg.Trans("stacks.new.options.name.description"),
			"",
		),
	)
	newCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"base",
			"b",
			abg.Trans("stacks.new.options.base.description"),
			"",
		),
	)
	newCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"packages",
			"p",
			abg.Trans("stacks.new.options.packages.description"),
			"",
		),
	)
	newCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"pkg-manager",
			"k",
			abg.Trans("stacks.new.options.pkgManager.description"),
			"",
		),
	)

	// Update subcommand
	updateCmd := cmdr.NewCommand(
		"update",
		abg.Trans("stacks.update.description"),
		abg.Trans("stacks.update.description"),
		updateStack,
	)
	updateCmd.WithBoolFlag(
		cmdr.NewBoolFlag(
			"no-prompt",
			"y",
			abg.Trans("stacks.update.options.noPrompt.description"),
			false,
		),
	)
	updateCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"name",
			"n",
			abg.Trans("stacks.update.options.name.description"),
			"",
		),
	)
	updateCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"base",
			"b",
			abg.Trans("stacks.update.options.base.description"),
			"",
		),
	)
	updateCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"packages",
			"p",
			abg.Trans("stacks.update.options.packages.description"),
			"",
		),
	)
	updateCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"pkg-manager",
			"k",
			abg.Trans("stacks.update.options.pkgManager.description"),
			"",
		),
	)

	// Rm subcommand
	rmStackCmd := cmdr.NewCommand(
		"rm",
		abg.Trans("stacks.rm.description"),
		abg.Trans("stacks.rm.description"),
		removeStack,
	)

	rmStackCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"name",
			"n",
			abg.Trans("stacks.rm.options.name.description"),
			"",
		),
	)
	rmStackCmd.WithBoolFlag(
		cmdr.NewBoolFlag(
			"force",
			"f",
			abg.Trans("stacks.rm.options.force.description"),
			false,
		),
	)

	// Export subcommand
	exportCmd := cmdr.NewCommand(
		"export",
		abg.Trans("stacks.export.description"),
		abg.Trans("stacks.export.description"),
		exportStack,
	)
	exportCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"name",
			"n",
			abg.Trans("stacks.export.options.name.description"),
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
		abg.Trans("stacks.import.description"),
		abg.Trans("stacks.import.description"),
		importStack,
	)
	importCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"input",
			"i",
			abg.Trans("stacks.import.options.input.description"),
			"",
		),
	)

	// Add subcommands to stacks
	cmd.AddCommand(listCmd)
	cmd.AddCommand(showCmd)
	cmd.AddCommand(newCmd)
	cmd.AddCommand(updateCmd)
	cmd.AddCommand(rmStackCmd)
	cmd.AddCommand(exportCmd)
	cmd.AddCommand(importCmd)

	return cmd
}

func listStacks(cmd *cobra.Command, args []string) error {
	jsonFlag, _ := cmd.Flags().GetBool("json")

	stacks := core.ListStacks()

	if !jsonFlag {
		stacksCount := len(stacks)
		if stacksCount == 0 {
			fmt.Println(abg.Trans("stacks.list.info.noStacks"))
			return nil
		}

		cmdr.Info.Printfln(abg.Trans("stacks.list.info.foundStacks"), stacksCount)

		table := core.CreateApxTable(os.Stdout)
		table.SetHeader([]string{abg.Trans("stacks.labels.name"), "Base", abg.Trans("stacks.labels.builtIn"), "Pkgs", "Pkg manager"})

		for _, stack := range stacks {
			builtIn := abg.Trans("abg.terminal.no")
			if stack.BuiltIn {
				builtIn = abg.Trans("abg.terminal.yes")
			}
			table.Append([]string{stack.Name, stack.Base, builtIn, fmt.Sprintf("%d", len(stack.Packages)), stack.PkgManager})
		}

		table.Render()
	} else {
		jsonStacks, err := json.MarshalIndent(stacks, "", "  ")
		if err != nil {
			return err
		}

		fmt.Println(string(jsonStacks))
	}

	return nil
}

func showStack(cmd *cobra.Command, args []string) error {
	stack, error := core.LoadStack(args[0])
	if error != nil {
		return error
	}

	table := core.CreateApxTable(os.Stdout)
	table.Append([]string{abg.Trans("stacks.labels.name"), stack.Name})
	table.Append([]string{"Base", stack.Base})
	table.Append([]string{"Packages", strings.Join(stack.Packages, ", ")})
	table.Append([]string{"Package manager", stack.PkgManager})
	table.Render()

	return nil
}

func newStack(cmd *cobra.Command, args []string) error {
	noPrompt, _ := cmd.Flags().GetBool("no-prompt")
	name, _ := cmd.Flags().GetString("name")
	base, _ := cmd.Flags().GetString("base")
	packages, _ := cmd.Flags().GetString("packages")
	pkgManager, _ := cmd.Flags().GetString("pkg-manager")

	if name == "" {
		if !noPrompt {
			cmdr.Info.Println(abg.Trans("stacks.new.info.askName"))
			fmt.Scanln(&name)
			if name == "" {
				cmdr.Error.Println(abg.Trans("stacks.new.error.emptyName"))
				return nil
			}
		} else {
			cmdr.Error.Println(abg.Trans("stacks.new.error.noName"))
			return nil
		}
	}

	ok := core.StackExists(name)
	if ok {
		cmdr.Error.Printfln(abg.Trans("stacks.new.error.alreadyExists"), name)
		return nil
	}

	if base == "" {
		if !noPrompt {
			cmdr.Info.Println(abg.Trans("stacks.new.info.askBase"))
			fmt.Scanln(&base)
			if base == "" {
				cmdr.Error.Println(abg.Trans("stacks.new.error.emptyBase"))
				return nil
			}
		} else {
			cmdr.Error.Println(abg.Trans("stacks.new.error.noBase"))
			return nil
		}
	}

	if pkgManager == "" {
		pkgManagers := core.ListPkgManagers()
		if len(pkgManagers) == 0 {
			cmdr.Error.Println(abg.Trans("stacks.new.error.noPkgManagers"))
			return nil
		}

		cmdr.Info.Println(abg.Trans("stacks.new.info.askPkgManager"))
		for i, manager := range pkgManagers {
			fmt.Printf("%d. %s\n", i+1, manager.Name)
		}
		cmdr.Info.Printfln(abg.Trans("stacks.new.info.selectPkgManager"), len(pkgManagers))
		var pkgManagerIndex int
		_, err := fmt.Scanln(&pkgManagerIndex)
		if err != nil {
			cmdr.Error.Println(abg.Trans("abg.errors.invalidInput"))
			return nil
		}

		if pkgManagerIndex < 1 || pkgManagerIndex > len(pkgManagers) {
			cmdr.Error.Println(abg.Trans("abg.errors.invalidInput"))
			return nil
		}

		pkgManager = pkgManagers[pkgManagerIndex-1].Name
	}

	ok = core.PkgManagerExists(pkgManager)
	if !ok {
		cmdr.Error.Println(abg.Trans("stacks.new.error.pkgManagerDoesNotExist"))
		return nil
	}

	packagesArray := strings.Fields(packages)
	if len(packagesArray) == 0 && !noPrompt {
		cmdr.Info.Println(abg.Trans("stacks.new.info.noPackages") + "[y/N]")
		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(answer)
		if answer == "y" || answer == "Y" {
			cmdr.Info.Println(abg.Trans("stacks.new.info.askPackages"))
			packagesInput, _ := reader.ReadString('\n')
			packagesInput = strings.TrimSpace(packagesInput)
			packagesArray = strings.Fields(packagesInput)
		} else {
			packagesArray = []string{}
		}
	}

	stack := core.NewStack(name, base, packagesArray, pkgManager, false)

	err := stack.Save()
	if err != nil {
		return err
	}

	cmdr.Success.Printfln(abg.Trans("stacks.new.info.success"), name)

	return nil
}

func updateStack(cmd *cobra.Command, args []string) error {
	noPrompt, _ := cmd.Flags().GetBool("no-prompt")
	name, _ := cmd.Flags().GetString("name")
	base, _ := cmd.Flags().GetString("base")
	packages, _ := cmd.Flags().GetString("packages")
	pkgManager, _ := cmd.Flags().GetString("pkg-manager")

	if name == "" {
		if len(args) != 1 || args[0] == "" {
			cmdr.Error.Println(abg.Trans("stacks.update.error.noName"))
			return nil
		}

		cmd.Flags().Set("name", args[0])
		name = args[0]
	}

	stack, error := core.LoadStack(name)
	if error != nil {
		return error
	}

	if stack.BuiltIn {
		cmdr.Error.Println(abg.Trans("stacks.update.error.builtIn"))
		os.Exit(126)
	}

	if base == "" {
		if !noPrompt {
			cmdr.Info.Printfln(abg.Trans("stacks.update.info.askBase"), stack.Base)
			fmt.Scanln(&base)
			if base == "" {
				base = stack.Base
			}
		} else {
			cmdr.Error.Println(abg.Trans("stacks.update.error.noBase"))
			return nil
		}
	}

	if pkgManager == "" {
		if !noPrompt {
			cmdr.Info.Printfln(abg.Trans("stacks.update.info.askPkgManager"), stack.PkgManager)
			fmt.Scanln(&pkgManager)
			if pkgManager == "" {
				pkgManager = stack.PkgManager
			}
		} else {
			cmdr.Error.Println(abg.Trans("stacks.update.error.noPkgManager"))
			return nil
		}
	}

	ok := core.PkgManagerExists(pkgManager)
	if !ok {
		cmdr.Error.Println(abg.Trans("stacks.update.error.pkgManagerDoesNotExist"))
		return nil
	}

	if len(packages) > 0 {
		stack.Packages = strings.Fields(packages)
	} else if !noPrompt {
		if len(stack.Packages) > 0 {
			cmdr.Info.Println(abg.Trans("stacks.update.info.confirmPackages") + "[y/N]"  + "\n\t -", strings.Join(stack.Packages, "\n\t - "))
		} else {
			cmdr.Info.Println(abg.Trans("stacks.update.info.noPackages") + "[y/N]")
		}
		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(answer)

		packagesArray := []string{}

		if answer == "y" || answer == "Y" {
			cmdr.Info.Println(abg.Trans("stacks.update.info.askPackages"))
			packagesInput, _ := reader.ReadString('\n')
			packagesInput = strings.TrimSpace(packagesInput)
			packagesArray = strings.Fields(packagesInput)
			stack.Packages = packagesArray
		}
	}

	stack.Base = base
	stack.PkgManager = pkgManager

	err := stack.Save()
	if err != nil {
		return err
	}

	cmdr.Info.Printfln(abg.Trans("stacks.update.info.success"), name)

	return nil
}

func removeStack(cmd *cobra.Command, args []string) error {
	stackName, _ := cmd.Flags().GetString("name")
	if stackName == "" {
		cmdr.Error.Println(abg.Trans("stacks.rm.error.noName"))
		return nil
	}

	subSystems, _ := core.ListSubsystemForStack(stackName)
	if len(subSystems) > 0 {
		cmdr.Error.Printfln(abg.Trans("stacks.rm.error.inUse"), len(subSystems))
		table := core.CreateApxTable(os.Stdout)
		table.SetHeader([]string{abg.Trans("subsystems.labels.name"), "Stack", abg.Trans("subsystems.labels.status"), "Pkgs"})
		for _, subSystem := range subSystems {
			table.Append([]string{
				subSystem.Name,
				subSystem.Stack.Name,
				subSystem.Status,
				fmt.Sprintf("%d", len(subSystem.Stack.Packages)),
			})
		}
		table.Render()
		return nil
	}

	force, _ := cmd.Flags().GetBool("force")

	if !force {
		reader := bufio.NewReader(os.Stdin)
		validChoice := false
		for !validChoice {
			cmdr.Info.Printfln(abg.Trans("stacks.rm.info.askConfirmation")+` [y/N]`, stackName)
			answer, _ := reader.ReadString('\n')
			if answer == "\n" {
				answer = "n\n"
			}
			answer = strings.ToLower(strings.ReplaceAll(answer, " ", ""))
			switch answer {
			case "y\n":
				validChoice = true
				force = true
			case "n\n":
				validChoice = true
			default:
				cmdr.Warning.Println(abg.Trans("abg.errors.invalidChoice"))
			}
		}
	}

	if !force {
		cmdr.Info.Printfln(abg.Trans("pkgmanagers.rm.info.aborting"), stackName)
		return nil
	}

	stack, error := core.LoadStack(stackName)
	if error != nil {
		return error
	}

	error = stack.Remove()
	if error != nil {
		return error
	}

	cmdr.Info.Printfln(abg.Trans("stacks.rm.info.success"), stackName)
	return nil
}

func exportStack(cmd *cobra.Command, args []string) error {
	stackName, _ := cmd.Flags().GetString("name")
	if stackName == "" {
		cmdr.Error.Println(abg.Trans("stacks.export.error.noName"))
		return nil
	}

	stack, error := core.LoadStack(stackName)
	if error != nil {
		return error
	}

	output, _ := cmd.Flags().GetString("output")
	if output == "" {
		cmdr.Error.Println(abg.Trans("stacks.export.error.noOutput"))
		return nil
	}

	error = stack.Export(output)
	if error != nil {
		return error
	}

	cmdr.Info.Printfln(abg.Trans("stacks.export.info.success"), stack.Name, output)
	return nil
}

func importStack(cmd *cobra.Command, args []string) error {
	input, _ := cmd.Flags().GetString("input")
	if input == "" {
		cmdr.Error.Println(abg.Trans("stacks.import.error.noInput"))
		return nil
	}

	stack, error := core.LoadStackFromPath(input)
	if error != nil {
		cmdr.Error.Printf(abg.Trans("stacks.import.error.cannotLoad"), input)
	}

	error = stack.Save()
	if error != nil {
		return error
	}

	cmdr.Info.Printfln(abg.Trans("stacks.import.info.success"), stack.Name)
	return nil
}
