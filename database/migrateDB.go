package database

import (
	"database/sql"
	"embed"
	"fmt"

	migrate "github.com/rubenv/sql-migrate"
)

//go:embed sql_migrations/*.sql
var dbMigrations embed.FS

var DbConnection *sql.DB

func DBMigrate(dbParam *sql.DB){
	migrations := &migrate.EmbedFileSystemMigrationSource{
		FileSystem: dbMigrations,
		Root: "sql_migrations",
	}


	// Cara migrate Up
	n, errs := migrate.Exec(dbParam, "postgres", migrations, migrate.Up)
	if errs != nil {
		panic(errs)
	}


	// Cara migrate Down
	// n, errs := migrate.Exec(dbParam, "postgres", migrations, migrate.Down)
	// if errs != nil {
	// 	panic(errs)
	// }

	// fmt.Println("Migration Down, applied", n ,"migrations")

	DbConnection = dbParam	

	fmt.Println("Migration Success, applied", n , "migrations")
}