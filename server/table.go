package server

// A table that is rendered correctly on the BBC
type Table struct {
	Headers []string
	Rows    []*Row
}

type Row struct {
	Cells []string
}

func (t *Table) AddHeader(h string) *Table {
	t.Headers = append(t.Headers, h)
	return t
}

func (t *Table) AddRow() *Row {
	row := &Row{}
	t.Rows = append(t.Rows, row)
	return row
}

func (r *Row) Append(s string) *Row {
	r.Cells = append(r.Cells, s)
	return r
}

// Table write format
// rowCount Number of rows
// colCount Number of columns
// header row
// 0..n data rows
//
// Each row consists of colCount 0 terminated strings
//
func (p *Packet) AppendTable(t *Table) *Packet {

	// Get max col count
	cols := 0
	for _, r := range t.Rows {
		c := len(r.Cells)
		if c > cols {
			cols = c
		}
	}

	p.Append(
		byte(len(t.Rows)), // Row length
		byte(cols),        // Col count
	)

	p.AppendCStrings(cols, t.Headers)

	for _, r := range t.Rows {
		p.AppendCStrings(cols, r.Cells)
	}

	return p
}
