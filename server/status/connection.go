package status

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

// Connection
type Connection struct {
	cat      *Category     // Link to underlying category
	ID       int           `json:"id" xml:"id,attr"`         // The connection id
	Name     string        `json:"name" xml:"name,attr"`     // Name of connection, e.g. Serial port or IP address
	Local    Addr          `json:"local" xml:"local"`        // Local network port
	Remote   Addr          `json:"remote" xml:"remote"`      // Remote network port
	Secure   bool          `json:"secure" xml:"secure,attr"` // Secure connection
	Time     time.Time     `json:"time" xml:"time,attr"`     // When the connection was made
	Duration time.Duration `json:"duration" xml:"duration"`  // How long the connecton has been running
}

// A network address
type Addr struct {
	Interface string `json:"interface" xml:"interface,attr,omitempty"` // Interface bound to
	Port      uint16 `json:"port" xml:"port,attr,omitempty"`           // Port bound to
	Valid     bool   `json:"-" xml:"valid,attr"`                       // Flag to say address is valid, serial usually doesn't have one
}

func (a Addr) Append(ary []string) []string {
	if a.Valid {
		ary = append(ary, fmt.Sprintf("<td>%s</td><td>%d</td>", a.Interface, a.Port))
	} else {
		ary = append(ary, "<td></td><td></td>")
	}
	return ary
}

func (c *Connection) clone() *Connection {
	return &Connection{
		ID:       c.ID,
		Name:     c.Name,
		Local:    c.Local,
		Remote:   c.Remote,
		Time:     c.Time,
		Duration: time.Now().Sub(c.Time).Truncate(time.Second),
	}
}

func ExtractPort(addr net.Addr) (string, uint16, bool) {
	if addr != nil {
		a := addr.String()
		i := strings.LastIndex(a, ":")
		port, err := strconv.Atoi(a[i+1:])
		if err == nil && port > 0 && port < 65536 {
			a = a[:i]
			if len(a) > 2 && a[0] == '[' && a[len(a)-1] == ']' {
				a = a[1 : len(a)-1]
			}
			return a, uint16(port), true
		}
	}
	return "", 0, false

}

// Add a connection to a Category, return it or the existing one
func (c *Category) Add(local, remote net.Addr) *Connection {

	ln, lp, lv := ExtractPort(local)
	rn, rp, rv := ExtractPort(remote)

	var con *Connection

	_ = c.s.invoke(func() error {

		c.ConnectionCount++
		con = &Connection{
			ID: c.ConnectionCount,
			Local: Addr{
				Interface: ln,
				Port:      lp,
				Valid:     lv,
			},
			Remote: Addr{
				Interface: rn,
				Port:      rp,
				Valid:     rv,
			},
			Time: time.Now().UTC(),
			cat:  c,
		}
		c.connections[con.ID] = con

		log.Println("Add", c.Name, con.ID)

		return nil
	})

	return con
}

// Remove a Connection from it's Category
func (con *Connection) Remove() {
	if con != nil && con.cat != nil {
		_ = con.cat.s.invoke(func() error {
			delete(con.cat.connections, con.ID)
			log.Println("Rem", con.cat.Name, con.ID)
			con.cat = nil
			return nil
		})
	}
}
