package gormigrate

import (
	"time"

	"github.com/jinzhu/gorm"
)

type MigrationVersion struct {
	SelfName string `gorm:"-"`
	Id       int64  `gorm:"AUTO_INCREMENT"`
	Version  int64  `gorm:"index"`
	Status   int8
	ExecTime time.Time
}

func (m MigrationVersion) TableName() string {
	if m.SelfName == "" {
		m.SelfName = DefaultOptions.VersionTableName
	}
	return m.SelfName
}

// ensure migration version table
func EnsureMigrationVersion(db *gorm.DB, m *MigrationVersion) error {
	if db.HasTable(m) == true {
		return nil
	} else {
		return CreateMigrationVersion(db, m)
	}
}

// create migration version table
func CreateMigrationVersion(db *gorm.DB, m *MigrationVersion) error {
	tx := db.Begin()
	if err := tx.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(m).Error; err != nil {
		tx.Rollback()
		return err
	}
	// first run version has 0
	m.Version = 0
	m.Status = 1
	m.ExecTime = Now()
	if err := tx.Create(m).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
