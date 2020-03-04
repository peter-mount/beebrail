package server

import (
	"log"
)

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

	r := NewResult(func(p *Page) {
		x := (40 - len(stationName)) >> 1
		for y := 1; y <= 2; y++ {
			p.Tab(0, y).
				AppendChars(AlphaBlue, NewBackground, AlphaWhite, DoubleHeight).
				Tab(x, y).
				Append(stationName).
				Tab(0, y+2).
				AppendChars(AlphaBlue, NewBackground, AlphaWhite, DoubleHeight).
				Append("Due Destination       Plat  Expctd")
		}
	})

	for _, s := range sr.Services {
		if !s.Location.Forecast.Suppressed {
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
					dest = "Terminates here"
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
			d := l.Forecast.Departure
			if l.Forecast.Arrived {
				expected = "Arrived"
			} else if l.Cancelled {
				expected = "Canc'ld"
				expectedColour = AlphaRed
			} else if d.Delayed {
				expected = "Delayed"
				expectedColour = AlphaRed
				expectedFlash = Flash
			} else if l.Delay == 0 {
				expected = "On Time"
			} else if d.ET != nil {
				expected = d.ET.String()
				expected = expected[0:2] + expected[3:5] + " "
				if l.Delay > 0 {
					expectedColour = AlphaYellow
					expectedFlash = Flash
				}
			}

			p := r.CurrentPage
			for i := 0; i < 2; i++ {
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
	}

	for _, m := range sr.Messages {
		if m.Active {
			r.AddPage().
				Tab(0, 4).
				Append(m.Message)
		}
	}

	return cmd.EmptyResponse(0).AppendPages(r)
}
