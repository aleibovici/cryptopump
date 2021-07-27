package mysql

import (
	"cryptopump/functions"
	"cryptopump/types"
	"database/sql"
	"fmt"
	"math"
	"os"

	log "github.com/sirupsen/logrus"

	_ "github.com/go-sql-driver/mysql" // This blank entry is required to enable mysql connectivity
)

// DBInit export
/* This function initializes GCP mysql database connectivity */
func DBInit() *sql.DB {

	var db *sql.DB
	var err error

	// If the optional DB_TCP_HOST environment variable is set, it contains
	// the IP address and port number of a TCP connection pool to be created,
	// such as "127.0.0.1:3306". If DB_TCP_HOST is not set, a Unix socket
	// connection pool will be created instead.
	if os.Getenv("DB_TCP_HOST") != "" {

		if db, err = InitTCPConnectionPool(); err != nil {

			log.Fatalf("initTCPConnectionPool: unable to connect: %v", err)

		}

	} else {

		if db, err = InitSocketConnectionPool(); err != nil {

			log.Fatalf("initSocketConnectionPool: unable to connect: %v", err)

		}

	}

	return db

}

// InitSocketConnectionPool initializes a Unix socket connection pool for
// a Cloud SQL instance of SQL Server.
func InitSocketConnectionPool() (*sql.DB, error) {

	var err error
	var dbPool *sql.DB

	// [START cloud_sql_mysql_databasesql_create_socket]
	var (
		dbUser                 = functions.MustGetenv("DB_USER")
		dbPwd                  = functions.MustGetenv("DB_PASS")
		instanceConnectionName = functions.MustGetenv("INSTANCE_CONNECTION_NAME")
		dbName                 = functions.MustGetenv("DB_NAME")
	)

	socketDir, isSet := os.LookupEnv("DB_SOCKET_DIR")
	if !isSet {
		socketDir = "/cloudsql"
	}

	var dbURI = fmt.Sprintf("%s:%s@unix(/%s/%s)/%s?parseTime=true", dbUser, dbPwd, socketDir, instanceConnectionName, dbName)

	// dbPool is the pool of database connections.
	if dbPool, err = sql.Open("mysql", dbURI); err != nil {

		return nil, fmt.Errorf("sql.Open: %v", err)

	}

	// [START_EXCLUDE]
	configureConnectionPool(dbPool)
	// [END_EXCLUDE]

	return dbPool, nil
	// [END cloud_sql_mysql_databasesql_create_socket]
}

// configureConnectionPool sets database connection pool properties.
// For more information, see https://golang.org/pkg/database/sql
func configureConnectionPool(dbPool *sql.DB) {
	// [START cloud_sql_mysql_databasesql_limit]

	// Set maximum number of connections in idle connection pool.
	dbPool.SetMaxIdleConns(5)

	// Set maximum number of open connections to the database.
	dbPool.SetMaxOpenConns(7)

	// [END cloud_sql_mysql_databasesql_limit]

	// [START cloud_sql_mysql_databasesql_lifetime]

	// Set Maximum time (in seconds) that a connection can remain open.
	dbPool.SetConnMaxLifetime(1800)

	// [END cloud_sql_mysql_databasesql_lifetime]
}

// InitTCPConnectionPool initializes a TCP connection pool for a Cloud SQL
// instance of SQL Server.
func InitTCPConnectionPool() (*sql.DB, error) {

	var err error
	var dbPool *sql.DB

	// [START cloud_sql_mysql_databasesql_create_tcp]
	var (
		dbUser    = functions.MustGetenv("DB_USER")
		dbPwd     = functions.MustGetenv("DB_PASS")
		dbTCPHost = functions.MustGetenv("DB_TCP_HOST")
		dbPort    = functions.MustGetenv("DB_PORT")
		dbName    = functions.MustGetenv("DB_NAME")
	)

	var dbURI = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPwd, dbTCPHost, dbPort, dbName)

	// dbPool is the pool of database connections.

	if dbPool, err = sql.Open("mysql", dbURI); err != nil {

		return nil, fmt.Errorf("sql.Open: %v", err)

	}

	// [START_EXCLUDE]
	configureConnectionPool(dbPool)
	// [END_EXCLUDE]

	return dbPool, nil
	// [END cloud_sql_mysql_databasesql_create_tcp]
}

