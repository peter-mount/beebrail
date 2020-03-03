package server

import (
	"log"
)

func (s *Server) boards(cmd Packet) *Packet {
	param := cmd.PayloadAsString()

	log.Println("CRS", param)

	sr, err := s.ldbClient.GetSchedule(param)
	if err != nil {
		return cmd.ErrorPacket(err)
	}

	r := cmd.EmptyResponse(0)

	r.AppendCString(sr.Crs)
	if len(sr.Station) == 0 {
		r.AppendCString(sr.Crs)
	} else {
		r.AppendCString(sr.Station[0])
	}

	r.AppendInt16(len(sr.Messages))
	for _, m := range sr.Messages {
		if m.Active {
			r.AppendInt16(m.Severity).AppendCString(m.Message)
		}
	}

	r.AppendInt16(len(sr.Messages))
	return r
}
