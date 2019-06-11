package symbols

import (
	"database/sql"
	"github.com/mattn/go-sqlite3"
	"github.com/prince1809/sourcegraph/pkg/env"
	"github.com/prometheus/common/log"
)

const maxFileSize = 1 << 19 // 512KB

var libSqlite3Pcre = env.Get("LIBSQLITE3_PCRE", "", "path to libsqlite3-pcre library")

func MustRegisterSqlite3WithPcre() {
	if libSqlite3Pcre == "" {
		env.PrintHelp()
		log.Fatal("Can't find the libsqlite3-pcre library becasue LIBSQLITE3_PCRE was not set")
	}
	sql.Register("sqlite3_with_pcre", &sqlite3.SQLiteDriver{Extensions: []string{libSqlite3Pcre}})
}
