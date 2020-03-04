package server

// The Mode7 Viewdata Teletext control coldes
const (
	AlphaRed = iota + 129
	AlphaGreen
	AlphaYellow
	AlphaBlue
	AlphaMagenta
	AlphaCyan
	AlphaWhite
	Flash
	Steady
	_ // 138 undefined
	_ // 139 undefined
	NormalHeight
	DoubleHeight
	_ // 142 undefined
	_ // 143 undefined
	_ // 144 undefined
	GraphRed
	GraphGreen
	GraphYellow
	GraphBlue
	GraphMagenta
	GraphCyan
	GraphWhite
	Conceal
	ContiguousGraphics
	SeparatedGraphics
	_ // 155 undefined
	BlackBackground
	NewBackground
	HoldGraphics
	ReleaseGraphics
)

// Pages is a collection of Mode7 Pages containing results
type Pages struct {
	Pages       []*Page     // The list of pages
	CurrentPage *Page       // The current page
	newpage     func(*Page) // hook called when a new page is added - think headers
}

func (p *Packet) AppendPages(r *Pages) *Packet {

	// Create a slice of each page's content.
	// We need this so we know the size of each page
	var pages [][]byte

	for _, page := range r.Pages {
		var buf []byte
		for ln, line := range page.Data {
			// skip line 0 as thats reserved for the UI
			if ln > 0 {
				if line == nil || len(line) == 0 {
					buf = append(buf, 13, 10)
				} else {
					buf = append(buf, line...)
					if len(line) < 40 {
						buf = append(buf, 13, 10)
					}
				}
			}
		}
		pages = append(pages, buf)
	}

	// Now the output
	//
	// This consists of the number of pages, a list of offsets to each page then the page data
	//
	// e.g. for a 2 page response
	//
	// Offs Len	Content
	// 		0	1		Number of pages = 2
	//	  1	2		Offset to page 0
	//		3 2		Offset to page 1
	//		5 2		Offset to end of all pages
	//		6	n		Start of page 0, value in 1,2 = 0006

	// So if we want page 0 then we look at bytes 1&2 for the start & output until we hit offset in bytes 3,4
	// For page 1 then its bytes 3&4 to 5&6
	//

	// Offset from begining to the first page, i.e. the page count & the offsets for pages + end offset
	base := 1 + (2 * len(pages)) + 2

	// number of pages max 127
	p.Append(byte(len(pages)))

	// Now the offsets to each page
	for _, pg := range pages {
		p.AppendInt16(base)
		base = base + len(pg)
	}
	// Final offset is for the last page
	p.AppendInt16(base)

	// Now the page content
	for _, pg := range pages {
		p.Append(pg...)
	}

	return p
}

func NewResult(newpage func(*Page)) *Pages {
	r := &Pages{newpage: newpage}
	r.AddPage()
	return r
}

type Page struct {
	Data [][]uint8 // 40x25 character display
	r    *Pages    // parent result
	x    int       // x position in screen
	y    int       // y position in screen
}

// Tab moves the cursor to a new location
// Note: X is 0..39 outside that range x is not changed
//       Y is 0..24 outside that range y is not changed
func (p *Page) Tab(x, y int) *Page {
	if x >= 0 && x <= 39 {
		p.x = x
	}
	if y >= 0 && y <= 24 {
		p.y = y
	}
	return p
}

func (p *Page) X() int {
	return p.x
}

func (p *Page) Y() int {
	return p.y
}

// Add a new page
func (r *Pages) AddPage() *Page {
	p := &Page{
		Data: make([][]byte, 25),
		r:    r,
		x:    0,
		y:    0,
	}

	r.Pages = append(r.Pages, p)
	r.CurrentPage = p
	if r.newpage != nil {
		r.newpage(p)
	}
	return p
}

// Append string
func (p *Page) Append(s string) *Page {
	for _, c := range s {
		p.AppendChar(byte(c))
	}
	return p
}

func (p *Page) AppendChars(s ...uint8) *Page {
	for _, c := range s {
		p.AppendChar(uint8(c))
	}
	return p
}

// Append char at position
func (p *Page) AppendChar(c uint8) *Page {
	for len(p.Data[p.y]) <= p.x {
		p.Data[p.y] = append(p.Data[p.y], ' ')
	}
	p.Data[p.y][p.x] = c
	p.x = p.x + 1
	if p.x >= 40 {
		return p.Newline()
	}
	return p
}

// Move to next line
func (p *Page) Newline() *Page {
	p.x = 0
	p.y = p.y + 1
	if p.y > 25 {
		p.y = 0
		return p.r.AddPage()
	}
	return p
}

// GraphChar returns the graphics character for a 6 bit pattern.
//
// Pattern matches the following bit set in v
// +-+-+
// |0|1|
// +-+-+
// |2|3|
// +-+-+
// |4|5|
// +-+-+
//
func GraphChar(v uint8) uint8 {
	return 160 + v
}
