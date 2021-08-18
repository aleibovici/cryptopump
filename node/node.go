package node

import (
	"os"
	"time"

	"github.com/aleibovici/cryptopump/functions"
	"github.com/aleibovici/cryptopump/logger"
	"github.com/aleibovici/cryptopump/types"
)

// GetRole Define node role Master or Slave
func GetRole(
	configData *types.Config,
	sessionData *types.Session) {

	var filename string = "master.lock"

	/* 	If TestNet is enabled will not check for "master.lock" to not affect production systems */
	if configData.TestNet {

		sessionData.MasterNode = false
		return

	}

	/* If Master Node already set to True */
	if sessionData.MasterNode {

		/* Set access time and modified time of the file to the current time */
		err := os.Chtimes(filename, time.Now().Local(), time.Now().Local())

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

		return

	}

	/* If Master Node set to False */
	if file, err := os.Stat(filename); err == nil { /* Check if "master.lock" is created and modified time */

		sessionData.MasterNode = false

		if time.Duration(time.Since(file.ModTime()).Seconds()) > 100 { /* Remove "master.lock" if old modified time */

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

	} else if os.IsNotExist(err) { /* Check if "master.lock" is created and modified time */

		var file *os.File
		if file, err = os.Create(filename); err != nil {

			logger.LogEntry{
				Config:   nil,
				Market:   nil,
				Session:  sessionData,
				Order:    &types.Order{},
				Message:  functions.GetFunctionName() + " - " + err.Error(),
				LogLevel: "DebugLevel",
			}.Do()

		}

		file.Close()

		sessionData.MasterNode = true

	}

}

// ReleaseRole Release node role if Master
func ReleaseRole(
	sessionData *types.Session) {

	/* Release node role if Master */
	if sessionData.MasterNode {

		var filename string = "master.lock"

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
