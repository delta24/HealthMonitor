package setup

import (
	"database/sql"
	"os"
	"path"

	"github.com/mattn/go-sqlite3"

	"health_monitor/utils"
)

// Database holds the active sqlite connection
var (
	Database  *sql.DB
	DBLogFile *os.File
)

func dbInit() {
	var (
		DBDriver    string
		err         error
		logFileName = path.Join(ConfigVars.HomeDir, "db.log")
	)
	DBLogFile, err = os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0666)
	if err != nil {
		utils.PLogError(err)
	}
	sql.Register(DBDriver, &sqlite3.SQLiteDriver{})
	Database, err = sql.Open(DBDriver, ConfigVars.DBFile)
	if err != nil {
		utils.ModuleError(DBLogFile, "Failed to create the handle", err.Error())
	}
	if err = Database.Ping(); err != nil {
		utils.ModuleError(DBLogFile, "Failed to keep connection alive", err.Error())
	}

	if stats, _ := os.Stat(ConfigVars.DBFile); stats.Size() == 0 {
		utils.ModuleError(DBLogFile, "Database not found", "Creating one with default values.")
		setupDB()
	}
}
