package server

import (
	"fmt"
	"github.com/peter-mount/nre-feeds/ldb"
	"github.com/peter-mount/nre-feeds/ldb/service"
	"log"
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

func (s *Server) departures(cmd Packet) *Packet {
	crs := cmd.PayloadAsString()

	log.Println("Departures", crs)

	sr, err := s.ldbClient.GetSchedule(crs)
	if err != nil {
		log.Println(err)
		return cmd.ErrorPacket(err)
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
}

func processDeparture(crs string, sr *service.StationResult, s ldb.Service, r *Pages) {
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
	expectedColour := uint8(AlphaGreen)
	expectedFlash := uint8(' ')

	forecast := l.Forecast.Departure

	if l.Forecast.Arrived {
		expected = "Arrived"
		expectedColour = AlphaWhite
	} else if l.Cancelled {
		expected = "Canc'ld"
		expectedColour = AlphaRed
	} else if forecast.Delayed {
		expected = "Delayed"
		expectedColour = AlphaRed
		expectedFlash = Flash
	} else if l.Delay == 0 {
		expected = "On Time"
	} else if forecast.ET != nil {
		expected = forecast.ET.String()
		expected = expected[0:2] + expected[3:5] + " "
		// TODO if terminates here delay can show wrong as its using WTT not PTT in the calculation
		//log.Println(forecast.Time(), l.Delay, forecast.ET, l.Times)
		if l.Delay > 0 {
			expectedColour = AlphaYellow
			expectedFlash = Flash
		}
	}

	p := r.CurrentPage
	for i := 0; i < 2; i++ {
		if p.Y() >= 22 {
			// If this means we go over the MAX_PAGES count then bail
			if r.PageCount() >= MAX_PAGES {
				return
			}
			p = r.AddPage()
		}
		y := p.Y() + 1

		p.Tab(0, y).
			AppendChars(DoubleHeight, AlphaYellow, ' ').
			Append(tm).
			AppendChar(AlphaWhite).
			Append(dest).
			Tab(24+5-len(plat), y).
			AppendChar(AlphaYellow).
			Append(plat).
			Tab(37-len(expected), y).
			AppendChars(expectedFlash, expectedColour).
			Append(expected)
	}
}
