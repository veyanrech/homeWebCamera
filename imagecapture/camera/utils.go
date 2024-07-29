package camera

import "time"

func generateFilename(additional string) string {

	timeNow := time.Now().Format("2006-01-02-15:04:05.000")

	return timeNow + "-" + additional + ".jpg"

}
