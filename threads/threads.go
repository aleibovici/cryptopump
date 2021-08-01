package threads

import (
	"cryptopump/functions"
	"cryptopump/logger"
	"cryptopump/mysql"
	"cryptopump/node"
	"cryptopump/types"
	"os"
	"time"
)

// ExitThreadID Cleanly exit a Thread
func ExitThreadID(
	sessionData *types.Session) {

	/* Verify wether buying/selling to allow graceful session exit */
	for sessionData.Busy {
		time.Sleep(time.Millisecond * 200)
	}

	/* Release node role if Master */
	if sessionData.MasterNode {

		node.ReleaseRole(sessionData)

	}

	/* Remove lock for threadID */
	unlockThreadID(sessionData)

	/* Delete session from Session table */
	_ = mysql.DeleteSession(sessionData)

	logger.LogEntry{
		Config:   nil,
		Market:   nil,
		Session:  sessionData,
		Order:    &types.Order{},
		Message:  "Clean Shutdown",
		LogLevel: "InfoLevel",
	}.Do()

	os.Exit(1)

}

/* Remove lock for threadID */
func unlockThreadID(
	sessionData *types.Session) {

	filename := sessionData.ThreadID + ".lock"

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
