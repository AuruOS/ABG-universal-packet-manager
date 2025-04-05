package core

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/olekukonko/tablewriter"
)

var ProcessPath string

// RootCheck checks if the current user has root privileges.
func RootCheck(display bool) bool {
	if os.Geteuid() != 0 {
		if display {
			fmt.Println("You must be root to run this command")
		}
		return false
	}
	return true
}

// AskConfirmation prompts the user for a yes/no confirmation.
func AskConfirmation(prompt string) bool {
	var response string
	fmt.Printf("%s [y/N]: ", prompt)
	fmt.Scanln(&response)
	return response == "y" || response == "Y"
}

// CopyToUserTemp copies a file to the user's temporary cache directory.
func CopyToUserTemp(path string) (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}

	cacheDir := filepath.Join(user.HomeDir, ".cache", "abg")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", err
	}

	fileName := filepath.Base(path)
	newPath := filepath.Join(cacheDir, fileName)

	if err := copyFile(path, newPath); err != nil {
		return "", err
	}

	return newPath, nil
}

// getPrettifiedDate returns a human-readable date from a timestamp.
func getPrettifiedDate(date string) string {
	t, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", date)
	if err != nil {
		return date // Return original date if parsing fails.
	}

	if t.After(time.Now().Add(-24 * time.Hour)) {
		duration := time.Since(t).Round(time.Hour)
		hours := int(duration.Hours())
		return fmt.Sprintf("%dh ago", hours)
	}

	return t.Format("02 Jan 2006 15:04:05")
}

// CreateApxTable initializes a table writer for displaying data.
func CreateApxTable(writer io.Writer) *tablewriter.Table {
	table := tablewriter.NewWriter(writer)
	table.SetColumnSeparator("┊")
	table.SetCenterSeparator("┼")
	table.SetRowSeparator("┄")
	table.SetHeaderLine(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetRowLine(true)

	return table
}

// CopyFile copies a file from src to dst.
func CopyFile(src, dst string) error {
	return copyFile(src, dst)
}

// SelectYamlFile selects a YAML file based on the provided base path and name.
func SelectYamlFile(basePath string, name string) string {
	const (
		YML  = ".yml"
		YAML = ".yaml"
	)

	yamlFile := filepath.Join(basePath, fmt.Sprintf("%s%s", name, YAML))
	ymlFile := filepath.Join(basePath, fmt.Sprintf("%s%s", name, YML))

	if _, err := os.Stat(yamlFile); errors.Is(err, os.ErrNotExist) {
		return ymlFile // Return .yml if .yaml does not exist.
	}

	return yamlFile // Return .yaml if it exists.
}

// copyFile is a helper function that performs the actual file copying.
func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	if _, err = io.Copy(destination, source); err != nil {
		return err
	}

	return nil
}
