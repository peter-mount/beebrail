package util

// Pagination handles pages, mainly for Table's
type Pagination struct {
	pageNo     int // Current page number from 1
	pageCount  int // Total number of pages
	pageHeight int // number of rows per page
	pageWidth  int // page width
	numRows    int // Number of rows
}

type Page struct {
	PageNo     int // Current page number from 1
	PageCount  int // Total number of pages
	PageHeight int // number of rows per page
	PageWidth  int // page width
}

// PaginationCallback allows customisation of the output
type PaginationCallback struct {
	PageHeader  func(p *Pagination, o *ResultWriter) error // Write optional page header
	TableHeader func(o *ResultWriter) error                // Write the TableHeader
}

func (t *Table) NewPagination() *Pagination {
	p := &Pagination{}
	return p.AddPages(t)
}

func (p *Pagination) Page() Page {
	return Page{
		PageNo:     p.pageNo,
		PageCount:  p.pageCount,
		PageHeight: p.pageHeight,
		PageWidth:  p.pageWidth,
	}
}

func (p *Pagination) IncPages(pageCount int) *Pagination {
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
		// FIXME this breaks so we currently show an incorrect page count
		//if (p.pageHeight * pageCount) > numRows {
		//  pageCount = pageCount - 1
		//}
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
