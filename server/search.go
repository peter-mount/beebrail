package server

import (
	"log"
)

func (s *Server) search(cmd Packet) *Packet {
	param := cmd.PayloadAsString()

	log.Println("Search", param)

	if len(param) < 3 {
		return cmd.EmptyResponse(1).
			AppendString("Search requires minimum 3 characters")
	}

	results, err := s.refClient.Search(param)
	if err != nil {
		log.Println(err)
		return cmd.ErrorPacket(err)
	}

	r := NewResult(func(p *Page) {
		m := "Search results"
		x := (40 - len(m)) >> 1
		for y := 1; y <= 2; y++ {
			p.Tab(0, y).
				AppendChars(AlphaBlue, NewBackground, AlphaWhite, DoubleHeight).
				Tab(x, y).
				Append(m)
		}
	})

	r.CurrentPage.Tab(0, 3).
		AppendChar(AlphaYellow).
		Append("CRS Station")

	for y, res := range results {
		if y < 20 {
			r.CurrentPage.Tab(0, 4+y)
			if res.Crs == param {
				r.CurrentPage.AppendChar(AlphaWhite)
			} else {
				r.CurrentPage.AppendChar(AlphaCyan)
			}
			r.CurrentPage.Append(res.Crs).
				Tab(5, -1).Append(res.Name)
			log.Println(res.Name, res.Label, res.Distance)
		}
	}

	return cmd.EmptyResponse(0).AppendPages(r)
}
