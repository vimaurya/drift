package drift

import (
	"log"

	"github.com/di-argus/drift/internal/config"
	"github.com/di-argus/drift/internal/core"
	"github.com/di-argus/drift/internal/driver"
)

type Migrator struct{
	Config *config.Config
	Driver driver.Driver
}

func Init() error {
	return nil
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
