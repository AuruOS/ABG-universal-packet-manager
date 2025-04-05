package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/AuruOS/abg/core"
	"github.com/AuruOS/orchid/cmdr"
)

func NewSubSystemsCommand() *cmdr.Command {
	// Root command
	cmd := cmdr.NewCommand(
		"subsystems",
		abg.Trans("subsystems.description"),
		abg.Trans("subsystems.description"),
		nil,
	)

	// List subcommand
	listCmd := cmdr.NewCommand(
		"list",
		abg.Trans("subsystems.list.description"),
		abg.Trans("subsystems.list.description"),
		listSubSystems,
	)

	listCmd.WithBoolFlag(
		cmdr.NewBoolFlag(
			"json",
			"j",
			abg.Trans("subsystems.list.options.json.description"),
			false,
		),
	)

	// New subcommand
	newCmd := cmdr.NewCommand(
		"new",
		abg.Trans("subsystems.new.description"),
		abg.Trans("subsystems.new.description"),
		newSubSystem,
	)

	newCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"stack",
			"s",
			abg.Trans("subsystems.new.options.stack.description"),
			"",
		),
	)
	newCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"name",
			"n",
			abg.Trans("subsystems.new.options.name.description"),
			"",
		),
	)
	newCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"home",
			"H",
			abg.Trans("subsystems.new.options.home.description"),
			"",
		),
	)
	newCmd.WithBoolFlag(
		cmdr.NewBoolFlag(
			"init",
			"i",
			abg.Trans("subsystems.new.options.init.description"),
			false,
		),
	)

	// Rm subcommand
	rmCmd := cmdr.NewCommand(
		"rm",
		abg.Trans("subsystems.rm.description"),
		abg.Trans("subsystems.rm.description"),
		rmSubSystem,
	)

	rmCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"name",
			"n",
			abg.Trans("subsystems.rm.options.name.description"),
			"",
		),
	)
	rmCmd.WithBoolFlag(
		cmdr.NewBoolFlag(
			"force",
			"f",
			abg.Trans("subsystems.rm.options.force.description"),
			false,
		),
	)

	// Reset subcommand
	resetCmd := cmdr.NewCommand(
		"reset",
		abg.Trans("subsystems.reset.description"),
		abg.Trans("subsystems.reset.description"),
		resetSubSystem,
	)

	resetCmd.WithStringFlag(
		cmdr.NewStringFlag(
			"name",
			"n",
			abg.Trans("subsystems.reset.options.name.description"),
			"",
		),
	)
	resetCmd.WithBoolFlag(
		cmdr.NewBoolFlag(
			"force",
			"f",
			abg.Trans("subsystems.reset.options.force.description"),
			false,
		),
	)

	// Add subcommands to subsystems
	cmd.AddCommand(listCmd)
	cmd.AddCommand(newCmd)
	cmd.AddCommand(rmCmd)
	cmd.AddCommand(resetCmd)

	return cmd
}

func listSubSystems(cmd *cobra.Command, args []string) error {
	jsonFlag, _ := cmd.Flags().GetBool("json")

	subSystems, err := core.ListSubSystems(false, false)
	if err != nil {
		return err
	}

	if !jsonFlag {
		subSystemsCount := len(subSystems)
		if subSystemsCount == 0 {
			cmdr.Info.Println(abg.Trans("subsystems.list.info.noSubsystems"))
			return nil
		}

		cmdr.Info.Printfln(abg.Trans("subsystems.list.info.foundSubsystems"), subSystemsCount)

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
	} else {
		jsonSubSystems, err := json.MarshalIndent(subSystems, "", "  ")
		if err != nil {
			return err
		}

		fmt.Println(string(jsonSubSystems))
	}

	return nil
}

