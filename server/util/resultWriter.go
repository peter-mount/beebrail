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

func (r *ResultWriter) Write(b []byte) (int, error) {
	if !r.initialised {
		// title
		r.w.Println(r.level, r.title)

		// STX
		if r.stxEtx {
			_, err := r.w.Write([]byte{2})
			if err != nil {
				return 0, err
			}
		}

		r.initialised = true
	}

	n, err := r.w.Write(b)
	return n, err
}

func (r *ResultWriter) GroupSeparator() error {
	if r.stxEtx {
		_, err := r.w.Write([]byte{29})
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ResultWriter) RecordSeparator() error {
	if r.stxEtx {
		_, err := r.w.Write([]byte{30})
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ResultWriter) Close() error {
	// ETX - also write a new line to terminate the output before the end result
	if r.stxEtx {
		_, err := r.w.Write([]byte{3, 13, 10})
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
