// psql/postgresql.go
package psql

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var credentialsMap = map[string]string{
	"DATABASE_HOST":     os.Getenv("DATABASE_HOST"),
	"DATABASE_USER":     os.Getenv("DATABASE_USER"),
	"DATABASE_PASSWORD": os.Getenv("DATABASE_PASSWORD"),
	"DATABASE_NAME":     os.Getenv("DATABASE_NAME"),
	"DATABASE_PORT":     os.Getenv("DATABASE_PORT"),
}

// DBCN a global db object will be used across different packages
var DBCN *gorm.DB

// GetDatabaseConnection create/return dabase connection
func GetDatabaseConnection() (err error) {
	for key, val := range credentialsMap {
		if val == "" {
			return fmt.Errorf("%s must be set", key)
		}
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", credentialsMap["DATABASE_HOST"], credentialsMap["DATABASE_USER"], credentialsMap["DATABASE_PASSWORD"], credentialsMap["DATABASE_NAME"], credentialsMap["DATABASE_PORT"])
	DBCN, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	return
}
