package driver

import (
	"fmt"
	"strings"
)

// func GetDriver(connURL string) (Driver, error) {
// 	if strings.HasPrefix(connURL, "postgres://") || strings.HasPrefix(connURL, "postgresql://") {
// 		return NewPostgresDriver(connURL)
// 	}
//
// 	if strings.HasPrefix(connURL, "mysql://") {
// 		dsn := strings.TrimPrefix(connURL, "mysql://")
// 		return NewMySQLDriver(dsn)
// 	}
//
// 	return nil, fmt.Errorf("unsupported database scheme name in url : %s", connURL)
// }

func GetDriver(connURL string) (Driver, error) {
	parts := strings.SplitN(connURL, "://", 2)
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid connection url format (missing scheme '://')")
	}
	
	scheme := parts[0]

	driver, err := getdbDriver(scheme)
	if err!=nil{
		return nil, err
	}

	return driver, nil
}