func newSubSystem(cmd *cobra.Command, args []string) error {
	home, _ := cmd.Flags().GetString("home")
	stackName, _ := cmd.Flags().GetString("stack")
	subSystemName, _ := cmd.Flags().GetString("name")
	isInit, _ := cmd.Flags().GetBool("init")

	stacks := core.ListStacks()
	if len(stacks) == 0 {
		cmdr.Error.Println(abg.Trans("subsystems.new.error.noStacks"))
		return nil
	}

	if subSystemName == "" {
		cmdr.Info.Println(abg.Trans("subsystems.new.info.askName"))
		fmt.Scanln(&subSystemName)
		if subSystemName == "" {
			cmdr.Error.Println(abg.Trans("subsystems.new.error.emptyName"))
			return nil
		}
	}

	if stackName == "" {
		cmdr.Info.Println(abg.Trans("subsystems.new.info.availableStacks"))
		for i, stack := range stacks {
			fmt.Printf("%d. %s\n", i+1, stack.Name)
		}
		cmdr.Info.Printfln(abg.Trans("subsystems.new.info.selectStack"), len(stacks))

		var stackIndex int
		_, err := fmt.Scanln(&stackIndex)
		if err != nil {
			cmdr.Error.Println(abg.Trans("abg.errors.invalidInput"))
			return nil
		}

		if stackIndex < 1 || stackIndex > len(stacks) {
			cmdr.Error.Println(abg.Trans("abg.errors.invalidInput"))
			return nil
		}

		stackName = stacks[stackIndex-1].Name
	}

	checkSubSystem, err := core.LoadSubSystem(subSystemName, false)
	if err == nil {
		cmdr.Error.Printf(abg.Trans("subsystems.new.error.alreadyExists"), checkSubSystem.Name)
		return nil
	}

	for _, existcommand := range cmd.Root().Commands() {
		if subSystemName == existcommand.Name() {
			cmdr.Error.Printfln(abg.Trans("subsystems.new.error.forbiddenName"), subSystemName)
			return nil
		}
	}

	stack, err := core.LoadStack(stackName)
	if err != nil {
		return err
	}

	subSystem, err := core.NewSubSystem(subSystemName, stack, home, isInit, false, false, false, true, "")
	if err != nil {
		return err
	}

	spinner, _ := cmdr.Spinner.Start(fmt.Sprintf(abg.Trans("subsystems.new.info.creatingSubsystem"), subSystemName, stackName))
	err = subSystem.Create()
	if err != nil {
		return err
	}

	spinner.UpdateText(fmt.Sprintf(abg.Trans("subsystems.new.info.success"), subSystemName))
	spinner.Success()

	return nil
}

func rmSubSystem(cmd *cobra.Command, args []string) error {
	subSystemName, _ := cmd.Flags().GetString("name")
	forceFlag, _ := cmd.Flags().GetBool("force")

	if subSystemName == "" {
		cmdr.Error.Println(abg.Trans("subsystems.rm.error.noName"))
		return nil
	}

	if !forceFlag {
		cmdr.Info.Printfln(abg.Trans("subsystems.rm.info.askConfirmation")+` [y/N]`, subSystemName)
		var confirmation string
		fmt.Scanln(&confirmation)
		if strings.ToLower(confirmation) != "y" {
			cmdr.Info.Println(abg.Trans("abg.info.aborting"))
			return nil
		}
	}

	subSystem, err := core.LoadSubSystem(subSystemName, false)
	if err != nil {
		return err
	}

	err = subSystem.Remove()
	if err != nil {
		return err
	}

	cmdr.Success.Printfln(abg.Trans("subsystems.rm.info.success"), subSystemName)

	return nil
}

func resetSubSystem(cmd *cobra.Command, args []string) error {
	subSystemName, _ := cmd.Flags().GetString("name")
	if subSystemName == "" {
		cmdr.Error.Println(abg.Trans("subsystems.reset.error.noName"))
		return nil
	}

	forceFlag, _ := cmd.Flags().GetBool("force")

	if !forceFlag {
		cmdr.Info.Printfln(abg.Trans("subsystems.reset.info.askConfirmation")+` [y/N]`, subSystemName)
		var confirmation string
		fmt.Scanln(&confirmation)
		if strings.ToLower(confirmation) != "y" {
			cmdr.Info.Println(abg.Trans("abg.info.aborting"))
			return nil
		}
	}

	subSystem, err := core.LoadSubSystem(subSystemName, false)
	if err != nil {
		return err
	}

	err = subSystem.Reset()
	if err != nil {
		return err
	}

	cmdr.Success.Printfln(abg.Trans("subsystems.reset.info.success"), subSystemName)

	return nil
}
