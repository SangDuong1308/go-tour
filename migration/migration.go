package migration

import (
	"fmt"
	"github.com/pkg/errors"
	"go-tour/internal/models"
	"gorm.io/gorm"
)

func Migration(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&models.User{},
		//&models.CustodialPaymentAddressAction{},
		//&models.CustodialPaymentAddressLog{},
	); err != nil {
		return errors.Wrap(err, "db.AutoMigrate")
	}

	return nil
}

func AutoSeedingData(db *gorm.DB) error {
	for _, seed := range All() {
		if err := seed.Run(db); err != nil {
			return errors.Wrap(err, fmt.Sprintf("Running seed '%s'", seed.Name))
		}
	}

	return nil
}
