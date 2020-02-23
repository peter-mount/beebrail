package server

import (
	"fmt"
	refclient "github.com/peter-mount/nre-feeds/darwinref/client"
	"log"
	"strings"
)

func (s *Server) search(param string) {
	log.Println("Search", param)
	refClient := &refclient.DarwinRefClient{Url: "https://ref.prod.a51.li"}

	results, err := refClient.Search(param)
	var b []byte
	if err != nil {
		log.Println(err)
		b = append(b, []byte(err.Error())...)
	} else {
		for _, r := range results {
			b = append(b, []byte(fmt.Sprintf("%s %s%c", r.Crs, r.Name, 13))...)
		}
	}
	b = append(b, 0)

	log.Println(strings.ReplaceAll(string(b[:]), "\r", "\r\n"))

	_, err = s.port.Write(b)
	if err != nil {
		log.Println(err)
	}

}
