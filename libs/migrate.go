package gormigrate

import (
	"errors"

	"github.com/jinzhu/gorm"
)

type MigrateUp func(*gorm.DB) error
type MigrateDown func(*gorm.DB) error

type Gormigrate struct {
	db         *gorm.DB
	tx         *gorm.DB
	options    *Options
	version    *MigrationVersion
	migrations []*Migration
}

type Options struct {
	VersionTableName string
	UseTransaction   bool
}

type Migration struct {
	Id   string
	Up   MigrateUp
	Down MigrateDown
}

var (
	DefaultOptions *Options = &Options{
		VersionTableName: "gormigrate_db_version",
		UseTransaction:   true,
	}
	ErrNoLastValidMigration = errors.New("Could not find last valid migration")
)

func New(db *gorm.DB, options *Options) *Gormigrate {
	return &Gormigrate{
		db:      db,
		options: options,
	}
}

func (this *Gormigrate) AddMigration(m *Migration) {
	//todo sort by Id
	this.migrations = append(this.migrations, m)
}

func (this *Gormigrate) Migrate() error {
	_, err := this.GetLastVersion()
	if err != nil {
		return err
	}
	this.begin()

	for _, m := range this.migrations {
		if err := this.runMigrate(m); err != nil {
			this.rollback()
			return err
		}
	}

	return this.commit()
}

func (this *Gormigrate) runMigrate(m *Migration) error {
	return nil
}

func (this *Gormigrate) GetLastVersion() (int64, error) {
	if err := this.ensureVersion(); err != nil {
		return 0, err
	}
	rows, err := this.db.Model(this.version).Select("version, status").Order("exec_time desc").Rows()
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	var skipList []int64
	for rows.Next() {
		var version int64
		var status int8
		err := rows.Scan(&version, &status)
		if err != nil {
			return 0, err
		}
		skip := false
		for _, s := range skipList {
			if version == s {
				skip = true
				break
			}
		}
		if skip {
			continue
		}
		if status == 1 {
			return version, nil
		}
		skipList = append(skipList, version)
	}
	return 0, ErrNoLastValidMigration
}

func (this *Gormigrate) ensureVersion() error {
	version := &MigrationVersion{SelfName: this.options.VersionTableName}
	err := EnsureMigrationVersion(this.db, version)
	if err != nil {
		return err
	} else {
		this.version = version
		return nil
	}
}

func (this *Gormigrate) begin() {
	if this.options.UseTransaction {
		this.tx = this.db.Begin()
	} else {
		this.tx = this.db
	}
}

func (this *Gormigrate) rollback() {
	if this.options.UseTransaction {
		this.tx.Rollback()
	}
}

func (this *Gormigrate) commit() error {
	if this.options.UseTransaction {
		return this.tx.Commit().Error
	} else {
		return nil
	}
}