// SaveOrder Save order to database
func SaveOrder(
	sessionData *types.Session,
	ClientOrderID string,
	CummulativeQuoteQuantity float64,
	ExecutedQuantity float64,
	OrderID int64,
	Price float64,
	Side string,
	Status string,
	Symbol string,
	TransactTime int64) (err error) {

	var rows *sql.Rows

	if rows, err = sessionData.Db.Query("call cryptopump.SaveOrder(?,?,?,?,?,?,?,?,?,?,?)",
		ClientOrderID,
		CummulativeQuoteQuantity,
		ExecutedQuantity,
		OrderID,
		Price,
		Side,
		Status,
		Symbol,
		TransactTime,
		sessionData.ThreadID,
		sessionData.ThreadIDSession); err != nil {

		functions.Logger(
			nil,
			nil,
			sessionData,
			log.DebugLevel,
			0,
			OrderID,
			Price,
			0,
			functions.GetFunctionName()+" - "+err.Error())

		return err

	}

	rows.Close()

	return nil

}

// UpdateOrder Update order
func UpdateOrder(
	sessionData *types.Session,
	OrderID int64,
	CummulativeQuoteQuantity float64,
	ExecutedQuantity float64,
	Price float64,
	Status string) (err error) {

	var rows *sql.Rows

	if rows, err = sessionData.Db.Query("call cryptopump.UpdateOrder(?,?,?,?,?)",
		OrderID,
		CummulativeQuoteQuantity,
		ExecutedQuantity,
		Price,
		Status); err != nil {

		functions.Logger(
			nil,
			nil,
			sessionData,
			log.DebugLevel,
			0,
			OrderID,
			Price,
			0,
			functions.GetFunctionName()+" - "+err.Error())

		return err

	}

	rows.Close()

	return nil

}

// UpdateSession Update existing session on Session table
func UpdateSession(
	configData *types.Config,
	sessionData *types.Session) (err error) {

	var rows *sql.Rows

	if rows, err = sessionData.Db.Query("call cryptopump.UpdateSession(?,?,?,?,?)",
		sessionData.ThreadID,
		sessionData.ThreadIDSession,
		configData.ExchangeName.(string),
		sessionData.Symbol_fiat,
		sessionData.Symbol_fiat_funds); err != nil {

		functions.Logger(
			configData,
			nil,
			sessionData,
			log.DebugLevel,
			0,
			0,
			0,
			0,
			functions.GetFunctionName()+" - "+err.Error())

		return err

	}

	rows.Close()

	return nil

}

// SaveSession Save new session to Session table.
func SaveSession(
	configData *types.Config,
	sessionData *types.Session) (err error) {

	var rows *sql.Rows

	if rows, err = sessionData.Db.Query("call cryptopump.SaveSession(?,?,?,?,?)",
		sessionData.ThreadID,
		sessionData.ThreadIDSession,
		configData.ExchangeName.(string),
		sessionData.Symbol_fiat,
		sessionData.Symbol_fiat_funds); err != nil {

		functions.Logger(
			configData,
			nil,
			sessionData,
			log.DebugLevel,
			0,
			0,
			0,
			0,
			functions.GetFunctionName()+" - "+err.Error())

		return err

	}

	rows.Close()

	return nil

}

// DeleteSession Delete session from Session table
func DeleteSession(
	sessionData *types.Session) (err error) {

	var rows *sql.Rows

	if rows, err = sessionData.Db.Query("call cryptopump.DeleteSession(?)",
		sessionData.ThreadID); err != nil {

		functions.Logger(
			nil,
			nil,
			sessionData,
			log.DebugLevel,
			0,
			0,
			0,
			0,
			functions.GetFunctionName()+" - "+err.Error())

		return err

	}

	rows.Close()

	return nil

}

// SaveThreadTransaction Save Thread cycle to database
func SaveThreadTransaction(
	sessionData *types.Session,
	OrderID int64,
	CummulativeQuoteQuantity float64,
	Price float64,
	ExecutedQuantity float64) (err error) {

	var rows *sql.Rows

	if rows, err = sessionData.Db.Query("call cryptopump.SaveThreadTransaction(?,?,?,?,?,?)",
		sessionData.ThreadID,
		sessionData.ThreadIDSession,
		OrderID,
		CummulativeQuoteQuantity,
		Price,
		ExecutedQuantity); err != nil {

		functions.Logger(
			nil,
			nil,
			sessionData,
			log.DebugLevel,
			0,
			OrderID,
			Price,
			0,
			functions.GetFunctionName()+" - "+err.Error())

		return err

	}

	rows.Close()

	return nil

}

