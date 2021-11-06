package ggoscmds

import "os"

// LogToFile - log the [msg] to [file].
func LogToFile(file string, msg string) (err error) {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	defer func() {
		if f != nil {
			f.Close()
		}
	}()
	if err != nil {
		return
	}
	if _, err2 := f.WriteString(msg); err2 != nil {
		err = err2
		return
	}
	f.WriteString("\n")

	return
}
