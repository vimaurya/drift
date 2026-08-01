package driver

import (
	"database/sql"
	"net/url"
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
	connurl, err:= url.Parse(connURL)
	if err!=nil{
		return nil, err 
	}

	db, err := sql.Open(connurl.Scheme, connURL)
	if err!=nil{
		return nil, err
	}

	driver, err := getdbDriver(connurl.Scheme, db)
	if err!=nil{
		return nil, err
	}

	return driver, nil
}