// DeleteThreadTransactionByOrderID function
func DeleteThreadTransactionByOrderID(
	sessionData *types.Session,
	orderID int) (err error) {

	var rows *sql.Rows

	if rows, err = sessionData.Db.Query("call cryptopump.DeleteThreadTransactionByOrderID(?)",
		orderID); err != nil {

		functions.Logger(
			nil,
			nil,
			sessionData,
			log.DebugLevel,
			orderID,
			0,
			0,
			0,
			functions.GetFunctionName()+" - "+err.Error())

		return err

	}

	rows.Close()

	return nil

}

// GetThreadTransactionCount Get Thread count
func GetThreadTransactionCount(
	sessionData *types.Session) (count int, err error) {

	var rows *sql.Rows

	if rows, err = sessionData.Db.Query("call cryptopump.GetThreadTransactionCount(?)",
		sessionData.ThreadID); err != nil {

		functions.Logger(
			nil,
			nil,
			sessionData,
			log.DebugLevel,
			0,
			0,
			0,
			0,
			functions.GetFunctionName()+" - "+err.Error())

		return 0, err

	}

	for rows.Next() {
		err = rows.Scan(&count)
	}

	rows.Close()

	return count, err

}

// GetLastOrderTransactionPrice Get time for last transaction the ThreadID
func GetLastOrderTransactionPrice(
	sessionData *types.Session,
	Side string) (price float64, err error) {

	var rows *sql.Rows

	if rows, err = sessionData.Db.Query("call cryptopump.GetLastOrderTransactionPrice(?,?)",
		sessionData.ThreadID,
		Side); err != nil {

		functions.Logger(
			nil,
			nil,
			sessionData,
			log.DebugLevel,
			0,
			0,
			0,
			0,
			functions.GetFunctionName()+" - "+err.Error())

		return 0, err

	}

	for rows.Next() {
		err = rows.Scan(&price)
	}

	rows.Close()

	return price, err

}

// GetLastOrderTransactionSide Get Side for last transaction the ThreadID
func GetLastOrderTransactionSide(
	sessionData *types.Session) (side string, err error) {

	var rows *sql.Rows

	if rows, err = sessionData.Db.Query("call cryptopump.GetLastOrderTransactionSide(?)",
		sessionData.ThreadID); err != nil {

		functions.Logger(
			nil,
			nil,
			sessionData,
			log.DebugLevel,
			0,
			0,
			0,
			0,
			functions.GetFunctionName()+" - "+err.Error())

		return "", err

	}

	for rows.Next() {
		err = rows.Scan(&side)
	}

	rows.Close()

	return side, err

}

// GetOrderTransactionSideLastTwo function
func GetOrderTransactionSideLastTwo(
	sessionData *types.Session) (side1 string, side2 string, err error) {

	var rows *sql.Rows

	if rows, err = sessionData.Db.Query("call cryptopump.GetOrderTransactionSideLastTwo(?)",
		sessionData.ThreadID); err != nil {

		functions.Logger(
			nil,
			nil,
			sessionData,
			log.DebugLevel,
			0,
			0,
			0,
			0,
			functions.GetFunctionName()+" - "+err.Error())

		return "", "", err

	}

	for rows.Next() {
		err = rows.Scan(&side1, &side2)
	}

	rows.Close()

	return side1, side2, err

}

// GetOrderSymbol Get symbol for ThreadID
func GetOrderSymbol(
	sessionData *types.Session) (symbol string, err error) {

	var rows *sql.Rows

	if rows, err = sessionData.Db.Query("call cryptopump.GetOrderSymbol(?)",
		sessionData.ThreadID); err != nil {

		functions.Logger(
			nil,
			nil,
			sessionData,
			log.DebugLevel,
			0,
			0,
			0,
			0,
			functions.GetFunctionName()+" - "+err.Error())

		return "", err

	}

	for rows.Next() {
		err = rows.Scan(&symbol)
	}

	rows.Close()

	return symbol, err

}

