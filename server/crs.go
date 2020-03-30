package server

import (
	"errors"
	"github.com/peter-mount/beebrail/server/util"
	"log"
	"strings"
)

func (s *Server) Crs(r *ShellRequest) error {
	if len(r.Args) != 1 {
		return errors.New("Syntax: crs code")
	}
	crs := strings.ToUpper(r.Args[0])

	log.Println("CRS", crs)

	response, err := s.refClient.GetCrs(crs)
	if err != nil {
		return err
	}

	t := util.Table{
		Title: "CRS " + crs,
		Style: util.Plain,
	}
	t.AddHeaders("CRS", "Tiploc", "Toc", "Name")

	if response != nil {
		for _, tpl := range response.Tiploc {
			t.NewRow().
				Append(tpl.Crs).
				Append(tpl.Tiploc).
				Append(tpl.Toc).
				Append(tpl.Name)
		}

		_ = t.Write(r.Stdout)
	}

	return nil
	/*
		r := NewResult(func(p *Page) {
			x := (40 - len(crs)) >> 1
			for y := 1; y <= 2; y++ {
				p.Tab(0, y).
					AppendChars(AlphaBlue, NewBackground, AlphaWhite, DoubleHeight).
					Tab(x, y).
					Append(crs)
			}
		})

		r.CurrentPage.Tab(0, 3).
			AppendChar(AlphaYellow).
			Append("CRS Tiploc  Toc Name")

		if response != nil {
			for y, tpl := range response.Tiploc {
				if y < 20 {
					r.CurrentPage.Tab(0, 4+y)
					if tpl.Crs == crs {
						r.CurrentPage.AppendChar(AlphaWhite)
					} else {
						r.CurrentPage.AppendChar(AlphaCyan)
					}
					r.CurrentPage.Append(tpl.Crs).
						Tab(5, -1).Append(tpl.Tiploc).
						Tab(13, -1).Append(tpl.Toc).
						Tab(17, -1).Append(tpl.Name)
				}
			}
		}

		return cmd.EmptyResponse(0).AppendPages(r)

	*/
}
