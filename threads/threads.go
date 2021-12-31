package threads

import (
	"os"
	"time"

	"github.com/aleibovici/cryptopump/functions"
	"github.com/aleibovici/cryptopump/logger"
	"github.com/aleibovici/cryptopump/mysql"
	"github.com/aleibovici/cryptopump/nodes"
	"github.com/aleibovici/cryptopump/types"
)

// Thread locking control
type Thread struct{}

// Terminate thread
func (Thread) Terminate(sessionData *types.Session, message string) {

	if message != "" {

		logger.LogEntry{ /* Log Entry */
			Config:   nil,
			Market:   nil,
			Session:  sessionData,
			Order:    &types.Order{},
			Message:  message,
			LogLevel: "DebugLevel",
		}.Do()

	}

	/* Verify wether buying/selling to allow graceful session exit */
	for sessionData.Busy {

		time.Sleep(time.Millisecond * 200)

	}

	/* Release node role if Master */
	if sessionData.MasterNode {

		nodes.Node{}.ReleaseMasterRole(sessionData)

	}

	// Unlock existing thread
	Thread{}.Unlock(sessionData)

	/* Delete session from Session table */
	if err := mysql.DeleteSession(sessionData); err != nil {

		logger.LogEntry{ /* Log Entry */
			Config:   nil,
			Market:   nil,
			Session:  sessionData,
			Order:    &types.Order{},
			Message:  "Clean Shutdown Failed",
			LogLevel: "DebugLevel",
		}.Do()

	} else {

		logger.LogEntry{ /* Log Entry */
			Config:   nil,
			Market:   nil,
			Session:  sessionData,
			Order:    &types.Order{},
			Message:  "Clean Shutdown",
			LogLevel: "InfoLevel",
		}.Do()

	}

	os.Exit(1)

}

// Lock existing thread
func (Thread) Lock(sessionData *types.Session) bool {

	if sessionData.ThreadID == "" {

		return false

	}

	filename := sessionData.ThreadID + ".lock"

	if _, err := os.Stat(filename); err == nil {

		return false

	} else if os.IsNotExist(err) {

		var file, err = os.Create(filename)

		if err != nil {

			return false

		}

		file.Close()

		return true

	}

	return false

} // //// // ExitThreadID Cleanly exit a Thread

// Unlock existing thread
func (Thread) Unlock(sessionData *types.Session) {

	if err := os.Remove(sessionData.ThreadID + ".lock"); err != nil {

		logger.LogEntry{ /* Log Entry */
			Config:   nil,
			Market:   nil,
			Session:  sessionData,
			Order:    &types.Order{},
			Message:  functions.GetFunctionName() + " - " + err.Error(),
			LogLevel: "DebugLevel",
		}.Do()

	}

}
