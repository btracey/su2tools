package remove_whitespace

import (
	"bufio"
	"bytes"
	"errors"
	"io"
)

// RemoveWhitespace trims the whitespace off of the ends of the lines
// in r and writes the result to writer. The writer will be written after
// the lines are all scanned. r is usually the contents of a file (from ReadFile)
func RemoveTrailingWhitespace(r []byte, w io.Writer) error {
	if len(r) == 0 {
		return nil
	}
	buf := bytes.NewBuffer(r)
	s := bufio.NewScanner(buf)
	b := make([]byte, 0, 1000)
	for s.Scan() {
		line := s.Bytes()
		b = append(b, bytes.TrimSpace(line)...)
		b = append(b, '\n')
	}
	if s.Err() != nil {
		return errors.New("error removing whitespace: " + s.Err().Error())
	}
	// remove final '\n'
	b = b[:len(b)-1]
	_, err := w.Write(b)
	if err != nil {
		return errors.New("error writing: " + err.Error())
	}
	return nil
}
