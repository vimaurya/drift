package tests

import (
	"fmt"
	"testing"

	"github.com/Di-Argus/Drift/internal/driver"
	_ "github.com/Di-Argus/Drift/internal/driver/mysql"
	_ "github.com/Di-Argus/Drift/internal/driver/postgres"
	"github.com/stretchr/testify/assert"
)


func TestGetDriver_Success(t *testing.T) {
	tests := []struct{
		testname string
		connurl string
		err bool
	}{
		{
			testname: "valid get driver (mysql)",
			connurl : "mysql://root:vikash@localhost:3306/test_db",
			err : false, 
		},
		{
			testname: "valid get driver (postgres)",
			connurl : "postgres://root:vikash@localhost:5432/test_db", 
			err : false, 
		},
		{
			testname: "valid get driver (postgresql)",
			connurl : "postgresql://root:vikash@localhost:5432/test_db", 
			err : false, 
		},
		{
			testname: "invalid connurl",
			connurl: "mysql://root:vikash@tcp(127.0.0.1:3306)/test_db",
			err : true, 
		},
		{
			testname: "invalid get driver",
			connurl : "postgresql:///root:vikash@localhost:5432/test_db", 
			err : false, 
		},
		{
			testname: "invalid get driver",
			connurl : "://root:vikash@localhost:5432/test_db", 
			err : true, 
		},

	}

	for _, test := range tests {
		t.Run(test.testname, func(t *testing.T) {
			_, err := driver.GetDriver(test.connurl)
			if test.err{
				msg := fmt.Sprintf("test : %s,\nconnurl : %s", test.testname, test.connurl)
				assert.Error(t, err, msg)
			} else {
				msg := fmt.Sprintf("test : %s,\nconnurl : %s", test.testname, test.connurl)
				assert.NoError(t, err, msg)
			}
		})	
	}
}
