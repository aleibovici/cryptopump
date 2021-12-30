package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/aleibovici/cryptopump/types"

	log "github.com/sirupsen/logrus"
)

// LogEntry struct
type LogEntry struct {
	Config   *types.Config  /* Config struct */
	Market   *types.Market  /* Market struct */
	Session  *types.Session /* Session struct */
	Order    *types.Order   /* Order struct */
	Message  string         /* Error message */
	LogLevel string         /* Logrus log level */
}

/* Level define log filename and the log level for the log entry */
func (logEntry LogEntry) level() (filename string) {

	/* Define the log level for the entry */
	switch strings.ToLower(logEntry.LogLevel) {
	case "infolevel":
		log.SetLevel(log.InfoLevel)
		filename = "cryptopump.log"
	case "debuglevel":
		log.SetLevel(log.DebugLevel)
		filename = "cryptopump_debug.log"
	default:
		log.SetLevel(log.DebugLevel)
		filename = "cryptopump_debug.log"
	}

	return filename

}

/* Set the log formatter */
func (logEntry LogEntry) formatter() {

	/* Log as JSON instead of the default ASCII formatter */
	log.SetFormatter(&log.TextFormatter{
		DisableColors:   false,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		DisableSorting:  false,
	})

}

// Do is LogEntry method to run system logging
func (logEntry LogEntry) Do() {

	var err error
	var file *os.File

	logEntry.formatter()         /* Set the log formatter */
	filename := logEntry.level() /* Define the log level for the entry */

	/* io.Writer output set for file */
	if file, err = os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666); err != nil { /* Open the file */

		log.Fatal(err) /* Log the error */

	}

	log.SetOutput(file) /* Set the output for the logger */

	switch {
	case log.StandardLogger().GetLevel() == log.InfoLevel:

		switch logEntry.Message {
		case "UP", "DOWN", "INIT":

			log.WithFields(log.Fields{
				"threadID":  logEntry.Session.ThreadID,
				"rsi3":      fmt.Sprintf("%.2f", logEntry.Market.Rsi3),
				"rsi7":      fmt.Sprintf("%.2f", logEntry.Market.Rsi7),
				"rsi14":     fmt.Sprintf("%.2f", logEntry.Market.Rsi14),
				"MACD":      fmt.Sprintf("%.2f", logEntry.Market.MACD),
				"high":      logEntry.Market.PriceChangeStatsHighPrice,
				"direction": logEntry.Market.Direction,
			}).Info(logEntry.Message)

		case "BUY":

			log.WithFields(log.Fields{
				"threadID":   logEntry.Session.ThreadID,
				"orderID":    logEntry.Order.OrderID,
				"orderPrice": fmt.Sprintf("%.4f", logEntry.Order.Price),
			}).Info(logEntry.Message)

		case "BUYDRYRUN":

			log.WithFields(log.Fields{
				"threadID":   logEntry.Session.ThreadID,
				"orderPrice": fmt.Sprintf("%.4f", logEntry.Order.Price),
			}).Info(logEntry.Message)

		case "SELL":

			log.WithFields(log.Fields{
				"threadID":      logEntry.Session.ThreadID,
				"OrderIDSource": logEntry.Order.OrderIDSource,
				"orderID":       logEntry.Order.OrderID,
				"orderPrice":    fmt.Sprintf("%.4f", logEntry.Order.Price),
			}).Info(logEntry.Message)

		case "SELLDRYRUN":

			log.WithFields(log.Fields{
				"threadID":   logEntry.Session.ThreadID,
				"orderPrice": fmt.Sprintf("%.4f", logEntry.Order.Price),
			}).Info(logEntry.Message)

		case "CANCELED":

			if logEntry.Config.Debug {

				log.WithFields(log.Fields{
					"threadID":      logEntry.Session.ThreadID,
					"OrderIDSource": logEntry.Order.OrderIDSource,
					"orderID":       logEntry.Order.OrderID,
				}).Info(logEntry.Message)

			}

		case "STOPLOSS":

			log.WithFields(log.Fields{
				"threadID": logEntry.Session.ThreadID,
				"orderID":  logEntry.Order.OrderID,
			}).Info(logEntry.Message)

		default:

			log.WithFields(log.Fields{
				"threadID": logEntry.Session.ThreadID,
			}).Info(logEntry.Message)

		}

	case log.StandardLogger().GetLevel() == log.DebugLevel:

		log.WithFields(log.Fields{
			"threadID": logEntry.Session.ThreadID,
			"orderID":  logEntry.Order.OrderID,
		}).Debug(logEntry.Message)

	}

}
