package driver

import (
	"fmt"
	"net/url"
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
	if len(parts) > 2 {
		return nil, fmt.Errorf("Some err")
	}
	
	driver, err := getdbDriver(parts[0], connURL)
	if err!=nil{
		return nil, err
	}
	
	formatConnURL(connURL)

	return driver, nil
}

func formatConnURL(connUrl string) (string, error){
	parsedUrl, err := url.Parse(connUrl)	
	if err!=nil{
		return "", err
	}
	
	fmt.Printf("inside formatting\n")
	scheme := parsedUrl.Scheme
	if scheme == "postgres"{
		user := parsedUrl.User.Username()
		password, _ := parsedUrl.User.Password()
		host := parsedUrl.Host
		database := parsedUrl.Path
		query := parsedUrl.RawQuery
		fmt.Printf("user : %s\npassword : %s\nhost : %s\n", user, password, host)
		fmt.Printf("database : %s\nquery : %s\n", database, query)
	}

	return "", nil
}
