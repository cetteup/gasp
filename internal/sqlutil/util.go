package sqlutil

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

func Connect(host, dbname, user, passwd string) *sql.DB {
	cfg := mysql.Config{
		User:                 user,
		Passwd:               passwd,
		Net:                  "tcp",
		Addr:                 host,
		DBName:               dbname,
		Loc:                  time.UTC,
		MaxAllowedPacket:     64 << 20, // same as mysql.defaultMaxAllowedPacket
		ParseTime:            true,
		AllowNativePasswords: true,
		CheckConnLiveness:    true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(76)
	db.SetMaxIdleConns(76)

	return db
}

func EscapeWildcards(s string) string {
	r := strings.NewReplacer(
		"%", "\\%",
		"_", "\\_",
	)
	return r.Replace(s)
}

func Quote(s string) string {
	return "`" + s + "`"
}

func QuoteJoin(table, column, sep string) string {
	return Quote(table) + sep + Quote(column)
}

func Qualify(table, column string) string {
	return QuoteJoin(table, column, ".")
}

func QualifyAlias(table, column string) string {
	return fmt.Sprintf("%s AS %s", Qualify(table, column), Predicate(table, column))
}

func Predicate(table, column string) string {
	return Quote(table + "_" + column)
}
