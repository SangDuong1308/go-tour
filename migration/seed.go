package migration

import (
	"gorm.io/gorm"
)

type Seed struct {
	Name string
	Run  func(*gorm.DB) error
}

func All() []Seed {
	return []Seed{
		//{
		//	Name: "Create user testing",
		//	Run: func(db *gorm.DB) error {
		//		return seed.FakeMember(db)
		//	},
		//},
	}
}
