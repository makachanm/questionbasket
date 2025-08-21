package frame

import (
	"database/sql"

	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3"
)

var dbManager DBModelManager

func init() {
	dbManager = DBModelManager{DbOpened: false}
}

type SQLDriverType int

type SQLDatabaseConnectionConfig struct {
	DBType SQLDriverType

	ConnectionAddress string
	Port              int
	Username          string
	Password          string
	DBName            string
}

const (
	DRIVER_SQLITE SQLDriverType = iota
)

type DatabaseConnector struct {
	Connector sqlx.DB
}

type DBModelManager struct {
	DbDriver sqlx.DB

	DbOpened bool

	//MigrationBindTargets []string
	//TODO: model migration management
	//TODO: Pagination
}

func InitalizeDatabaseConnection(connConfig SQLDatabaseConnectionConfig) {
	if dbManager.DbOpened {
		return
	}

	switch connConfig.DBType {
	case DRIVER_SQLITE:
		dbManager.DbDriver = *sqlx.MustConnect("sqlite3", connConfig.ConnectionAddress)

		//currently not consider
		/*
			case DRIVER_POSTGRESQL:
				//host={} port={} user={} password={} dbname={}
				connStr := fmt.Sprintf("user=%s password=%s dbname=%s port=%d host=%s",
					connConfig.Username, connConfig.Password, connConfig.DBName, connConfig.Port, connConfig.ConnectionAddress)
				dbManager.DbDriver = *sqlx.MustConnect("postgres", connStr)
		*/
	}
}

func CloseDatabaseConnection() {
	dbManager.DbOpened = false
	dbManager.DbDriver.Close()
}

func DatabaseBind(model DatabaseModel) {
	connetor := DatabaseConnector{Connector: dbManager.DbDriver}
	model.DatabaseConnect(connetor)
}

func GetDB() *sqlx.DB {
	return &dbManager.DbDriver
}

func (dbc *DatabaseConnector) RunPrepared(query string, v any) (sql.Result, error) {
	tx := dbc.Connector.MustBegin()

	statement, qerr := tx.PrepareNamed(query)
	if qerr != nil {
		return nil, qerr
	}
	defer statement.Close()

	result := statement.MustExec(v)
	return result, tx.Commit()
}

// RunPreparedSelect executes a query that returns multiple rows and scans them into a slice.
func (dbc *DatabaseConnector) RunPreparedSelect(dest any, query string, arg any) error {
	tx, err := dbc.Connector.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback() // Rollback is a no-op if the transaction is committed

	nstmt, err := tx.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer nstmt.Close()

	err = nstmt.Select(dest, arg)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// RunPreparedRow executes a query that is expected to return a single row.
func (dbc *DatabaseConnector) RunPreparedRow(query string, v any) (*sqlx.Row, error) {
	nstmt, err := dbc.Connector.PrepareNamed(query)
	if err != nil {
		return nil, err
	}
	// Not closing the prepared statement here because sqlx caches them.

	return nstmt.QueryRowx(v), nil
}

type ModelInfo struct {
	Name string
}

type DatabaseModel interface {
	DatabaseConnect(dbc DatabaseConnector)
	GetModelInfo() ModelInfo
}

// DatabaseModel
/*
	- RunMigrate()
	- Get()
	- Update()
	- Insert()
	- Delete()
	- Custom Query things....
*/
