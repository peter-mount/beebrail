package util

import (
	"fmt"
)

// ResultWriter implements a WriteClosable which will wrap it's content within a block suitable for the API
type ResultWriter struct {
	level       int            // API level
	stxEtx      bool           // true to wrap result in stx,etx also enables GS between groups
	title       string         // Title string
	footer      string         // footer string
	w           ResponseWriter // Output writer
	initialised bool           // true once inside document
}

func NewResultWriter(w ResponseWriter) *ResultWriter {
	return &ResultWriter{
		level: ERR_RESULT,
		w:     w,
	}
}

// Level sets the API level (default is ERR_RESULT)
func (r *ResultWriter) Level(level int) *ResultWriter {
	r.level = level
	return r
}

// StxEtx enables/disables wrapping the result within an STX/ETX pair. It also enables GS (Group Separator) for between
// pages
func (r *ResultWriter) StxEtx(stxEtx bool) *ResultWriter {
	r.stxEtx = stxEtx
	return r
}

// Sets the result title. This is sent with the initial level string
func (r *ResultWriter) Title(format string, args ...interface{}) *ResultWriter {
	r.title = fmt.Sprintf(format, args...)
	return r
}

func (r *ResultWriter) Footer(format string, args ...interface{}) *ResultWriter {
	r.footer = fmt.Sprintf(format, args...)
	return r
}

func (r *ResultWriter) WriteBytes(b ...byte) error {
	_, err := r.Write(b)
	return err
}

func (r *ResultWriter) WriteString(s string) error {
	_, err := r.Write([]byte(s[:]))
	return err
}

func (r *ResultWriter) Write(b []byte) (int, error) {
	if !r.initialised {
		// title
		r.w.Println(r.level, r.title)

		// STX
		if r.stxEtx {
			_, err := r.w.Write(stx)
			if err != nil {
				return 0, err
			}
		}

		r.initialised = true
	}

	n, err := r.w.Write(b)
	return n, err
}

var (
	stx = []byte{2, 13, 10}  // ASCII STX Start TeXt - mark start of result
	etx = []byte{3, 13, 10}  // ASCII ETX End TeXT - mark end of result
	gs  = []byte{29, 13, 10} // ASCII Group Separator - mark new table
	rs  = []byte{30, 13, 10} // ASCII Record Separator - mark new page in table
)

func (r *ResultWriter) GroupSeparator() error {
	if r.stxEtx {
		_, err := r.w.Write(gs)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ResultWriter) RecordSeparator() error {
	if r.stxEtx {
		_, err := r.w.Write(rs)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ResultWriter) Close() error {
	// ETX - also write a new line to terminate the output before the end result
	if r.stxEtx {
		_, err := r.w.Write(etx)
		if err != nil {
			return err
		}
	}

	if r.footer != "" {
		// Footer text
		r.w.Println(ERR_INFORMATIONAL, r.footer)
	}

	r.initialised = false
	return nil
}
