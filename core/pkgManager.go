package core

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// PkgManager represents a package manager in ABG.
type PkgManager struct {
	// Model defines the command model:
	// 1: name + command + args (deprecated)
	// 2: full command string (recommended)
	Model         int
	Name          string
	NeedSudo      bool
	CmdAutoRemove string
	CmdClean      string
	CmdInstall    string
	CmdList       string
	CmdPurge      string
	CmdRemove     string
	CmdSearch     string
	CmdShow       string
	CmdUpdate     string
	CmdUpgrade    string
	BuiltIn       bool // Built-in managers can't be removed
}

// NewPkgManager creates a new instance of PkgManager.
func NewPkgManager(name string, needSudo bool, autoRemove, clean, install, list, purge, remove, search, show, update, upgrade string, builtIn bool) *PkgManager {
	return &PkgManager{
		Name:          name,
		NeedSudo:      needSudo,
		CmdAutoRemove: autoRemove,
		CmdClean:      clean,
		CmdInstall:    install,
		CmdList:       list,
		CmdPurge:      purge,
		CmdRemove:     remove,
		CmdSearch:     search,
		CmdShow:       show,
		CmdUpdate:     update,
		CmdUpgrade:    upgrade,
		BuiltIn:       builtIn,
		Model:         2,
	}
}

// LoadPkgManager attempts to load a user-defined or built-in package manager.
func LoadPkgManager(name string) (*PkgManager, error) {
	userFile := SelectYamlFile(abg.Cnf.UserPkgManagersPath, name)
	pm, err := loadPkgManagerFromPath(userFile)
	if err == nil {
		return pm, nil
	}

	sharedFile := SelectYamlFile(abg.Cnf.PkgManagersPath, name)
	return loadPkgManagerFromPath(sharedFile)
}

// Save persists the PkgManager to user storage.
func (pm *PkgManager) Save() error {
	filePath := SelectYamlFile(abg.Cnf.UserPkgManagersPath, pm.Name)
	data, err := yaml.Marshal(pm)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}

// Remove deletes the user-defined package manager.
func (pm *PkgManager) Remove() error {
	if pm.BuiltIn {
		return errors.New("cannot remove built-in package manager")
	}
	filePath := SelectYamlFile(abg.Cnf.UserPkgManagersPath, pm.Name)
	return os.Remove(filePath)
}

// Export writes the package manager definition to a custom location.
func (pm *PkgManager) Export(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err = os.MkdirAll(path, 0755); err != nil {
			return err
		}
	}
	filePath := SelectYamlFile(path, pm.Name)
	data, err := yaml.Marshal(pm)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}

// GenCmd builds the full command for the container environment.
func (pm *PkgManager) GenCmd(cmd string, args ...string) []string {
	var finalArgs []string

	if pm.NeedSudo {
		finalArgs = append(finalArgs, "sudo")
	}

	if pm.Model == 0 || pm.Model == 1 {
		fmt.Println("!!! DEPRECATION WARNING: Model 1 is deprecated. Please update your ABG package manager.")
		finalArgs = append(finalArgs, pm.Name, cmd)
	} else {
		finalArgs = append(finalArgs, strings.Fields(cmd)...)
	}

	return append(finalArgs, args...)
}

// ListPkgManagers lists all available package managers.
func ListPkgManagers() []*PkgManager {
	var managers []*PkgManager
	managers = append(managers, listPkgManagersFromPath(abg.Cnf.UserPkgManagersPath)...)

	if abg.Cnf.PkgManagersPath != abg.Cnf.UserPkgManagersPath {
		managers = append(managers, listPkgManagersFromPath(abg.Cnf.PkgManagersPath)...)
	}
	return managers
}

// listPkgManagersFromPath scans for package managers in a given path.
func listPkgManagersFromPath(path string) []*PkgManager {
	var managers []*PkgManager

	files, err := os.ReadDir(path)
	if err != nil {
		return managers
	}

	for _, f := range files {
		ext := filepath.Ext(f.Name())
		if !f.IsDir() && (ext == ".yaml" || ext == ".yml") {
			name := f.Name()[:len(f.Name())-len(ext)]
			pm, err := LoadPkgManager(name)
			if err == nil {
				managers = append(managers, pm)
			}
		}
	}
	return managers
}

// PkgManagerExists checks if a manager by name is present.
func PkgManagerExists(name string) bool {
	_, err := LoadPkgManager(name)
	return err == nil
}

// LoadPkgManagerFromPath loads a manager from a direct file path.
func LoadPkgManagerFromPath(path string) (*PkgManager, error) {
	return loadPkgManagerFromPath(path)
}

func loadPkgManagerFromPath(path string) (*PkgManager, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.New("package manager not found")
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	pm := &PkgManager{}
	if err := yaml.Unmarshal(data, pm); err != nil {
		return nil, err
	}

	// Backward compatibility for old configs
	if pm.Model == 0 {
		pm.Model = 1
	}
	return pm, nil
}
