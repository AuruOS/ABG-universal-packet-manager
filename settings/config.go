package settings

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	// Paths
	AbgPath       string `json:"abgPath"`
	DistroboxPath string `json:"distroboxPath"`
	StorageDriver string `json:"storageDriver"`

	// Virtual
	UserAbgPath         string
	AbgStoragePath      string
	StacksPath          string
	UserStacksPath      string
	PkgManagersPath     string
	UserPkgManagersPath string
}

func GetAbgDefaultConfig() (*Config, error) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	// dev paths
	viper.AddConfigPath("config/")

	// tests paths
	viper.AddConfigPath("../config/")

	// user paths
	viper.AddConfigPath(filepath.Join(userHome, ".config/abg/"))

	// prod paths
	viper.AddConfigPath("/etc/abg/")
	viper.AddConfigPath("/usr/share/abg/")

	viper.SetConfigName("abg")
	viper.SetConfigType("json")

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Unable to read config file: \n\t%s\n", err)
		os.Exit(1)
	}

	distroboxPath := viper.GetString("distroboxPath")

	_, err = os.Stat(distroboxPath)
	if err != nil {
		if os.IsNotExist(err) {
			path, err := exec.LookPath("distrobox")
			if err != nil {
				fmt.Printf("Unable to find distrobox in PATH.\n")
			} else {
				distroboxPath = path
			}
		}
	}

	Cnf := NewAbgConfig(
		viper.GetString("abgPath"),
		distroboxPath,
		viper.GetString("storageDriver"),
	)
	return Cnf, nil
}

func NewAbgConfig(abgPath, distroboxPath, storageDriver string) *Config {
	userHome, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	Cnf := &Config{
		// Common
		AbgPath:       abgPath,
		DistroboxPath: distroboxPath,
		StorageDriver: storageDriver,

		// Virtual
		UserAbgPath:         "",
		AbgStoragePath:      "",
		StacksPath:          "",
		UserStacksPath:      "",
		PkgManagersPath:     "",
		UserPkgManagersPath: "",
	}

	Cnf.UserAbgPath = filepath.Join(userHome, ".local/share/abg")
	Cnf.AbgStoragePath = filepath.Join(Cnf.UserAbgPath, "storage")
	Cnf.StacksPath = filepath.Join(Cnf.AbgPath, "stacks")
	Cnf.UserStacksPath = filepath.Join(Cnf.UserAbgPath, "stacks")
	Cnf.PkgManagersPath = filepath.Join(Cnf.AbgPath, "package-managers")
	Cnf.UserPkgManagersPath = filepath.Join(Cnf.UserAbgPath, "package-managers")

	return Cnf
}
