package nodes

import (
	"os"
	"time"

	"github.com/aleibovici/cryptopump/functions"
	"github.com/aleibovici/cryptopump/logger"
	"github.com/aleibovici/cryptopump/mysql"
	"github.com/aleibovici/cryptopump/types"
)

// Node functions
type Node struct{}

// GetRole retrieve correct role for node
func (Node) GetRole(
	configData *types.Config,
	sessionData *types.Session) {

	var filename = "master.lock"
	var err error

	/* Conditional defer logging when there is an error retriving data */
	defer func() {
		if err != nil {
			logger.LogEntry{
				Config:   nil,
				Market:   nil,
				Session:  sessionData,
				Order:    &types.Order{},
				Message:  functions.GetFunctionName() + " - " + err.Error(),
				LogLevel: "DebugLevel",
			}.Do()
		}
	}()

	/* 	If TestNet is enabled will not check for "master.lock" to not affect production systems */
	if configData.TestNet {

		sessionData.MasterNode = false
		return

	}

	/* If Master Node already set to True */
	if sessionData.MasterNode {

		/* Set access time and modified time of the file to the current time */
		if err = os.Chtimes(filename, time.Now().Local(), time.Now().Local()); err != nil {

			return

		}

		return

	}

	/* If Master Node set to False */
	if file, err := os.Stat(filename); err == nil { /* Check if "master.lock" is created and modified time */

		sessionData.MasterNode = false

		if time.Duration(time.Since(file.ModTime()).Seconds()) > 100 { /* Remove "master.lock" if old modified time */

			if err := os.Remove(filename); err != nil {

				return

			}

		}

	} else if os.IsNotExist(err) { /* If "master.lock" is not created create it */

		var file *os.File
		if file, err = os.Create(filename); err != nil { /* Create "master.lock" if not exists */

			file.Close()                  /* Close file */
			sessionData.MasterNode = true /* Set Master Node to True */

		} else {

			file.Close() /* Close file */

		}

	}

}

// ReleaseMasterRole Release node role if Master
func (Node) ReleaseMasterRole(sessionData *types.Session) {

	/* Release node role if Master */
	if sessionData.MasterNode {

		var filename = "master.lock"

		if err := os.Remove(filename); err != nil {

			logger.LogEntry{
				Config:   nil,
				Market:   nil,
				Session:  sessionData,
				Order:    &types.Order{},
				Message:  functions.GetFunctionName() + " - " + err.Error(),
				LogLevel: "DebugLevel",
			}.Do()

		}

	}

}

// CheckStatus check for errors on node
func (Node) CheckStatus(configData *types.Config,
	sessionData *types.Session) {

	/* Check last WsBookTicker */
	if time.Duration(time.Since(sessionData.LastWsBookTickerTime).Seconds()) > time.Duration(30) {

		sessionData.Status = true

	}

	/* Check last WsKline */
	if time.Duration(time.Since(sessionData.LastWsKlineTime).Seconds()) > time.Duration(100) {

		sessionData.Status = true

	}

	/* Update Session table */
	if err := mysql.UpdateSession(
		configData,
		sessionData); err != nil {

		logger.LogEntry{
			Config:   configData,
			Market:   nil,
			Session:  sessionData,
			Order:    &types.Order{},
			Message:  functions.GetFunctionName() + " - " + err.Error(),
			LogLevel: "DebugLevel",
		}.Do()

	}

	sessionData.Status = false
}
