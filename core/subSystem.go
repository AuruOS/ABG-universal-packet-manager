
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type SubSystem struct {
	InternalName         string
	Name                 string
	Stack                *Stack
	Status               string
	HasInit              bool
	IsManaged            bool
	IsRootfull           bool
	IsUnshared           bool
	HasNvidiaIntegration bool
	ExportedPrograms     map[string]map[string]string
}

func findExported(internalName string, name string) map[string]map[string]string {
	bins := findExportedBinaries(internalName)
	progs := findExportedPrograms(internalName, name)

	for k, v := range progs {
		bins[k] = v // Give priority to application if duplicate is found.
	}

	return bins
}

func (s *SubSystem) Create() error {
	dbox, err := NewDbox()
	if err != nil {
		return err
	}

	labels := map[string]string{
		"stack": strings.ReplaceAll(s.Stack.Name, " ", "\\ "),
		"name":  strings.ReplaceAll(s.Name, " ", "\\ "),
	}

	if s.IsManaged {
		labels["managed"] = "true"
	}
	if s.HasInit {
		labels["hasInit"] = "true"
	}
	if s.IsUnshared {
		labels["unshared"] = "true"
	}
	if s.HasNvidiaIntegration {
		labels["nvidia"] = "true"
	}

	return dbox.CreateContainer(
		s.InternalName,
		s.Stack.Base,
		s.Stack.Packages,
		s.Home,
		labels,
		s.HasInit,
		s.IsRootfull,
		s.IsUnshared,
		s.HasNvidiaIntegration,
		s.Hostname,
    )
}

func LoadSubSystem(name string, isRootFull bool) (*SubSystem, error) {
	dbox, err := NewDbox()
	if err != nil {
        return nil, err
    }

	internalName := genInternalName(name)
	container, err := dbox.GetContainer(internalName, isRootFull)
	if err != nil {
        return nil, err
    }

	stack, err := LoadStack(container.Labels["stack"])
	if err != nil {
        return nil, err
    }

	return &SubSystem{
        InternalName: internalName,
        Name:         container.Labels["name"],
        Stack:        stack,
        Status:       container.Status,
        HasInit:      container.Labels["hasInit"] == "true",
        IsManaged:    container.Labels["managed"] == "true",
        IsRootfull:   isRootFull,
        IsUnshared:   container.Labels["unshared"] == "true",
    }, nil
}

func ListSubSystems(includeManaged bool, includeRootFull bool) ([]*SubSystem, error) {
	dbox, err := NewDbox()
	if err != nil {
        return nil, err
    }

    containers, err := dbox.ListContainers(includeRootFull)
    if err != nil {
        return nil, err
    }

	subsystems := make([]*SubSystem, 0)
	for _, container := range containers {
	    if _, ok := container.Labels["name"]; !ok {
	        continue // Skip containers without a name label.
	    }

	    if !includeManaged && container.Labels["managed"] == "true" {
	        continue // Skip managed containers if not included.
	    }

	    stack, err := LoadStack(container.Labels["stack"])
	    if err != nil {
	        log.Printf("Error loading stack %s: %s", container.Labels["stack"], err)
	        continue
	    }

	    internalName := genInternalName(container.Labels["name"])
	    subsystem := &SubSystem{
	        InternalName:     internalName,
	        Name:             container.Labels["name"],
	        Stack:            stack,
	        Status:           container.Status,
	        ExportedPrograms: findExported(internalName, container.Labels["name"]),
	    }

	    subsystems = append(subsystems, subsystem)
    }

	return subsystems, nil
}

// ListSubsystemForStack returns a list of subsystems for the specified stack.
func ListSubsystemForStack(stackName string) ([]*SubSystem, error) {
	dbox, err := NewDbox()
	if err != nil {
        return nil, err
    }

	rootlessContainers, err := dbox.ListContainers(false)
	if err != nil {
        return nil, err
    }

	rootfullContainers, err := dbox.ListContainers(true)
	if err != nil {
        return nil, err
    }

	var containers []Container // Assuming Container is a defined type.
	for _, c := range rootlessContainers {
	    containers = append(containers, c)
    }
	for _, c := range rootfullContainers {
	    containers = append(containers, c)
    }

	subsystems := make([]*SubSystem, 0)
	for _, container := range containers {
	    if _, ok := container.Labels["name"]; !ok {
	        continue // Skip containers without a name label.
	    }

	    stack, err := LoadStack(stackName)
	    if err != nil {
	        log.Printf("Error loading stack %s: %s", stackName, err)
	        continue
	    }

	    internalName := genInternalName(container.Labels["name"])
	    subsystem := &SubSystem{
	        InternalName:     internalName,
	        Name:             container.Labels["name"],
	        Stack:            stack,
	        Status:           container.Status,
	        ExportedPrograms: findExported(internalName, container.Labels["name"]),
	    }

	    if subsystem.Stack.Name == stack.Name { // Check for matching stack names.
	        subsystems = append(subsystems, subsystem)
	    }
    }

	return subsystems, nil
}

