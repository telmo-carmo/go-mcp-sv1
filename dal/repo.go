package dal

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/phuslu/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
)

type Book struct {
	ID     uint   `json:"id" bun:",pk,autoincrement"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

type Product struct {
	ID        uint `bun:",pk,autoincrement"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	Code      string
	Price     uint
}

type Bonus struct {
	Ename string `json:"ename"`
	Job   string `json:"job"`
	Sal   int    `json:"sal"`
	Comm  int    `json:"comm"`
}

type Dept struct {
	bun.BaseModel `bun:"table:dept"`

	Deptno int    `json:"deptno" bun:"DEPTNO,pk"`
	Name   string `json:"name"   bun:"NAME,notnull"`
	Loc    string `json:"loc"    bun:"LOC,notnull"`
}

var DB *bun.DB

// InitDatabase initializes the database and creates the application schema.
func InitDatabase() error {
	dbname := os.Getenv("DB_CONNECTION")
	if dbname == "" {
		dbname = "mydata1.db"
	}
	sqldb, err := sql.Open(sqliteshim.ShimName, dbname)
	if err != nil {
		return fmt.Errorf("open database %q: %w", dbname, err)
	}
	db := bun.NewDB(sqldb, sqlitedialect.New())
	db.WithQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	log.Info().Msgf("using db: %s", dbname)
	if _, err := db.NewCreateTable().Model((*Book)(nil)).IfNotExists().Exec(context.Background()); err != nil {
		db.Close()
		return fmt.Errorf("create book table: %w", err)
	}

	DB = db
	return nil
}
