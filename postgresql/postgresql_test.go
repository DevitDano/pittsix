// postgresql/postgresql_test.go
package postgresql_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
)

func TestGetDatabaseConnection(t *testing.T) {
	ass := assert.New(t)
	err := postgres.GetDatabaseConnection()
	ass.NoError(err)
}

func TestGetDatabaseConnection_Missing_variable(t *testing.T) {
	ass := assert.New(t)
	os.Setenv("DATABASE_HOST", "")
	err := postgres.GetDatabaseConnection()
	ass.NotNil(err)
	ass.Equal(err.Error(), "DATABASE_USER must be set")
}
