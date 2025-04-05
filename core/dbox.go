package core

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

type DBox struct {
	Engine       string
	EngineBinary string
	Version      string
}

type DBoxContainer struct {
	ID        string
	CreatedAt string
	Status    string
	Labels    map[string]string
	Name      string
}

func NewDBox() (*DBox, error) {
	engineBinary, engine := getEngine()

	version, err := getDBoxVersion()
	if err != nil {
		return nil, err
	}

	return &DBox{
		Engine:       engine,
		EngineBinary: engineBinary,
		Version:      version,
	}, nil
}

func getEngine() (string, string) {
	if podmanBinary, err := exec.LookPath("podman"); err == nil {
		return podmanBinary, "podman"
	}
	if dockerBinary, err := exec.LookPath("docker"); err == nil {
		return dockerBinary, "docker"
	}
	log.Fatal("no container engine found. Please install Podman or Docker.")
	return "", ""
}

func getDBoxVersion() (string, error) {
	output, err := exec.Command(abg.Cnf.DistroboxPath, "version").Output()
	if err != nil {
		return "", err
	}

	parts := strings.Split(string(output), "distrobox: ")
	if len(parts) != 2 {
		return "", errors.New("can't retrieve distrobox version")
	}

	return strings.TrimSpace(parts[1]), nil
}

func (d *DBox) RunCommand(command string, args, engineFlags []string, useEngine, captureOutput, muteOutput, rootFull, detached bool) ([]byte, error) {
	entrypoint := abg.Cnf.DistroboxPath
	finalArgs := []string{command}

	if useEngine {
		entrypoint = d.EngineBinary
	}

	if rootFull && useEngine {
		entrypoint = "pkexec"
		finalArgs = []string{d.EngineBinary, command}
	}

	cmd := exec.Command(entrypoint, finalArgs...)

	if detached {
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	}

	if !captureOutput && !muteOutput {
		cmd.Stdout = os.Stdout
	}
	if !muteOutput {
		cmd.Stderr = os.Stderr
	}
	cmd.Stdin = os.Stdin

	cmd.Env = append(os.Environ(), "DBX_SUDO_PROGRAM=pkexec")

	switch d.Engine {
	case "podman":
		cmd.Env = append(cmd.Env, "CONTAINER_STORAGE_DRIVER="+abg.Cnf.StorageDriver)
	case "docker":
		cmd.Env = append(cmd.Env, "DOCKER_STORAGE_DRIVER="+abg.Cnf.StorageDriver)
	}

	if len(engineFlags) > 0 {
		cmd.Args = append(cmd.Args, "--additional-flags", strings.Join(engineFlags, " "))
	}
	if rootFull && !useEngine {
		cmd.Args = append(cmd.Args, "--root")
	}
	cmd.Args = append(cmd.Args, args...)

	if os.Getenv("ABG_VERBOSE") == "1" {
		fmt.Println("Running a command:")
		fmt.Printf("\tCommand: %s\n", cmd.String())
		fmt.Printf("\tcaptureOutput: %v\n\tmuteOutput: %v\n\trootFull: %v\n\tdetachedMode: %v\n",
			captureOutput, muteOutput, rootFull, detached)
	}

	if detached {
		return nil, cmd.Start()
	}

	if captureOutput {
		output, err := cmd.Output()
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				return output, errors.New(string(exitErr.Stderr))
			}
		}
		return output, err
	}

	return nil, cmd.Run()
}

func (d *DBox) ListContainers(rootFull bool) ([]DBoxContainer, error) {
	output, err := d.RunCommand("ps", []string{
		"-a",
		"--format", "{{.ID}}|{{.CreatedAt}}|{{.Status}}|{{.Labels}}|{{.Names}}",
	}, nil, true, true, true, rootFull, false)
	if err != nil {
		return nil, err
	}

	var containers []DBoxContainer
	for _, row := range strings.Split(string(output), "\n") {
		if row == "" {
			continue
		}

		parts := strings.Split(row, "|")
		if len(parts) != 5 {
			continue
		}

		container := DBoxContainer{
			ID:        strings.Trim(parts[0], "\""),
			CreatedAt: strings.Trim(parts[1], "\""),
			Status:    strings.Trim(parts[2], "\""),
			Name:      strings.Trim(parts[4], "\""),
			Labels:    parseLabels(parts[3]),
		}
		containers = append(containers, container)
	}

	return containers, nil
}

