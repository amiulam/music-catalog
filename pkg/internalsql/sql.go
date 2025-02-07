package internalsql

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(dataSourceName string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dataSourceName), &gorm.Config{})

	if err != nil {
		log.Fatalf("Error when try to connect to database %+v\n", err)
		return nil, err
	}

	return db, nil
}
