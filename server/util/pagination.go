package util

import "log"

// Pagination handles pages, mainly for Table's
type Pagination struct {
	pageNo     int // Current page number from 1
	pageCount  int // Total number of pages
	pageHeight int // number of rows per page
	numRows    int // Number of rows
}

// PaginationCallback allows customisation of the output
type PaginationCallback struct {
	PageHeader  func(p *Pagination, o *ResultWriter) error // Write optional page header
	TableHeader func(o *ResultWriter) error                // Write the TableHeader
}

func (t *Table) NewPagination() *Pagination {
	p := &Pagination{}
	// Default callbacks
	t.Callback.PageHeader = func(p *Pagination, o *ResultWriter) error { return nil }
	t.Callback.TableHeader = t.WriteHeader

	return p.AddPages(t)
}

func (p *Pagination) PageNo() (int, int) {
	return p.pageNo, p.pageCount
}

func (p *Pagination) IncPages(pageCount int) *Pagination {
	log.Println(p.pageCount, pageCount, p.pageCount+pageCount)
	p.pageCount = p.pageCount + pageCount
	return p
}

func (p *Pagination) AddPages(t *Table) *Pagination {
	p.pageHeight = t.Style.PageHeight

	numRows := len(t.Rows)

	// Work out how many pages we have
	var pageCount int
	if p.pageHeight > 0 {
		// Number of pages, must be >1
		pageCount = 1 + (numRows / p.pageHeight)

		// if we have exactly the right number of rows to fill the last page don't have a blank one
		if (p.pageHeight * pageCount) > numRows {
			pageCount--
		}
	} else {
		// Just one big page for the entire table
		pageCount = 1
	}

	return p.IncPages(pageCount)
}

// IsPageBreak returns true if we should break the table for a specific rowNum
func (p *Pagination) IsPageBreak(rowNum int) bool {
	return rowNum == 0 || (p.pageHeight > 0 && (rowNum%p.pageHeight) == 0)
}

func (p *Pagination) NextPage() bool {
	p.pageNo++
	return p.pageNo > 1
}
