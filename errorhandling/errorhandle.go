package errorhandling

import "log"

// helper function to make one-liner error handling
func HandleError(err error, message ...string) {
	if err != nil {
		errorMsg := err.Error()
		if len(message) > 0 {
			errorMsg = message[0] + ": " + errorMsg
		}
		log.Printf("%s", errorMsg)
	}
}