func parseLabels(raw string) map[string]string {
	labels := make(map[string]string)
	raw = strings.TrimPrefix(strings.TrimSuffix(raw, "]"), "map[")
	for _, item := range strings.Fields(raw) {
		kv := strings.SplitN(item, ":", 2)
		if len(kv) == 2 {
			labels[kv[0]] = kv[1]
		}
	}
	return labels
}

func (d *DBox) GetContainer(name string, rootFull bool) (*DBoxContainer, error) {
	containers, err := d.ListContainers(rootFull)
	if err != nil {
		return nil, err
	}
	for _, c := range containers {
		if c.Name == name {
			return &c, nil
		}
	}
	return nil, errors.New("container not found")
}

func (d *DBox) ContainerDelete(name string, rootFull bool) error {
	_, err := d.RunCommand("rm", []string{"--force", name}, nil, false, false, true, rootFull, false)
	return err
}

func (d *DBox) CreateContainer(name, image string, packages []string, home string, labels map[string]string, withInit, rootFull, unshared, withNvidia bool, hostname string) error {
	args := []string{
		"--image", image,
		"--name", name,
		"--no-entry",
		"--yes",
		"--pull",
	}

	if home != "" {
		args = append(args, "--home", home)
	}
	if hasNvidiaGPU() && withNvidia {
		args = append(args, "--nvidia")
	}
	if withInit {
		args = append(args, "--init")
	}
	if unshared {
		args = append(args, "--unshare-all")
	}
	if hostname != "" {
		args = append(args, "--hostname", hostname)
	}
	if len(packages) > 0 {
		args = append(args, "--additional-packages", strings.Join(packages, " "))
	}

	var engineFlags []string
	for k, v := range labels {
		engineFlags = append(engineFlags, fmt.Sprintf("--label=%s=%s", k, v))
	}
	engineFlags = append(engineFlags, "--label=manager=abg")

	_, err := d.RunCommand("create", args, engineFlags, false, false, false, rootFull, false)
	return err
}

func (d *DBox) RunContainerCommand(name string, command []string, rootFull, detached bool) error {
	args := append([]string{"--name", name, "--"}, command...)
	_, err := d.RunCommand("run", args, nil, false, false, false, rootFull, detached)
	return err
}

func (d *DBox) ContainerExec(name string, captureOutput, muteOutput, rootFull, detached bool, args ...string) (string, error) {
	fullArgs := append([]string{name, "--"}, args...)
	out, err := d.RunCommand("enter", fullArgs, nil, false, captureOutput, muteOutput, rootFull, detached)
	if err != nil && err.Error() == "exit status 130" {
		return string(out), nil
	}
	return string(out), err
}

func (d *DBox) ContainerEnter(name string, rootFull bool) error {
	_, err := d.RunCommand("enter", []string{name}, nil, false, false, false, rootFull, false)
	if err != nil && err.Error() == "exit status 130" {
		return nil
	}
	return err
}

func (d *DBox) ContainerStart(name string, rootFull bool) error {
	_, err := d.RunCommand("start", []string{name}, nil, true, false, false, rootFull, false)
	return err
}

func (d *DBox) ContainerStop(name string, rootFull bool) error {
	_, err := d.RunCommand("stop", []string{name, "--yes"}, nil, false, false, false, rootFull, false)
	return err
}

func (d *DBox) ContainerExport(name string, delete, rootFull bool, args ...string) error {
	finalArgs := append([]string{"distrobox-export"}, args...)
	if delete {
		finalArgs = append([]string{"--delete"}, finalArgs...)
	}
	_, err := d.ContainerExec(name, true, true, rootFull, false, finalArgs...)
	return err
}

func (d *DBox) ContainerExportDesktopEntry(name, app, label string, rootFull bool) error {
	return d.ContainerExport(name, false, rootFull, "--app", app, "--export-label", label)
}

func (d *DBox) ContainerUnexportDesktopEntry(name, app string, rootFull bool) error {
	return d.ContainerExport(name, true, rootFull, "--app", app)
}

func (d *DBox) ContainerExportBin(name, binary, path string, rootFull bool) error {
	return d.ContainerExport(name, false, rootFull, "--bin", binary, "--export-path", path)
}

func (d *DBox) ContainerUnexportBin(name, binary string, rootFull bool) error {
	return d.ContainerExport(name, true, rootFull, "--bin", binary)
}
