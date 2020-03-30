package util

import (
	"fmt"
	"io"
)

// Table of results
type Table struct {
	Title   string    // Table title (optional)
	Style   TableSyle // Style of table
	Headers []*Header // Column headers
	Rows    []*Row    // Response rows
}

// Table Header
type Header struct {
	Label string // Label of header
	Align int    // Column alignment
	Width int    // Width of column (used at write time only)
}

type Row struct {
	parent *Table   // link to parent table
	cells  []string // Cell contents
}

const (
	Left   = iota // Left justify (default)
	Center        // Center justify
	Right         // Right justify
)

// Table styling
type TableSyle struct {
	HSep       byte // Char separating columns in row separator
	HLine      byte // Char forming row separator
	CSep       byte // Char separating columns in rows
	Border     bool // Outer border
	PageHeight int  // Page height, 0=no paging
	SOH        byte // Start of heading i.e. table, 0=ignore
	STX        byte // Start of text i.e. page, 0=ignore
	ETX        byte // End of text i.e. page, 0=ignore
	EOT        byte // End of Transmission i.e. table, 0=ignore
}

// Plain Table Style
var Plain = TableSyle{
	HSep:   '=',
	HLine:  '=',
	CSep:   ' ',
	Border: false,
}

// SQL Table Style
var SQL = TableSyle{
	HSep:   '+',
	HLine:  '-',
	CSep:   '|',
	Border: true,
}

// BBC Mode 8 - i.e. for the BBC Master 128 ROM
var MODE7 = TableSyle{
	HSep:       0,     // no separator
	HLine:      0,     // no separator
	CSep:       ' ',   // Space
	Border:     false, // no border
	PageHeight: 19,    // Mode7 page height minus title & header
	SOH:        1,     // ASCII code
	STX:        2,     // ASCII code
	ETX:        3,     // ASCII code
	EOT:        4,     // ASCII code
}

func (t *TableSyle) WriteCSep(o io.Writer) error {
	if t != nil {
		err := write(o, t.CSep)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *TableSyle) WriteHSep(o io.Writer) error {
	if t != nil && t.Border {
		err := write(o, t.HSep)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *TableSyle) WriteBorder(o io.Writer, tab *Table) error {
	if t.Border {
		return t.WriteSeparator(o, tab)
	}
	return nil
}

func (t *TableSyle) WriteSeparator(o io.Writer, tab *Table) error {
	err := t.WriteHSep(o)
	if err != nil {
		return err
	}

	for i, h := range tab.Headers {
		if i > 0 {
			err := write(o, t.HSep)
			if err != nil {
				return err
			}
		}

		var b []byte
		for i := 0; i < h.Width; i++ {
			b = append(b, t.HLine)
		}
		_, err = o.Write(b)
		if err != nil {
			return err
		}
	}

	err = t.WriteHSep(o)
	if err != nil {
		return err
	}

	return write(o, '\n')
}

func (t *TableSyle) VisitRow(o io.Writer, v func(o io.Writer) error) error {
	err := t.WriteCSep(o)
	if err != nil {
		return err
	}

	err = v(o)
	if err != nil {
		return err
	}

	err = t.WriteCSep(o)
	if err != nil {
		return err
	}

	return write(o, '\n')
}

func (t *TableSyle) VisitCell(o io.Writer, i int, v func(o io.Writer) error) error {
	if i > 0 {
		err := t.WriteCSep(o)
		if err != nil {
			return err
		}
	}

	return v(o)
}

func (t *Table) AddHeader(label string) *Table {
	t.Headers = append(t.Headers, &Header{Label: label})
	return t
}

func (t *Table) AddHeaders(label ...string) *Table {
	for _, l := range label {
		t.AddHeader(l)
	}
	return t
}

func (t *Table) NewRow() *Row {
	r := &Row{parent: t}
	t.Rows = append(t.Rows, r)
	return r
}

func (r *Row) NewRow() *Row {
	return r.End().NewRow()
}

func (r *Row) End() *Table {
	return r.parent
}

func (r *Row) Append(v string) *Row {
	r.cells = append(r.cells, v)
	return r
}

func (r *Row) AppendInt(v int) *Row {
	return r.Append(fmt.Sprintf("%d", v))
}

func (r *Row) Appendf(f string, v ...interface{}) *Row {
	return r.Append(fmt.Sprintf(f, v...))
}

func (r *Row) Visit(f func(i int, h *Header, c string) error) error {
	t := r.parent
	for i, c := range r.cells {
		if i >= len(t.Headers) {
			t.AddHeader("")
		}

		h := t.Headers[i]
		err := f(i, h, c)
		if err != nil {
			return err
		}
	}
	return nil
}

func write(o io.Writer, v byte) error {
	_, err := o.Write([]byte{v})
	return err
}

func (h *Header) append(o io.Writer, v string) error {
	var f string
	switch h.Align {
	case Left:
		f = "%%-%d.%ds"
	case Center:
		f = "%%-%d.%ds"
	case Right:
		f = "%%%d.%ds"
	}

	_, err := fmt.Fprintf(o, fmt.Sprintf(f, h.Width, h.Width), v)
	return err
}

func (t *Table) Write(out io.Writer) error {
	// init column width
	for _, h := range t.Headers {
		h.Width = len(h.Label)
	}

	// Now get max width
	for _, r := range t.Rows {
		_ = r.Visit(func(i int, h *Header, c string) error {
			l := len(c)
			if l > 0 {
				h := t.Headers[i]
				if h.Width < l {
					h.Width = l
				}
			}
			return nil
		})
	}

	err := t.Style.WriteBorder(out, t)
	if err != nil {
		return err
	}

	err = t.Style.VisitRow(out, func(o io.Writer) error {
		for i, h := range t.Headers {
			err := t.Style.VisitCell(o, i, func(o io.Writer) error {
				return h.append(o, h.Label)
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	err = t.Style.WriteSeparator(out, t)
	if err != nil {
		return err
	}

	for _, r := range t.Rows {
		err = t.Style.VisitRow(out, func(o io.Writer) error {
			return r.Visit(func(i int, h *Header, c string) error {
				return t.Style.VisitCell(o, i, func(o io.Writer) error {
					return h.append(o, c)
				})
			})
		})
		if err != nil {
			return err
		}
	}

	err = t.Style.WriteBorder(out, t)
	if err != nil {
		return err
	}

	return write(out, '\n')
}
