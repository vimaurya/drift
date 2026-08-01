package driver

import (
	"fmt"
	"sync"
)

type MigrationRecord struct {
	Version  int64
	Checksum string
}

var (
	mu sync.Mutex
	drivers = make(map[string]factory)
)

type Driver interface {
	InitializeMigrations() error
	GetAppliedMigrations() (map[int64]string, error)
	Apply(version int64, name, checksum, sql string) error
	Down(version int64, sql string) error
	Close()
}

type factory func(connURL string) (Driver, error) 

func DriverRegister(name string, factory factory) {
	mu.Lock()
	defer mu.Unlock()

	if factory==nil{
		panic("drift/driver : Register factory nil")
	}
	if _, dup := drivers[name]; dup{
		panic("drift/driver : Register called twice for driver "+name)
	}

	drivers[name] = factory
}


func getdbDriver(driverName string) (Driver, error) {
	if driverName=="postgresql"{
		driverName = "postgres"
	}
	
	mu.Lock()
	factory, ok := drivers[driverName]
	mu.Unlock()

	if !ok{
		return nil, fmt.Errorf("driver/dirft : unknown driver %q (forgotten import ?)", driverName)
	}

	return factory(driverName)
}