// Exec executes a command in the subsystem.
func (s *SubSystem) Exec(captureOutput bool, detachedMode bool, args ...string) (string, error) {
	dbox, err := NewDbox()
	if err != nil {
        return "", err
    }

	outStrg ,err:= dbox.ContainerExec(s.InternalName,captureOutput,false,s.IsRootfull ,detachedMode,args...)

	if captureOutput{
	  return outStrg,nil
	  }

	return "",nil
}

// Enter enters the subsystem's environment.
func (s *SubSystem) Enter() error {
	dbox ,err:= NewDbox()
	if(err!=nil){
	  return 	err
	  }
	  return dbox.ContainerEnter(s.InternalName,s.IsRootfull)
}

// Start starts the subsystem.
func (s *SubSystem) Start() error {
	dbox ,err:= NewDbox()
	if(err!=nil){
	  return 	err
	  }
	  return dbox.ContainerStart(s.InternalName,s.IsRootfull)
}

// Stop stops the subsystem.
func (s *SubSystem) Stop() error {
	dbox ,err:= NewDbox()
	if(err!=nil){
	  return 	err
	  }
	  return dbox.ContainerStop(s.InternalName,s.IsRootfull)
}

// Remove deletes the subsystem.
func (s *SubSystem) Remove() error {
	dbox ,err:= NewDbox()
	if(err!=nil){
	  return 	err
	  }

	return dbox.ContainerDelete(s.InternalName,s.IsRootfull)
}

// Reset removes and recreates the subsystem.
func (s *SubSystem) Reset() error {
	err:= s.Remove()
	if(err!=nil){
	   return 	err
	   }

	return s.Create()
}

// ExportDesktopEntry exports a desktop entry for an application.
func (s *SubSystem) ExportDesktopEntry(appName string) error {
	dbox ,err:= NewDbox()
	if(err!=nil){
	   return 	err
	   }

	return dbox.ContainerExportDesktopEntry(s.InternalName ,appName ,fmt.Sprintf("on %s", s.Name), s.IsRootfull )
}

// ExportDesktopEntries exports multiple desktop entries for applications.
func (s *SubSystem) ExportDesktopEntries(args ...string) (int,error){
	exportedN:=0

	for _, appName:=range args{
	   if(err:= s.ExportDesktopEntry(appName);err!=nil){
	      return exportedN ,err
	      }

	   exportedN++
	   }

	return exportedN,nil
}

// UnexportDesktopEntries unexports multiple desktop entries for applications.
func (s *SubSystem) UnexportDesktopEntries(args ...string)(int,error){
	exportedN:=0

	for _, appName:=range args{
	   if(err:= s.UnexportDesktopEntry(appName);err!=nil){
	      return exportedN ,err
	      }

	   exportedN++
	   }

	return exportedN,nil
}

// ExportBin exports a binary to a specified path.
func (s *SubSystem) ExportBin(binary string ,exportPath string)(error){
	if !strings.HasPrefix(binary,"/"){
	   binaryPath ,err:= s.Exec(true,false,"which",binary )
	   if(err!=nil){
	      return 	err
	      }

	   binary=strings.TrimSpace(binaryPath )
	   }

	binaryName:=filepath.Base(binary )

	dbox ,err:= NewDbox()
	if(err!=nil){
	   return 	err
	   }

	var homeDir string
	var homeErr error

	if homeDir ,homeErr= os.UserHomeDir();homeErr!=nil{
	     return homeErr
	     }

	if exportPath==""{
	     exportPath=filepath.Join(homeDir,".local","bin")
	     }

   joinedPath:=filepath.Join(exportPath,binaryName )
   if _,err= os.Stat(joinedPath);err==nil{
      tmpExportPath:=fmt.Sprintf("/tmp/%s",uuid.New().String())
      if mkErr:=os.MkdirAll(tmpExportPath ,0o755);mkErr!=nil{
          return mkErr
          }

      if expErr:=dbox.ContainerExportBin(s.InternalName,binary,tmpExportPath,s.IsRootfull );expErr!=nil{
          return expErr
          }

      copyErr:=CopyFile(filepath.Join(tmpExportPath,binaryName),filepath.Join(exportPath ,fmt.Sprintf("%s-%s",binaryName,s.InternalName)))
      if copyErr!=nil{
          return copyErr
          }

      removeErr:=os.RemoveAll(tmpExportPath )
      if removeErr!=nil{
          return removeErr
          }

      chmodErr:=os.Chmod(filepath.Join(exportPath,fmt.Sprintf("%s-%s",binaryName,s.InternalName)),0o755 )
      if chmodErr!=nil{
          return chmodErr
          }

      return nil
   }

   mkDirErr=os.MkdirAll(exportPath ,0o755 )
   if mkDirErr!=nil{
       return mkDirErr
       }

   expBinErr=dbox.ContainerExportBin(s.InternalName,binary ,exportPath,s.IsRootfull )
   if expBinErr!=nil{
       return expBinErr
       }

   return nil
}

// UnexportDesktopEntry unexports a desktop entry for an application.
func (s *SubSystem) UnexportDesktopEntry(appName string)(error){
	dbox ,err:= NewDbox()
	if(err!=nil){
	   return 	err
	   }

	return dbox.ContainerUnexportDesktopEntry(s.InternalName ,appName,s.IsRootfull )
}

// UnexportBin unexports a binary from the
