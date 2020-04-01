package util

import (
	"fmt"
	"strings"
)

// Table of results
type Table struct {
	Style      TableStyle         // Style of table
	Headers    []*Header          // Column headers
	Rows       []*Row             // Response rows
	pagination *Pagination        // Page pagination
	linked     bool               // True if this Table was created from another
	nextTable  *Table             // next table in the chain
	Callback   PaginationCallback // Callbacks
}

// Table Header
type Header struct {
	Label  string // Label of header
	Align  int    // Column alignment
	Width  int    // Width of column (used at write time only)
	Prefix string // Optional prefix
}

type Row struct {
	parent *Table  // link to parent table
	cells  []*Cell // Cell contents
}

type Cell struct {
	Label  string // Cell contents
	Prefix string // Optional prefix
}

const (
	Left   = iota // Left justify (default)
	Center        // Center justify
	Right         // Right justify
	None          // No justification
)

// Table styling
type TableStyle struct {
	ShowTitle  bool // Show table title
	HSep       byte // Char separating columns in row separator
	HLine      byte // Char forming row separator
	Separator  bool // true to enable the separator bar
	CSep       byte // Char separating columns in rows
	HCSep      byte // Char separating columns in headers, normally same as CSep
	Border     bool // Outer border
	PageHeight int  // Page height, 0=no paging
	Mode7      bool // true for BBC mode 7 Teletext
}

// Plain Table Style
var Plain = TableStyle{
	HSep:      '=',
	HLine:     '=',
	Separator: true,
	HCSep:     ' ',
	CSep:      ' ',
	Border:    false,
}

// SQL Table Style
var SQL = TableStyle{
	HSep:      '+',
	HLine:     '-',
	Separator: true,
	HCSep:     '|',
	CSep:      '|',
	Border:    true,
}

// BBC Mode 8 - i.e. for the BBC Master 128 ROM
var MODE7 = TableStyle{
	ShowTitle:  true, // Show table title
	HSep:       0,    // no separator
	HLine:      0,    // no separator
	Separator:  false,
	CSep:       0, // no column separator as we replace it
	HCSep:      ' ',
	Border:     false, // no border
	PageHeight: 19,    // Mode7 page height minus title & header
	Mode7:      true,  // Flag as Mode7
}

// NewTable creates a new table using this one's configuration
func (t *Table) NewTable() *Table {
	t1 := &Table{
		Style:      t.Style,        // Duplicate the style
		pagination: t.Pagination(), // use the same instance as previous Table
		linked:     true,           // mark as linked
		Callback:   t.Callback,     // Duplicate callbacks, this can be changed later
	}

	// Link to old table
	t.nextTable = t1

	return t1
}

