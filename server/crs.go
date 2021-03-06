package server

import (
	"errors"
	"github.com/peter-mount/beebrail/server/util"
	"log"
	"strings"
)

func (s *Server) Crs(r *ShellRequest) error {
	defer func() {
		if r := recover(); r != nil {
			log.Println("PANIC", r)
		}
	}()
	return s.crs(r)
}
func (s *Server) crs(r *ShellRequest) error {
	if len(r.Args) != 1 {
		return errors.New("Syntax: crs code")
	}
	crs := strings.ToUpper(r.Args[0])

	ctx := r.Context()

	ctx.Connection().Println("CRS", crs)

	response, err := s.refClient.GetCrs(crs)
	if err != nil {
		return err
	}

	w := r.ResultWriter().
		Title("CRS %s", crs)
	defer w.Close()

	t := r.NewTable().
		AddHeaders("CRS", "Tiploc", "Toc", "Name")

	if response != nil {
		if ctx.IsAPI() {
			w.Footer("%d rows", len(response.Tiploc))
		}

		for _, tpl := range response.Tiploc {
			r := t.NewRow().
				Append(tpl.Crs).
				Append(tpl.Tiploc).
				Append(tpl.Toc).
				Append(tpl.Name)
			if t.Style.Mode7 {
				col := cyan
				if tpl.Crs == crs {
					col = white
				}
				for i := 0; i < 4; i++ {
					r.GetCell(i).Prefix = col
				}
			}
		}
	}

	if t.Style.Mode7 {
		t.Callback.TableHeader = func(o *util.ResultWriter) error {
			err := o.WriteBytes(AlphaBlue, NewBackground, AlphaWhite)
			if err == nil {
				err = t.WriteHeader(o)
			}
			return err
		}

		t.Callback.TableRow = func(t *util.Table, r *util.Row, o *util.ResultWriter) error {
			err := o.WriteString("  ")
			if err == nil {
				err = t.WriteRow(t, r, o)
			}
			return nil
		}

		t.Style.HCSep = AlphaWhite
	}

	_ = t.Write(w)

	return nil
	/*
		reader := NewResult(func(p *Page) {
			x := (40 - len(crs)) >> 1
			for y := 1; y <= 2; y++ {
				p.Tab(0, y).
					AppendChars(AlphaBlue, NewBackground, AlphaWhite, DoubleHeight).
					Tab(x, y).
					Append(crs)
			}
		})

		reader.CurrentPage.Tab(0, 3).
			AppendChar(AlphaYellow).
			Append("CRS Tiploc  Toc Name")

		if response != nil {
			for y, tpl := range response.Tiploc {
				if y < 20 {
					reader.CurrentPage.Tab(0, 4+y)
					if tpl.Crs == crs {
						reader.CurrentPage.AppendChar(AlphaWhite)
					} else {
						reader.CurrentPage.AppendChar(AlphaCyan)
					}
					reader.CurrentPage.Append(tpl.Crs).
						Tab(5, -1).Append(tpl.Tiploc).
						Tab(13, -1).Append(tpl.Toc).
						Tab(17, -1).Append(tpl.Name)
				}
			}
		}

		return cmd.EmptyResponse(0).AppendPages(reader)

	*/
}
