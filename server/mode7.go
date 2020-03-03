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
	//p.Append(byte(len(r.Pages)))
	for _, page := range r.Pages {
		p.Append(page.Data...)
	}
	return p
}

func NewResult(newpage func(*Page)) *Pages {
	r := &Pages{newpage: newpage}
	r.AddPage()
	return r
}

type Page struct {
	Data []uint8 // 40x25 character display
	r    *Pages  // parent result
	x    int     // x position in screen
	y    int     // y position in screen
}

func (p *Page) offset() int {
	return p.x + (40 * p.y)
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

// Add a new page
func (r *Pages) AddPage() *Page {
	s := 24 * 40
	p := &Page{
		Data: make([]byte, s),
		r:    r,
		x:    0,
		y:    0,
	}

	// Fill page with spaces
	for i := 0; i < s; i++ {
		p.Data[i] = ' '
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
	p.Data[p.offset()] = c
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
