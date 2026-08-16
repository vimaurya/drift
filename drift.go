package drift

import (
	"log"

	"github.com/Di-Argus/Drift/pkg/config"
	"github.com/Di-Argus/Drift/pkg/core"
	"github.com/Di-Argus/Drift/pkg/driver"
)

type Migrator struct{
	Config *config.Config
	Driver driver.Driver
}

func (u *Migrator) RunUp() error {
	err := core.RunUp(*u.Config, u.Driver)
	if err!=nil{
		log.Printf("%v", err)
	}
	return err
}

func (u *Migrator) RunDown() error {
	err := core.RunDown(*u.Config, u.Driver)
	if err!=nil{
		log.Printf("%v", err)
	}

	return err
}
