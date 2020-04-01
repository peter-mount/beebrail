package server

import (
	"errors"
	"fmt"
	"github.com/peter-mount/beebrail/server/util"
	"github.com/peter-mount/nre-feeds/darwind3"
	"github.com/peter-mount/nre-feeds/ldb"
	"github.com/peter-mount/nre-feeds/ldb/service"
	"regexp"
	"strings"
)

const (
	// Max pages of departures
	MAX_PAGES = 5
)

var stripHtml = []string{
	"<p>", "</p>",
}

var stripMore = []string{
	"More detail",
	"More information",
}

var (
	blueDoubleHeight = []byte{AlphaBlue, NewBackground, AlphaWhite, DoubleHeight}
)

func (s *Server) departures(r *ShellRequest) error {
	if len(r.Args) != 1 {
		return errors.New("Syntax: depart code")
	}
	crs := strings.ToUpper(r.Args[0])

	ctx := r.Context()
	ctx.Connection().Println("Departures", crs)

	sr, err := s.ldbClient.GetSchedule(crs)
	if err != nil {
		return err
	}

	var stationName string
	if len(sr.Station) == 0 {
		stationName = sr.Crs
	} else {
		stationName = sr.Station[0]
	}
	if d, ok := sr.Tiplocs.Get(stationName); ok {
		stationName = d.Name
	}

	w := r.ResultWriter().
		Title("CRS %s", crs).
		StxEtx(ctx.IsStxEtx())
	defer w.Close()

	t1 := r.NewTable().
		AddHeaders("Due", "Destination", "Plat")
	// Expected column is abbreviated for mode 7
	if t1.Style.Mode7 {
		t1.AddHeader("Expcted")
	} else {
		t1.AddHeader("Expected")
	}

	// Default minimum widths
	t1.Headers[0].Width = 4
	t1.Headers[1].Width = 17
	t1.Headers[2].Width = 4

	// Page header with Station Name
	t1.Callback.PageHeader = boardHeader(t1, stationName)

	include := true
	for _, s := range sr.Services {
		if include {
			processDeparture(crs, sr, s, t1)
		}
	}

	if len(t1.Rows) == 0 {
		t1.NewRow().Append("").Append("No services")
	}

	t2 := t1.NewTable().
		AddHeader("Message")

	// No table headers for message pages
	t2.Callback.TableHeader = func(o *util.ResultWriter) error {
		return o.WriteString("\n")
	}

	// We want 1 message per page
	t2.Style.PageHeight = 1

	for _, m := range sr.Messages {
		processMessage(m, t2)
	}

	return t1.Write(w)

	/*
		header := true
		r := NewResult(func(p *Page) {
			x := (40 - len(stationName)) >> 1
			for y := 1; y <= 2; y++ {
				p.Tab(0, y).
					AppendChars(AlphaBlue, NewBackground, AlphaWhite, DoubleHeight).
					Tab(x, y).
					Append(stationName)

				if header {
					p.Tab(0, y+2).
						AppendChars(AlphaBlue, NewBackground, AlphaWhite, DoubleHeight).
						Append("Due Destination       Plat  Expctd")
				}
			}
		})

		include := true
		for _, s := range sr.Services {
			if include {
				processDeparture(crs, sr, s, r)
				if r.PageCount() > MAX_PAGES {
					include = false
				}
			}
		}

		header = false
		for _, m := range sr.Messages {
			//if m.Active {
			s := m.Message

			// Strip out More detail... text
			for _, v := range stripMore {
				i := strings.Index(s, v)
				if i > -1 {
					s = s[:i]
				}
			}
			for _, v := range stripHtml {
				s = strings.ReplaceAll(s, v, "")
			}

			var v []string
			for len(s) > 37 {
				j := 37
				for s[j] != ' ' && j > 0 {
					j = j - 1
				}
				if j <= 0 {
					// Just split - should never happen
					v = append(v, s[:37])
					s = s[37:]
				} else {
					v = append(v, s[:j])
					if (j + 1) >= len(s) {
						s = ""
					} else {
						s = s[j+1:]
					}
				}
			}
			if s != "" {
				v = append(v, s)
			}

			p := r.AddPage()
			for y, s := range v {
				for i := 0; i < 2; i++ {
					p.Tab(0, 5+(2*y)+i).
						AppendChar(DoubleHeight).
						Append(s)
				}
			}
			//}
		}

		// Now run through and add page numering if required
		pageCount := len(r.Pages)
		if pageCount > 1 {
			for pn, p := range r.Pages {
				s := fmt.Sprintf("%d/%d", pn+1, pageCount)
				p.Tab(39-len(s), 1).
					Append(s).
					Tab(39-len(s), 2).
					Append(s)
			}
		}

		return cmd.EmptyResponse(0).AppendPages(r)
	*/
}