// GetThreadTransactionDistinct Get Thread Distinct
func GetThreadTransactionDistinct(
	sessionData *types.Session) (threadID string, threadIDSession string, err error) {

	var rows *sql.Rows

	if rows, err = sessionData.Db.Query("call cryptopump.GetThreadTransactionDistinct()"); err != nil {

		functions.Logger(
			nil,
			nil,
			sessionData,
			log.DebugLevel,
			0,
			0,
			0,
			0,
			functions.GetFunctionName()+" - "+err.Error())

		return "", "", err

	}

	for rows.Next() {
		err = rows.Scan(
			&threadID,
			&threadIDSession)

		if functions.LockThreadID(threadID) { /* Create lock for threadID */

			break

		} else {

			threadID = ""
			threadIDSession = ""

		}
	}

	rows.Close()

	return threadID, threadIDSession, err

}

// GetOrderTransactionPending Get 1 order with pending FILLED status
func GetOrderTransactionPending(
	sessionData *types.Session) (orderID int64, symbol string, err error) {

	var rows *sql.Rows

	if rows, err = sessionData.Db.Query("call cryptopump.GetOrderTransactionPending(?)",
		sessionData.ThreadID); err != nil {

		functions.Logger(
			nil,
			nil,
			sessionData,
			log.DebugLevel,
			0,
			0,
			0,
			0,
			functions.GetFunctionName()+" - "+err.Error())

		return 0, "", err

	}

	for rows.Next() {
		err = rows.Scan(
			&orderID,
			&symbol)
	}

	rows.Close()

	return orderID, symbol, err

}

// GetThreadTransactionByPrice function
func GetThreadTransactionByPrice(
	marketData *types.Market,
	sessionData *types.Session) (orderID int, price float64, executedQuantity float64, cummulativeQuoteQty float64, transactTime int64, err error) {

	var rows *sql.Rows

	if rows, err = sessionData.Db.Query("call cryptopump.GetThreadTransactionByPrice(?,?)",
		sessionData.ThreadID,
		marketData.Price); err != nil {

		functions.Logger(
			nil,
			nil,
			sessionData,
			log.DebugLevel,
			0,
			0,
			0,
			0,
			functions.GetFunctionName()+" - "+err.Error())

		return 0, 0, 0, 0, 0, err

	}

	for rows.Next() {
		err = rows.Scan(
			&cummulativeQuoteQty,
			&orderID,
			&price,
			&executedQuantity,
			&transactTime)
	}

	rows.Close()

	return orderID, price, executedQuantity, cummulativeQuoteQty, transactTime, err

}

// GetThreadLastTransaction Return the last 'active' BUY transaction for a Thread
func GetThreadLastTransaction(
	sessionData *types.Session) (orderID int, price float64, executedQuantity float64, cummulativeQuoteQty float64, transactTime int64, err error) {

	var rows *sql.Rows

	if rows, err = sessionData.Db.Query("call cryptopump.GetThreadLastTransaction(?)",
		sessionData.ThreadID); err != nil {

		functions.Logger(
			nil,
			nil,
			sessionData,
			log.DebugLevel,
			0,
			0,
			0,
			0,
			functions.GetFunctionName()+" - "+err.Error())

		return 0, 0, 0, 0, 0, err

	}

	for rows.Next() {
		err = rows.Scan(
			&cummulativeQuoteQty,
			&orderID,
			&price,
			&executedQuantity,
			&transactTime)
	}

	rows.Close()

	return orderID, price, executedQuantity, cummulativeQuoteQty, transactTime, err

}

// GetThreadTransactiontUpmarketPriceCount function
func GetThreadTransactiontUpmarketPriceCount(
	sessionData *types.Session,
	price float64) (count int, err error) {

	var rows *sql.Rows

	if rows, err = sessionData.Db.Query("call cryptopump.GetThreadTransactiontUpmarketPriceCount(?,?)",
		sessionData.ThreadID,
		price); err != nil {

		functions.Logger(
			nil,
			nil,
			sessionData,
			log.DebugLevel,
			0,
			0,
			0,
			0,
			functions.GetFunctionName()+" - "+err.Error())

		return 0, err

	}

	for rows.Next() {
		err = rows.Scan(&count)
	}

	rows.Close()

	return count, err

}

