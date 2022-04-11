// postgresql/postgresql_test.go
package postgresql_test

import (
	"os"
	"pittsix/m/postgresql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDatabaseConnection(t *testing.T) {
	ass := assert.New(t)
	err := postgresql.GetDatabaseConnection()
	ass.NoError(err)
}

func TestGetDatabaseConnection_Missing_variable(t *testing.T) {
	ass := assert.New(t)
	os.Setenv("DATABASE_HOST", "")
	err := postgresql.GetDatabaseConnection()
	ass.NotNil(err)
	ass.Equal(err.Error(), "DATABASE_USER must be set")
}
