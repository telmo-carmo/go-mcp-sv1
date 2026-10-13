// using ORM Bun - https://github.com/uptrace/bun
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
	"github.com/uptrace/bun"
        "github.com/uptrace/bun/dialect/sqlitedialect"
        "github.com/uptrace/bun/driver/sqliteshim"

)

// User defines the database model using struct tags for Bun mapping.
type User struct {
	ID        int64     `bun:"id,pk,autoincrement"`
	Name      string    `bun:"name,notnull"`
	Email     string    `bun:"email,unique,notnull"`
	CreatedAt time.Time `bun:"created_at,default:current_timestamp"`
}

func main() {
	ctx := context.Background()

	// 1. Open the underlying sql.DB connection using the Bun SQLite shim.
	// For an in-memory DB use "file::memory:?cache=shared". For a file use a filename like "dev.db".
	sqldb, err := sql.Open(sqliteshim.ShimName, "file::memory:?cache=shared")


	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer sqldb.Close()

	// 2. Wrap the sql.DB instance with a Bun Database wrapper using the SQLite dialect
	db := bun.NewDB(sqldb, sqlitedialect.New()) // pgdialect.New())
	defer db.Close()

	// 3. Create the users table if it doesn't already exist
	_, err = db.NewCreateTable().
		Model((*User)(nil)).
		IfNotExists().
		Exec(ctx)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
	log.Println("Database and table successfully initialized.")

	// 4. Insert a new user into the database
	newUser := &User{
		Name:  "Jane Doe",
		Email: "jane.doe@example.com",
	}

	_, err = db.NewInsert().
		Model(newUser).
		Exec(ctx)
	if err != nil {
		log.Fatalf("Failed to insert user: %v", err)
	}
	fmt.Printf("Inserted User ID: %d\n", newUser.ID)

	// 5. Query the user back from the database by ID
	var fetchedUser User
	err = db.NewSelect().
		Model(&fetchedUser).
		Where("id = ?", newUser.ID).
		Scan(ctx)
	if err != nil {
		log.Fatalf("Failed to select user: %v", err)
	}

	fmt.Printf("Fetched User from DB: ID=%d, Name=%s, Email=%s, CreatedAt=%v\n",
		fetchedUser.ID, fetchedUser.Name, fetchedUser.Email, fetchedUser.CreatedAt)

	var users []User
	db.NewSelect().Model(&users).Scan(ctx)
	fmt.Printf("Users: %+v\n", users)

}


