package models

import "gorm.io/gorm"

type Artifacts struct {
	ID      uint    `gorm:"primary key;autoIncrement" json:"id""`
	Student *string `json:"student"`
	Type    *string `json:"type"`
	Site    *string `json:"site"`
}

func MigrateArtifacts(db *gorm.DB) error {
	err := db.AutoMigrate(&Artifacts{})
	return err
}
