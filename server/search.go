package server

import (
	"fmt"
	refclient "github.com/peter-mount/nre-feeds/darwinref/client"
	"log"
)

func (s *Server) search(cmd *Packet) *Packet {
	param := cmd.PayloadAsString()

	log.Println("Search", param)

	if len(param) < 3 {
		return cmd.EmptyResponse(1).
			AppendString("Search requires minimum 3 characters")
	}

	refClient := &refclient.DarwinRefClient{Url: "https://ref.prod.a51.li"}

	results, err := refClient.Search(param)
	if err != nil {
		log.Println(err)
		return cmd.ErrorPacket(err)
	}

	resp := cmd.EmptyResponse(0)
	for _, r := range results {
		resp.AppendString(fmt.Sprintf("%s %s%c", r.Crs, r.Name, 13))
	}

	//log.Println(strings.ReplaceAll(string(b[:]), "\r", "\r\n"))

	return resp
}
