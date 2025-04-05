package core

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func (a *Abg) EssentialChecks() error {
	if err := a.CheckContainerTools(); err != nil {
		fmt.Println(`One or more core components are not available.
Please refer to our documentation at https://documentation.auruos.org/`)
		return err
	}

	if err := a.EnsureDirectory(a.Cnf.UserStacksPath, "user stacks"); err != nil {
		return err
	}

	if err := a.EnsureDirectory(a.Cnf.AbgStoragePath, "ABG storage"); err != nil {
		return err
	}

	if err := a.EnsureDirectory(a.Cnf.UserPkgManagersPath, "user package managers"); err != nil {
		return err
	}

	return nil
}

func (a *Abg) CheckContainerTools() error {
	if _, err := os.Stat(a.Cnf.DistroboxPath); err != nil {
		if os.IsNotExist(err) {
			return errors.New("distrobox is not installed")
		}
		return err
	}

	if _, err := exec.LookPath("docker"); err != nil {
		if _, err := exec.LookPath("podman"); err != nil {
			return errors.New("no container engine (docker or podman) found")
		}
	}

	return nil
}

func IsOverlayTypeFS() bool {
	output, err := exec.Command("df", "-T", "/").Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), "overlay")
}

func ExitIfOverlayTypeFS() {
	if IsOverlayTypeFS() {
		log.Println("ABG does not work with overlay-type filesystem.")
		os.Exit(1)
	}
}

func (a *Abg) EnsureDirectory(path, label string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			if mkErr := os.MkdirAll(path, 0755); mkErr != nil {
				return fmt.Errorf("failed to create %s directory: %w", label, mkErr)
			}
		} else {
			return fmt.Errorf("failed to check %s directory: %w", label, err)
		}
	}
	return nil
}

func hasNvidiaGPU() bool {
	_, err := os.Stat("/dev/nvidia0")
	return err == nil
}
