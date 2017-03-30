package gormigrate_test

import (
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	gormigrate "github.com/lvzhihao/gormigrate/libs"
)

type TestTable struct {
	gorm.Model
	Name string
}

var (
	db          *gorm.DB
	m           *gormigrate.Gormigrate
	createTable = &gormigrate.Migration{
		Id: "20170329175903",
		Up: func(tx *gorm.DB) error {
			return tx.CreateTable(&TestTable{}).Error
		},
		Down: func(tx *gorm.DB) error {
			return tx.DropTable(&TestTable{}).Error
		},
	}
)

func init() {
	db, _ = gorm.Open("mysql", "root:@/gormigrate_test?parseTime=True&loc=Asia%2FShanghai")
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "prefix_" + defaultTableName
	}
	m = gormigrate.New(db, gormigrate.DefaultOptions)
}

func Test_001_LastVersion(t *testing.T) {
	v, e := m.GetLastVersion()
	if e != nil {
		t.Error(e)
	} else {
		t.Log("LastVersion:", v)
	}
}

func Test_002_MigrateUp(t *testing.T) {
	m.AddMigration(createTable)
	m.Migrate()
}
