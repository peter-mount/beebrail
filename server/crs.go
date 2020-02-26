package server

import (
	"fmt"
	"log"
)

func (s *Server) crs(cmd Packet) *Packet {
	param := cmd.PayloadAsString()

	log.Println("CRS", param)

	response, err := s.refClient.GetCrs(param)
	if err != nil {
		return cmd.ErrorPacket(err)
	}

	resp := cmd.EmptyResponse(0)

	for i, tpl := range response.Tiploc {
		if i > 0 {
			resp.Append(13)
		}
		resp.AppendString(fmt.Sprintf("%3.3s %-7.7s %-2.2s %s",
			tpl.Crs,
			tpl.Tiploc,
			tpl.Toc,
			tpl.Name,
		))
	}

	return resp
}
