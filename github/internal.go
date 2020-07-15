package github

import (
	"os"
	"strings"
)

// WriteLineToFile writes a line to file
func WriteLineToFile(filename, line string) (err error) {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		if er := file.Close(); er != nil {
			err = er
		}
	}()

	if _, err := file.WriteString(line); err != nil {
		return err
	}
	return nil
}

// Replace replaces a string
func Replace(input, from, to string) string {
	return strings.Replace(input, from, to, -1)
}
