package core

import (
	"fmt"
	"github.com/AuruOS/apx/v2/settings"
)

var abg *Abg

type Abg struct {
	Cnf *settings.Config
}

func NewAbg(cnf *settings.Config) *Abg {
	abg = &Abg{
		Cnf: cnf,
	}

	err := abg.EssentialChecks()
	if err != nil {
		// localisation features aren't available at this stage, so this error can't be translated
		fmt.Println("ERROR: Unable to find abg configuration files")
		return nil
	}

	return abg
}

func NewStandardAbg() *Abg {
	cnf, err := settings.GetApxDefaultConfig()
	if err != nil {
		panic(err)
	}

	abg = &Abg{
		Cnf: cnf,
	}

	err = abg.EssentialChecks()
	if err != nil {
		// localisation features aren't available at this stage, so this error can't be translated
		fmt.Println("ERROR: Unable to find abg configuration files")
		return nil
	}
	return abg
}
