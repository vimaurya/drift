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

type connectionUrl struct {
	scheme string
	url    string
}

func GetDriver(connURL string) (Driver, error) {
	connectionURL, err := formatConnURL(connURL)
	if err!=nil{
		return nil, err
	}

	driver, err := getdbDriver(connectionURL.scheme, connectionURL.url)

	return driver, nil
}

// Standard Drift URL format is :
// scheme://username:passowrd@host:port/database?queryparam=something

// mysql://root:vikash@tcp(127.0.0.1:3306)/test_db
// mysql://root:vikash@localhost:3306/test_db
func formatConnURL(connUrl string) (connectionUrl, error) {
	parsedUrl, err := url.Parse(connUrl)
	if err != nil {
		return connectionUrl{}, err
	}

	conn := &connectionUrl{}

	scheme := parsedUrl.Scheme

	user := cleanString(parsedUrl.User.Username())
	password, _ := parsedUrl.User.Password()
	password = cleanString(password)
	host := cleanString(parsedUrl.Host)
	database := cleanString(parsedUrl.Path)
	query := cleanString(parsedUrl.RawQuery)

	// fmt.Printf("user : %s\npassword : %s\nhost : %s\n", user, password, host)
	// fmt.Printf("database : %s\nquery : %s\n", database, query)

	var connectionString string

	if scheme == "postgres" {
		connectionString = fmt.Sprintf("%s://%s:%s@%s/%s", scheme, user, password, host, database)
		if query != "" {
			connectionString += "?" + query
		}
	} else if scheme == "mysql" {
		connectionString = fmt.Sprintf("%s://%s:%s@tcp(%s)/%s", scheme, user, password, host, database)
	}

	conn.scheme = scheme
	conn.url = connectionString

	return *conn, nil
}

func cleanString(dirtyString string) string {
	clString, _ := strings.CutPrefix(dirtyString, ":")
	clString, _ = strings.CutSuffix(clString, ":")
	clString, _ = strings.CutPrefix(clString, "/")
	clString, _ = strings.CutSuffix(clString, "/")
	clString, _ = strings.CutPrefix(clString, "@")
	clString, _ = strings.CutSuffix(clString, "@")
	return clString
}