func (t *TableStyle) WriteCSep(o *ResultWriter) error {
	if t != nil && t.Border {
		err := write(o, t.CSep)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *TableStyle) WriteHCSep(o *ResultWriter) error {
	if t != nil && t.Border {
		err := write(o, t.HCSep)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *TableStyle) WriteHSep(o *ResultWriter) error {
	if t != nil && t.Border {
		err := write(o, t.HSep)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *TableStyle) WriteBorder(o *ResultWriter, tab *Table) error {
	if t.Border {
		return t.WriteSeparator(o, tab)
	}
	return nil
}

func (t *TableStyle) WriteSeparator(o *ResultWriter, tab *Table) error {
	if t.Separator {
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

	return nil
}

func (t *TableStyle) VisitHeaders(o *ResultWriter, v func(o *ResultWriter) error) error {
	err := t.WriteHCSep(o)
	if err != nil {
		return err
	}

	err = v(o)
	if err != nil {
		return err
	}

	err = t.WriteHCSep(o)
	if err != nil {
		return err
	}

	return write(o, '\n')
}

func (t *TableStyle) VisitHeader(o *ResultWriter, i int, v func(o *ResultWriter) error) error {
	if i > 0 && t.HCSep > 0 {
		err := write(o, t.HCSep)
		if err != nil {
			return err
		}
	}

	return v(o)
}

func (t *TableStyle) VisitRow(tab *Table, r *Row, o *ResultWriter, v func(t *Table, r *Row, o *ResultWriter) error) error {
	err := t.WriteCSep(o)
	if err != nil {
		return err
	}

	err = v(tab, r, o)
	if err != nil {
		return err
	}

	err = t.WriteCSep(o)
	if err != nil {
		return err
	}

	return write(o, '\n')
}

func (t *TableStyle) VisitCell(o *ResultWriter, i int, v func(o *ResultWriter) error) error {
	if i > 0 && t.CSep > 0 {
		err := write(o, t.CSep)
		if err != nil {
			return err
		}
	}

	return v(o)
}

func (t *Table) AddHeader(label string) *Table {
	return t.Header(&Header{Label: label})
}

func (t *Table) Header(h *Header) *Table {
	t.Headers = append(t.Headers, h)
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

func (r *Row) IsEmpty() bool {
	return len(r.cells) == 0
}

func (r *Row) GetCell(i int) *Cell {
	return r.cells[i]
}

func (r *Row) Cell(v *Cell) *Row {
	r.cells = append(r.cells, v)
	return r
}

func (r *Row) Append(v string) *Row {
	return r.Cell(&Cell{Label: v})
}

func (r *Row) AppendInt(v int) *Row {
	return r.Append(fmt.Sprintf("%d", v))
}

func (r *Row) Appendf(f string, v ...interface{}) *Row {
	return r.Append(fmt.Sprintf(f, v...))
}

func (r *Row) Visit(f func(i int, h *Header, c *Cell) error) error {
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

func write(o *ResultWriter, v byte) error {
	_, err := o.Write([]byte{v})
	return err
}

func (h *Header) append(o *ResultWriter, v string) error {
	var f string
	switch h.Align {
	case Left:
		f = "%%-%d.%ds"
	case Center:
		f = "%%-%d.%ds"
	case Right:
		f = "%%%d.%ds"
	case None:
		return o.WriteString(v)
	}

	_, err := fmt.Fprintf(o, fmt.Sprintf(f, h.Width, h.Width), v)

	return err
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func (t *Table) layout() {

	// Ensure we have a Pagination instance
	if t.linked {
		// Add out pages to the output
		t.pagination.AddPages(t)
	} else {
		_ = t.Pagination()
	}

	// init column width
	for _, h := range t.Headers {
		h.Width = max(h.Width, len(h.Label))
	}

	// Now get max width
	for _, r := range t.Rows {
		_ = r.Visit(func(i int, h *Header, c *Cell) error {
			for _, v := range strings.Split(c.Label, "\n") {
				l := len(v)
				if l > 0 {
					h := t.Headers[i]
					h.Width = max(h.Width, l)
				}
			}
			return nil
		})
	}

	// Finally the page width
	w := 0
	for _, h := range t.Headers {
		w = w + h.Width
	}
	if w > t.pagination.pageWidth {
		t.pagination.pageWidth = w
	}

	// layout the next table
	if t.nextTable != nil {
		t.nextTable.layout()
	}

}

func (t *Table) Write(out *ResultWriter) error {
	t.layout()
	return t.write(out)
}

func (t *Table) write(out *ResultWriter) error {
	// Ugly but need this tables pageHeight here
	t.pagination.pageHeight = t.Style.PageHeight

	err := t.writeTable(out)
	if err != nil {
		return err
	}

	err = t.Style.WriteBorder(out, t)
	if err != nil {
		return err
	}

	if t.nextTable != nil {
		t.pagination.Reset()
		return t.nextTable.write(out)
	}

	return nil
}

func (t *Table) WriteHeader(out *ResultWriter) error {
	err := t.Style.WriteBorder(out, t)
	if err != nil {
		return err
	}

	err = t.Style.VisitHeaders(out, func(o *ResultWriter) error {
		for i, h := range t.Headers {
			err := t.Style.VisitHeader(o, i, func(o *ResultWriter) error {
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

	return t.Style.WriteSeparator(out, t)
}

func (t *Table) Pagination() *Pagination {
	if t.pagination == nil {
		t.pagination = t.NewPagination()
	}
	return t.pagination
}

func (t *Table) writeTable(out *ResultWriter) error {

	for rowNum, r := range t.Rows {
		if t.pagination.IsPageBreak() {
			// Start new page
			if t.pagination.NextPage() {
				if rowNum > 0 {
					// Page 2 up add a record separator
					err := out.RecordSeparator()
					if err != nil {
						return err
					}
				} else if t.linked {
					// Page 1 but linked then add a Group separator
					err := out.GroupSeparator()
					if err != nil {
						return err
					}
				}
			}
			t.pagination.NextRow()

			if t.Callback.PageHeader != nil {
				err := t.Callback.PageHeader(t.pagination, out)
				if err != nil {
					return err
				}
			}

			if t.Callback.TableHeader != nil {
				err := t.Callback.TableHeader(out)
				if err != nil {
					return err
				}
			}
		}

		if t.Callback.TableRow != nil {
			err := t.Callback.TableRow(t, r, out)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (t *Table) WriteRow(_ *Table, r *Row, o *ResultWriter) error {
	return t.Style.VisitRow(t, r, o, func(t *Table, r *Row, o *ResultWriter) error {
		return r.Visit(func(i int, h *Header, c *Cell) error {
			return t.Style.VisitCell(o, i, func(o *ResultWriter) error {
				// FIXME this only works for multi line for tables of 1 column
				for i, s := range strings.Split(c.Label, "\n") {
					if i > 0 {
						err := o.WriteBytes('\n')
						if err != nil {
							return err
						}
					}

					p := c.Prefix
					if p == "" && h.Prefix != "" {
						p = h.Prefix
					}
					if p != "" {
						err := o.WriteString(p)
						if err != nil {
							return err
						}
					}

					err := h.append(o, s)
					if err != nil {
						return err
					}
				}
				return nil
			})
		})
	})
}
