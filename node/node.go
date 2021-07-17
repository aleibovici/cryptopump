package node

import (
	"cryptopump/functions"
	"cryptopump/types"
	"os"

	log "github.com/sirupsen/logrus"
)

/* Define node role Master or Slave */
func GetRole(
	sessionData *types.Session) {

	var file *os.File
	var filename string = "master.lock"

	/* Exit function if Master Node already Master */
	if sessionData.MasterNode {

		return

	}

	/* Check if "master.lock" is created */
	/* Create "master.lock" is it doesn't exist */
	if _, err := os.Stat(filename); err == nil {

		sessionData.MasterNode = false

	} else if os.IsNotExist(err) {

		if file, err = os.Create(filename); err != nil {

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

		}

		file.Close()

		sessionData.MasterNode = true

	}

}

/* Release node role if Master */
func ReleaseRole(
	sessionData *types.Session) {

	/* Release node role if Master */
	if sessionData.MasterNode {

		var filename string = "master.lock"

		if err := os.Remove(filename); err != nil {

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

		}

	}

}
