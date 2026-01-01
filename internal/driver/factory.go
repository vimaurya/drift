package driver

import (
	"fmt"
	"net/url"
)

func GetDriver(connURL string) (Driver, error) {
	u, err := url.Parse(connURL)
	if err!=nil {
		return nil, err
	}
	switch u.Scheme {
	case "postgres", "postgresql":
		fmt.Println("calling nw")
		return NewPostgresDriver(connURL)
	default:
		return nil, fmt.Errorf("unsupported database: %s", u.Scheme)	
	}
}