// GetOrderTransactionCount Retrieve transaction count by Side and minutes
func GetOrderTransactionCount(
	sessionData *types.Session,
	side string) (count float64, err error) {

	var rows *sql.Rows

	if rows, err = sessionData.Db.Query("call cryptopump.GetOrderTransactionCount(?,?,?)",
		sessionData.ThreadID,
		side,
		(60 * -1)); err != nil {

		functions.Logger(
			nil,
			nil,
			sessionData,
			log.DebugLevel,
			0,
			0,
			0,
			0,
			functions.GetFunctionName()+" - "+err.Error())

		return 0, err

	}

	for rows.Next() {
		err = rows.Scan(&count)
	}

	rows.Close()

	return count, err

}

// GetThreadTransactionByThreadID  Retrieve transaction count by Side and minutes
func GetThreadTransactionByThreadID(
	sessionData *types.Session) (orders []types.Order, err error) {

	var rows *sql.Rows

	order := types.Order{}

	if rows, err = sessionData.Db.Query("call cryptopump.GetThreadTransactionByThreadID(?)",
		sessionData.ThreadID); err != nil {

		functions.Logger(
			nil,
			nil,
			sessionData,
			log.DebugLevel,
			0,
			0,
			0,
			0,
			functions.GetFunctionName()+" - "+err.Error())

		return nil, err

	}

	for rows.Next() {

		var orderID int
		var cummulativeQuoteQty, price string
		err = rows.Scan(&orderID, &cummulativeQuoteQty, &price)

		order.OrderID = orderID
		order.CummulativeQuoteQuantity = math.Round(functions.StrToFloat64(cummulativeQuoteQty)*100) / 100
		order.Price = math.Round(functions.StrToFloat64(price)*1000) / 1000
		orders = append(orders, order)

	}

	rows.Close()

	return orders, err

}

// GetProfitByThreadID Retrieve thread profit
func GetProfitByThreadID(
	sessionData *types.Session) (profit float64, err error) {

	var rows *sql.Rows

	if rows, err = sessionData.Db.Query("call cryptopump.GetProfitByThreadID(?)",
		sessionData.ThreadID); err != nil {

		functions.Logger(
			nil,
			nil,
			sessionData,
			log.DebugLevel,
			0,
			0,
			0,
			0,
			functions.GetFunctionName()+" - "+err.Error())

		return 0, err

	}

	for rows.Next() {
		err = rows.Scan(&profit)
	}

	rows.Close()

	return math.Round(profit*100) / 100, err

}

// GetProfit Retrieve total profit
func GetProfit(
	sessionData *types.Session) (profit float64, err error) {

	var rows *sql.Rows

	if rows, err = sessionData.Db.Query("call cryptopump.GetProfit()"); err != nil {

		functions.Logger(
			nil,
			nil,
			sessionData,
			log.DebugLevel,
			0,
			0,
			0,
			0,
			functions.GetFunctionName()+" - "+err.Error())

		return 0, err

	}

	for rows.Next() {
		err = rows.Scan(&profit)
	}

	rows.Close()

	return math.Round(profit*100) / 100, err

}

// GetThreadCount Retrieve Running Thread Count
func GetThreadCount(
	sessionData *types.Session) (count int, err error) {

	var rows *sql.Rows

	if rows, err = sessionData.Db.Query("call cryptopump.GetThreadCount()"); err != nil {

		functions.Logger(
			nil,
			nil,
			sessionData,
			log.DebugLevel,
			0,
			0,
			0,
			0,
			functions.GetFunctionName()+" - "+err.Error())

		return 0, err

	}

	for rows.Next() {
		err = rows.Scan(&count)
	}

	rows.Close()

	return count, err

}

// GetThreadAmount Retrieve Thread Dollar Amount
func GetThreadAmount(
	sessionData *types.Session) (amount float64, err error) {

	var rows *sql.Rows

	if rows, err = sessionData.Db.Query("call cryptopump.GetThreadTransactionAmount()"); err != nil {

		functions.Logger(
			nil,
			nil,
			sessionData,
			log.DebugLevel,
			0,
			0,
			0,
			0,
			functions.GetFunctionName()+" - "+err.Error())

		return 0, err

	}

	for rows.Next() {
		err = rows.Scan(&amount)
	}

	rows.Close()

	return math.Round(amount*100) / 100, err

}