func processDeparture(crs string, sr *service.StationResult, s ldb.Service, t *util.Table) {
	if s.Location.Forecast.Suppressed {
		return
	}

	l := s.Location

	// Time minus :
	var tm string //tm := l.Time.String()
	if l.Times.Ptd != nil {
		tm = l.Times.Ptd.String()
	} else {
		tm = l.Times.Pta.String()
	}
	tm = tm[0:2] + tm[3:5]

	// Destination
	dest := s.Destination
	if d, ok := sr.Tiplocs.Get(dest); ok {
		if d.Crs == crs {
			return
		} else {
			dest = d.Name
		}
	}
	if len(dest) > 17 {
		dest = dest[:17]
	}

	// Platform, might be suppressed
	var plat string
	if !l.Forecast.Platform.Suppressed && !l.Forecast.Platform.CISSuppressed {
		plat = l.Forecast.Platform.Platform
	}

	var expected string
	expectedColour := steadyGreen

	forecast := l.Forecast.Departure

	if l.Forecast.Arrived {
		expected = "Arrived"
		expectedColour = steadyWhite
	} else if l.Cancelled {
		expected = "Canc'ld"
		expectedColour = steadyRed
	} else if forecast.Delayed {
		expected = "Delayed"
		expectedColour = steadyRed
	} else if l.Delay == 0 {
		expected = "On Time"
	} else if forecast.ET != nil {
		expected = forecast.ET.String()
		expected = expected[0:2] + expected[3:5] + " "
		// TODO if terminates here delay can show wrong as its using WTT not PTT in the calculation
		//log.Println(forecast.Time(), l.Delay, forecast.ET, l.Times)
		if l.Delay > 0 {
			expectedColour = flashYellow
		}
	}

	t.NewRow().
		Cell(&util.Cell{
			Label:  tm,
			Prefix: dblWhite,
		}).
		Cell(&util.Cell{
			Label:  dest,
			Prefix: white,
		}).
		Cell(&util.Cell{
			Label:  plat,
			Prefix: yellow,
		}).
		Cell(&util.Cell{
			Label:  expected,
			Prefix: expectedColour,
		})
}

var (
	dblWhite     = string([]byte{DoubleHeight, AlphaYellow})
	white        = string([]byte{AlphaWhite})
	yellow       = string([]byte{AlphaYellow})
	steadyWhite  = string([]byte{' ', AlphaWhite})
	steadyGreen  = string([]byte{' ', AlphaGreen})
	steadyRed    = string([]byte{' ', AlphaRed})
	steadyYellow = string([]byte{' ', AlphaYellow})
	flashYellow  = string([]byte{Flash, AlphaYellow})
)

func processMessage(m *darwind3.StationMessage, t *util.Table) {
	s := m.Message

	// Strip out More detail... text
	for _, v := range stripMore {
		i := strings.Index(s, v)
		if i > -1 {
			s = s[:i]
		}
	}
	for _, v := range stripHtml {
		s = strings.ReplaceAll(s, v, "")
	}

	// Remove html links
	re := regexp.MustCompile("(<a.+/a>)")
	msg := re.ReplaceAllString(s, "")

	var v []string
	for _, s := range strings.Split(msg, "\n") {
		s = strings.Trim(s, " ")
		for len(s) > 37 {
			j := 37
			for s[j] != ' ' && j > 0 {
				j = j - 1
			}
			if j <= 0 {
				// Just split - should never happen
				v = append(v, s[:37])
				s = s[37:]
			} else {
				v = append(v, s[:j])
				if (j + 1) >= len(s) {
					s = ""
				} else {
					s = s[j+1:]
				}
			}
		}
		s = strings.Trim(s, " ")
		if s != "" {
			v = append(v, s)
		}
	}

	t.NewRow().
		Append(strings.Join(v, "\n"))
}

func boardHeader(t1 *util.Table, stationName string) func(*util.Pagination, *util.ResultWriter) error {
	return func(p *util.Pagination, o *util.ResultWriter) error {

		page := p.Page()

		// Station name
		var title string
		if page.PageCount > 1 {
			title = fmt.Sprintf(" %d/%d", page.PageNo, page.PageCount)
		}
		w := (page.PageWidth - len(title) - len(stationName) - 1) / 2
		if t1.Style.Mode7 {
			w = w - len(blueDoubleHeight)
		}
		if w < 0 {
			w = 0
		}
		title = fmt.Sprintf(fmt.Sprintf("%%%ds%%s%%s\n", w), "", stationName, title)

		n := 1
		if t1.Style.Mode7 {
			n = 2
		}
		for i := 0; i < n; i++ {
			if t1.Style.Mode7 {
				_, err := o.Write(blueDoubleHeight)
				if err != nil {
					return err
				}
			}

			err := o.WriteString(title)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
